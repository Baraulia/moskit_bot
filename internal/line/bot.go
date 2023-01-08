package line

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/tv"
	"moskitbot/pkg/logging"
)

type Bot struct {
	Bot     *tgbotapi.BotAPI
	ChatID  int64
	logger  *logging.Logger
	newLine chan tv.Line
}

func NewBot(bot *tgbotapi.BotAPI, chatID int64, logger *logging.Logger, newLine chan tv.Line) *Bot {
	return &Bot{Bot: bot, ChatID: chatID, logger: logger, newLine: newLine}
}

func (b *Bot) Start() error {
	updates := b.initUpdatesChannel()

	if err := b.handleUpdates(updates); err != nil {
		b.logger.Errorf("handleUpdates: %s", err)
		return fmt.Errorf("handleUpdates: %w", err)
	}
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil { // Ignore any non_Message Updates
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.logger.Errorf("handleCommand:%s", err)
				return fmt.Errorf("handleCommand:%w", err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			b.logger.Errorf("handleMessage:%s", err)
			return fmt.Errorf("handleMessage:%w", err)
		}
	}
	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return b.Bot.GetUpdatesChan(u)
}
