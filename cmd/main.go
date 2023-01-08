package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"moskitbot/internal/config"
	"moskitbot/internal/line"
	"moskitbot/internal/tv"
	"moskitbot/internal/tvscrapper"
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

	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	//
	//<-quit
	//cancel()
	tvsocket, err := tvscrapper.Connect(
		func(symbol string, data *tvscrapper.QuoteData) {
			fmt.Println(*data.Price)
		},
		func(err error, context string) {
			fmt.Printf("%#v", "error -> "+err.Error())
			fmt.Printf("%#v", "context -> "+context)
		},
	)
	if err != nil {
		panic("Error while initializing the trading view socket -> " + err.Error())
	}

	err = tvsocket.AddSymbol("BINANCE:BTCUSDT")
	if err != nil {
		return
	}
	err = tvsocket.AddSymbol("BINANCE:ETHUSDT")
	if err != nil {
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

}
