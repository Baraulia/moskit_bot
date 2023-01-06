package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/config"
	"moskitbot/internal/line"
	"moskitbot/internal/tv"
	"moskitbot/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.MoskitToken)
	if err != nil {
		logger.Panicf("Error while creating new bot API:%s", err)
	}

	tgBot := line.NewBot(bot, cfg.ChatID)

	tvClient := tv.NewClient(cfg.TVURI, tgBot, logger)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := tvClient.StartAndServe(ctx); err != nil {
			logger.Panicf(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	cancel()

}
