package tv

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
	"moskitbot/internal/line"
	"moskitbot/internal/repository"
	"moskitbot/pkg/logging"
	"time"
)

var ErrInvalidScreener = errors.New("invalid screener")
var ErrNoSymbols = errors.New("no symbols given")

type Client struct {
	logger   *logging.Logger
	uri      string
	watchers []Watcher
	bot      *line.Bot
}

type Alarm struct {
	instrument  string
	lineType    string
	value       float32
	description string
}

type Request struct {
	Symbols struct {
		Tickers []string `json:"tickers"`
		Query   struct {
			Types []interface{} `json:"types"`
		} `json:"query"`
	} `json:"symbols"`
	Columns []string `json:"columns"`
}

type Response struct {
	Data []struct {
		Instrument   string            `json:"s"`
		ColumnsValue []decimal.Decimal `json:"d"`
	} `json:"data"`
}

func NewClient(uri string, bot *line.Bot, logger *logging.Logger, repo repository.Repository, cache *redis.Client) *Client {
	var watchers []Watcher

	rsiWatcher := NewLineWatcher(
		1*time.Minute,
		logger,
		uri,
		repo,
		cache,
	)
	watchers = append(watchers, rsiWatcher)

	return &Client{
		logger:   logger,
		uri:      uri,
		bot:      bot,
		watchers: watchers,
	}
}

func (c *Client) StartAndServe(ctx context.Context) error {
	responseChan := make(chan Alarm)
	errorChan := make(chan error)

	for _, watcher := range c.watchers {
		go watcher.StartWatching(ctx, errorChan, responseChan)
	}

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
