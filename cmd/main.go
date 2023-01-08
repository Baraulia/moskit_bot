package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/client"
	"moskitbot/internal/config"
	"moskitbot/internal/line"
	"moskitbot/internal/repository"
	"moskitbot/internal/tv"
	"moskitbot/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	mysqlDB, err := repository.NewMysqlDB(repository.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "PREFERRED",
	})
	if err != nil {
		logger.Panicf("Error while initialization database:%s", err)
	}

	repo := repository.NewLineRepository(mysqlDB, logger)

	bot, err := tgbotapi.NewBotAPI(cfg.MoskitToken)
	if err != nil {
		logger.Panicf("Error while creating new bot API:%s", err)
	}

	var newLine = make(chan tv.Line)

	tgBot := line.NewBot(bot, cfg.ChatID, logger, newLine)

	tvsocket, err := tv.Connect(logger)
	if err != nil {
		logger.Panicf("Error while initializing the trading view socket -> " + err.Error())
	}

	cli := client.NewClient(repo, tvsocket, logger, tgBot, newLine)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err = cli.StartAndServe(ctx); err != nil {
			logger.Panicf(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	cancel()

}
