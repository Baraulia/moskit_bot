package client

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/line"
	"moskitbot/internal/repository"
	"moskitbot/internal/tv"
	"moskitbot/pkg/logging"
)

type Client struct {
	repo    repository.Repository
	socket  *tv.Socket
	logger  *logging.Logger
	bot     *line.Bot
	lines   []tv.Line
	newLine chan tv.Line
}

func NewClient(repo repository.Repository, socket *tv.Socket, logger *logging.Logger, bot *line.Bot, newLine chan tv.Line) *Client {
	return &Client{
		repo:    repo,
		socket:  socket,
		logger:  logger,
		bot:     bot,
		newLine: newLine,
	}
}

type Response map[string]*float64

func (c *Client) StartAndServe(ctx context.Context) error {
	responseChan := make(chan Response)
	errorChan := make(chan error)

	lines, symbols, err := c.repo.GetAll()
	c.lines = lines

	go func() {
		err = c.socket.InitWatching(symbols, errorChan, responseChan)
		if err != nil {
			errorChan <- err
			return
		}
	}()

	msg := tgbotapi.NewMessage(c.bot.ChatID, fmt.Sprintf(`*Start watching\.\.\.*`))
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err = c.bot.Bot.Send(msg)
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
			for key, value := range resp {
				for _, l := range c.lines {
					if l.Pair == key {
						if *value >= *value*0.999 && *value <= *value*1.001 {
							text := fmt.Sprintf("Цена (%v) пары %s приблизилась k %s %s. Примечание: %s", value, l.Pair, l.Typ, l.Timeframe, l.Description)
							msg = tgbotapi.NewMessage(c.bot.ChatID, text)
							_, err = c.bot.Bot.Send(msg)
							if err != nil {
								return err
							}
							count, err := c.repo.Delete(l.ID)
							if err != nil {
								return err
							} else if count == 0 {
								return fmt.Errorf("can not delete line with id = %d", l.ID)
							}
						}
					}
				}
			}
		case newLine := <-c.newLine:
			id, err := c.repo.Create(newLine)
			if err != nil {
				return err
			}
			newLine.ID = id
			c.socket.AddNewLine(newLine)

		case <-ctx.Done():
			return nil
		}
	}
}
