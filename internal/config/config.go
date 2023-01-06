package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"moskitbot/pkg/logging"
	"sync"
)

var logger = logging.GetLogger()

type Config struct {
	IsDebug             bool   `yaml:"is_debug" env:"IS_DEBUG"  env-default:"true"`
	BotToken            string `yaml:"bot_token" env:"BOT_TOKEN" env-default:""`
	MoskitToken         string `yaml:"moskit_token" env:"MOSKIT_TOKEN" env-default:""`
	AvailableTimeframes string `yaml:"available_timeframes" env:"AVAILABLE_TIMEFRAMES" env-default:"1m, 5m, 15m, 30m, 1h, 2h, 4h, 1d, 1w, 1M"`
	TVURI               string `yaml:"tv_uri" env:"TV_URI" env-default:"https://scanner.tradingview.com/crypto/scan"`
	ChatID              int64  `yaml:"chat_id" env:"CHAT_ID" env_default:"-1001857192414"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger.Info("read application configuration")
		instance = &Config{}

		err := cleanenv.ReadConfig("configs/config.yaml", instance)
		if err != nil {
			logger.Info(err)
			if err := cleanenv.ReadConfig(".env", instance); err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				logger.Info(help)
				logger.Fatal(err)
			}
		}

	})
	return instance
}
