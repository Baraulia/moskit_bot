package line

//
//import (
//	"context"
//	"fmt"
//	tradingview "github.com/Baraulia/MyTelBot.git/internal/tv"
//	"github.com/go-line-bot-api/line-bot-api/v5"
//	"strings"
//)
//
//const commandStart = "start"
//
//var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//	bot.NewInlineKeyboardRow(
//		bot.NewInlineKeyboardButtonData("BTC", "BTC"),
//		bot.NewInlineKeyboardButtonData("BNB", "BNB"),
//		bot.NewInlineKeyboardButtonData("ETH", "ETH"),
//	),
//	bot.NewInlineKeyboardRow(
//		bot.NewInlineKeyboardButtonData("TRX", "TRX"),
//		bot.NewInlineKeyboardButtonData("BCH", "BCH"),
//		bot.NewInlineKeyboardButtonData("XRP", "XRP"),
//	),
//)
//
//func (b *Bot) handleCommand(message *bot.Message) error {
//	b.logger.Printf("working command: %s", message.Command())
//	msg := tgbotapi.NewMessage(message.Chat.ID, "Неподдерживаемая команда")
//	msg.ParseMode = bot.ModeMarkdownV2
//	switch message.Command() {
//	case commandStart:
//		msg.Text = fmt.Sprintf(`Введите название биржы, валютную пару и таймфрейм в соответствии с примером:
//	*"binance:BTCUSDT/5m"*
//__Доступные таймфреймы:__ %s`, b.timeFrames)
//	default:
//		msg.Text = "Неподдерживаемая команда"
//	}
//	_, err := b.bot.Send(msg)
//	if err != nil {
//		b.logger.Errorf("send message to Telegram:%s", err)
//		return fmt.Errorf("send message to Telegram:%w", err)
//	}
//	return nil
//}
//
//func (b *Bot) handleMessage(message *bot.Message) error {
//	var symbols []string
//	pair, timeframe, err := b.parseMessage(message.Text)
//	if err != nil {
//		b.logger.Errorf("error while parsing message: %s", err)
//		return err
//	}
//
//	symbols = append(symbols, pair)
//
//	analysis, err := b.tv.GetAnalysis(context.Background(), "crypto", symbols, tradingview.Interval(timeframe))
//	if err != nil {
//		b.logger.Errorf("error while getting data from tradingview: %s", err)
//		return err
//	}
//
//	msg := bot.NewMessage(message.Chat.ID, fmt.Sprintf("RSI на таймфрейме %s: %d", timeframe, analysis.IntPart()))
//
//	_, err = b.bot.Send(msg)
//	if err != nil {
//		b.logger.Errorf("send message to Telegram:%s", err)
//		return fmt.Errorf("send message to Telegram:%w", err)
//	}
//
//	return nil
//
//}
//
//func (b *Bot) parseMessage(message string) (string, string, error) {
//	sliceParam := strings.Split(strings.ReplaceAll(message, " ", ""), "/")
//	if len(sliceParam) != 2 {
//		return "", "", fmt.Errorf("invalid input data")
//	}
//
//	if !strings.Contains(b.timeFrames, sliceParam[1]) {
//		return "", "", fmt.Errorf("invalid input timeframe")
//	}
//
//	return strings.ToUpper(sliceParam[0]), sliceParam[1], nil
//}
