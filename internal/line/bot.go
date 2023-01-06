package line

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Bot    *tgbotapi.BotAPI
	ChatID int64
}

func NewBot(bot *tgbotapi.BotAPI, chatID int64) *Bot {
	return &Bot{Bot: bot, ChatID: chatID}
}

//func (b *Bot) Start() error {
//	b.logger.Infof("Authorized on account %s", b.bot.Self.UserName)
//
//	updates := b.initUpdatesChannel()
//
//	if err := b.handleUpdates(updates); err != nil {
//		b.logger.Errorf("handleUpdates: %s", err)
//		return fmt.Errorf("handleUpdates: %w", err)
//	}
//	return nil
//}

//func (b *Bot) handleUpdates(updates bot.UpdatesChannel) error {
//	for update := range updates {
//		if update.Message == nil { // Ignore any non_Message Updates
//			continue
//		}
//		if update.Message.IsCommand() {
//			if err := b.handleCommand(update.Message); err != nil {
//				b.logger.Errorf("handleCommand:%s", err)
//				return fmt.Errorf("handleCommand:%w", err)
//			}
//			continue
//		}
//		if err := b.handleMessage(update.Message); err != nil {
//			b.logger.Errorf("handleMessage:%s", err)
//			return fmt.Errorf("handleMessage:%w", err)
//		}
//	}
//	return nil
//}
//
//func (b *Bot) initUpdatesChannel() bot.UpdatesChannel {
//	u := bot.NewUpdate(0)
//	u.Timeout = 60
//	return b.bot.GetUpdatesChan(u)
//}
