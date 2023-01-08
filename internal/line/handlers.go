package line

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/tv"
	"strconv"
	"strings"
)

const commandStart = "start"

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	b.logger.Printf("working command: %s", message.Command())
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неподдерживаемая команда")
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	switch message.Command() {
	case commandStart:
		msg.Text = fmt.Sprintf(`Введите цену которую необходимо отслеживать в формате: pair/type/timeframe/value/description
	*"BTCUSDT/BearOB/5m/16200/Подойдя к этой цене ждем отскока"*`)
	default:
		msg.Text = "Неподдерживаемая команда"
	}
	_, err := b.Bot.Send(msg)
	if err != nil {
		b.logger.Errorf("send message to Telegram:%s", err)
		return fmt.Errorf("send message to Telegram:%w", err)
	}
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	line, err := b.parseMessage(message.Text)
	if err != nil {
		b.logger.Errorf("error while parsing message: %s", err)
		return err
	}

	b.newLine <- line

	return nil
}

func (b *Bot) parseMessage(message string) (tv.Line, error) {
	sliceParam := strings.Split(message, "/")
	if len(sliceParam) != 5 {
		return tv.Line{}, fmt.Errorf("invalid input data")
	}
	value, err := strconv.ParseFloat(sliceParam[3], 64)
	if err != nil {
		return tv.Line{}, err
	}
	return tv.Line{
		Pair:        fmt.Sprintf("BINANCE:%s", strings.ReplaceAll(sliceParam[0], " ", "")),
		Val:         value,
		Description: sliceParam[4],
		Typ:         sliceParam[1],
		Timeframe:   sliceParam[2],
	}, nil
}
