package tv

import (
	"encoding/json"
	"errors"
	"moskitbot/internal/client"
	"moskitbot/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

// Socket ...
type Socket struct {
	conn      *websocket.Conn
	isClosed  bool
	sessionID string
	logger    *logging.Logger
	pairs     []string
	levels    map[string][]float32
}

// Connect - Connects and returns the trading view socket object
func Connect(logger *logging.Logger) (socket *Socket, err error) {
	conn, _, err := (&websocket.Dialer{}).Dial("wss://data.tradingview.com/socket.io/websocket", getHeaders())
	if err != nil {
		logger.Errorf("error while socKet connecting: %s", err)
		return
	}

	socket = &Socket{
		conn:     conn,
		logger:   logger,
		isClosed: false,
	}
	return
}

// InitWatching initializes data tracking
func (s *Socket) InitWatching(symbols []string, errorChan chan error, responseChan chan client.Response) (err error) {
	err = s.checkFirstReceivedMessage()
	if err != nil {
		return
	}
	s.generateSessionID()

	err = s.sendConnectionSetupMessages()
	if err != nil {
		s.onError()
		s.logger.Errorf("error while sending setup message: %s", err)
		return
	}

	for _, symbol := range symbols {
		if err = s.AddSymbol(symbol); err != nil {
			s.logger.Errorf("error while adding symbol %s: %s", symbol, err)
			return
		}
	}

	go s.connectionLoop(errorChan, responseChan)

	return
}

// Close ...
func (s *Socket) Close() (err error) {
	s.isClosed = true
	return s.conn.Close()
}

// AddSymbol ...
func (s *Socket) AddSymbol(symbol string) (err error) {
	err = s.sendSocketMessage(
		getSocketMessage("quote_add_symbols", []interface{}{s.sessionID, symbol, getFlags()}),
	)
	return
}

// RemoveSymbol ...
func (s *Socket) RemoveSymbol(symbol string) (err error) {
	err = s.sendSocketMessage(
		getSocketMessage("quote_remove_symbols", []interface{}{s.sessionID, symbol}),
	)
	return
}

func (s *Socket) checkFirstReceivedMessage() (err error) {
	var msg []byte

	_, msg, err = s.conn.ReadMessage()
	if err != nil {
		s.logger.Errorf("error while checking first message: %s", err)
		return
	}

	payload := msg[getPayloadStartingIndex(msg):]
	var p map[string]interface{}

	err = json.Unmarshal(payload, &p)
	if err != nil {
		s.logger.Errorf("error while decoding first message: %s", err)
		return
	}

	if p["session_id"] == nil {
		err = errors.New("cannot recognize the first received message after establishing the connection")
		s.logger.Errorf(err.Error())
		return
	}

	return
}

func (s *Socket) generateSessionID() {
	s.sessionID = "qs_" + GetRandomString(12)
}

func (s *Socket) sendConnectionSetupMessages() (err error) {
	messages := []*SocketMessage{
		getSocketMessage("set_auth_token", []string{"unauthorized_user_token"}),
		getSocketMessage("quote_create_session", []string{s.sessionID}),
		getSocketMessage("quote_set_fields", []string{s.sessionID, "lp", "volume", "bid", "ask"}),
	}

	for _, msg := range messages {
		err = s.sendSocketMessage(msg)
		if err != nil {
			return
		}
	}

	return
}

func (s *Socket) sendSocketMessage(p *SocketMessage) (err error) {
	payload, _ := json.Marshal(p)
	payloadWithHeader := "~m~" + strconv.Itoa(len(payload)) + "~m~" + string(payload)

	err = s.conn.WriteMessage(websocket.TextMessage, []byte(payloadWithHeader))
	if err != nil {
		s.logger.Errorf("error while sending socket message: %s", err)
		return
	}
	return
}

func (s *Socket) connectionLoop(errorChan chan error, responseChan chan client.Response) {
	var readMsgError error
	var writeKeepAliveMsgError error

	for readMsgError == nil && writeKeepAliveMsgError == nil {
		if s.isClosed {
			break
		}

		var msgType int
		var msg []byte
		msgType, msg, readMsgError = s.conn.ReadMessage()

		go func(msgType int, msg []byte) {
			if msgType != websocket.TextMessage {
				return
			}

			if isKeepAliveMsg(msg) {
				writeKeepAliveMsgError = s.conn.WriteMessage(msgType, msg)
				return
			}

			response, err := s.parsePacket(msg)
			if err != nil {
				errorChan <- readMsgError
			} else {
				responseChan <- response
			}
		}(msgType, msg)
	}

	if readMsgError != nil {
		errorChan <- readMsgError
	}
	if writeKeepAliveMsgError != nil {
		errorChan <- writeKeepAliveMsgError
	}
}

func (s *Socket) parsePacket(packet []byte) (client.Response, error) {
	var resp = make(map[string]*float64)

	index := 0
	for index < len(packet) {
		payloadLength, err := getPayloadLength(packet[index:])
		if err != nil {
			return client.Response{}, err
		}

		headerLength := 6 + len(strconv.Itoa(payloadLength))
		payload := packet[index+headerLength : index+headerLength+payloadLength]
		index = index + headerLength + len(payload)

		symbol, data, err := s.parseJSON(payload)
		if err != nil {
			return client.Response{}, err
		}
		resp[symbol] = data.Price

	}
	return resp, nil
}

func (s *Socket) parseJSON(msg []byte) (symbol string, data *QuoteData, err error) {
	var decodedMessage *SocketMessage

	err = json.Unmarshal(msg, &decodedMessage)
	if err != nil {
		s.logger.Errorf("error while decoding message: %s", err)
		return
	}

	if decodedMessage.Message == "critical_error" || decodedMessage.Message == "error" {
		err = errors.New("Error -> " + string(msg))
		s.logger.Errorf("error while decoding message: %s", err)
		return
	}

	if decodedMessage.Message != "qsd" {
		err = errors.New("ignored message - Not QSD")
		return
	}

	if decodedMessage.Payload == nil {
		err = errors.New("Msg does not include 'p' -> " + string(msg))
		s.logger.Errorf("Decoded message does not include payload: %s", err)
		return
	}

	p, isPOk := decodedMessage.Payload.([]interface{})
	if !isPOk || len(p) != 2 {
		err = errors.New("There is something wrong with the payload - can't be parsed -> " + string(msg))
		s.logger.Errorf(err.Error())
		return
	}

	var decodedQuoteMessage *QuoteMessage
	err = mapstructure.Decode(p[1].(map[string]interface{}), &decodedQuoteMessage)
	if err != nil {
		s.logger.Errorf("payload can not be parsed: %s", err)
		return
	}

	if decodedQuoteMessage.Status != "ok" || decodedQuoteMessage.Symbol == "" || decodedQuoteMessage.Data == nil {
		err = errors.New("There is something wrong with the payload - couldn't be parsed -> " + string(msg))
		s.logger.Errorf(err.Error())
		return
	}
	symbol = decodedQuoteMessage.Symbol
	data = decodedQuoteMessage.Data
	return
}

func (s *Socket) onError() {
	if s.conn != nil {
		s.conn.Close()
		s.isClosed = true
	}
}

func (s *Socket) addNewLine(newLine Line) error {
	err := s.AddSymbol(newLine.Pair)
	if err != nil {
		return err
	}
	s.pairs = append(s.pairs, newLine.Pair)

}

func getSocketMessage(m string, p interface{}) *SocketMessage {
	return &SocketMessage{
		Message: m,
		Payload: p,
	}
}

func getFlags() *Flags {
	return &Flags{
		Flags: []string{"force_permission"},
	}
}

func isKeepAliveMsg(msg []byte) bool {
	return string(msg[getPayloadStartingIndex(msg)]) == "~"
}

func getPayloadStartingIndex(msg []byte) int {
	char := ""
	index := 3
	for char != "~" {
		char = string(msg[index])
		index++
	}
	index += 2
	return index
}

func getPayloadLength(msg []byte) (length int, err error) {
	char := ""
	index := 3
	lengthAsString := ""
	for char != "~" {
		char = string(msg[index])
		if char != "~" {
			lengthAsString += char
		}
		index++
	}
	length, err = strconv.Atoi(lengthAsString)
	return
}

func getHeaders() http.Header {
	headers := http.Header{}

	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "en-US,en;q=0.9,es;q=0.8")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Host", "data.tradingview.com")
	headers.Set("Origin", "https://www.tradingview.com")
	headers.Set("Pragma", "no-cache")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.193 Safari/537.36")

	return headers
}
