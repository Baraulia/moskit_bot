package client

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/line"
	"moskitbot/internal/repository"
	"moskitbot/internal/tvscrapper"
	"moskitbot/pkg/logging"
)

type Client struct {
	repo   repository.Repository
	socket *tvscrapper.Socket
	logger *logging.Logger
	bot    *line.Bot
}

type Response map[string]*float64

func (c *Client) StartAndServe(ctx context.Context) error {
	responseChan := make(chan Response)
	errorChan := make(chan error)

	lines, symbols, err := c.repo.GetAll()

	go c.socket.InitWatching(symbols, errorChan, responseChan)

	msg := tgbotapi.NewMessage(c.bot.ChatID, fmt.Sprintf(`*Start watching\.\.\.*`))
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := c.bot.Bot.Send(msg)
	if err != nil {
		return err
	}

	for {
		select {
		case err = <-errorChan:
			msg = tgbotapi.NewMessage(c.bot.ChatID, fmt.Sprintf("Ошибка: %s", err))
			_, err = c.bot.Bot.Send(msg)
			if err != nil {
				return err
			}
		case resp := <-responseChan:
			text := fmt.Sprintf("Цена (%v) пары %s приблизилась k уровню %s. Примечание: %s", resp.value, resp.instrument, resp.lineType, resp.description)
			msg = tgbotapi.NewMessage(c.bot.ChatID, text)
			_, err := c.bot.Bot.Send(msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
