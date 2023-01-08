package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"moskitbot/pkg/logging"
	"sync"
)

var logger = logging.GetLogger()

type Config struct {
	IsDebug             bool   `yaml:"is_debug" env:"IS_DEBUG"  env-default:"true"`
	MoskitToken         string `yaml:"moskit_token" env:"MOSKIT_TOKEN" env-default:""`
	AvailableTimeframes string `yaml:"available_timeframes" env:"AVAILABLE_TIMEFRAMES" env-default:"1m, 5m, 15m, 30m, 1h, 2h, 4h, 1d, 1w, 1M"`
	ChatID              int64  `yaml:"chat_id" env:"CHAT_ID" env_default:"-1001857192414"`
	DBHost              string `yaml:"db_host" env:"DB_HOST"  env-default:"0.0.0.0"`
	DBPort              string `yaml:"db_port" env:"DB_PORT"  env-default:"3306"`
	DBUsername          string `yaml:"mysql_user" env:"MYSQL_USER"  env-default:"moskit"`
	DBPassword          string `yaml:"mysql_password" env:"MYSQL_PASSWORD"  env-default:"moskitpassword"`
	DBName              string `yaml:"db_name" env:"DB_NAME"  env-default:"moskit"`
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
