package models

import (
	"app/internal/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Config *ConfigStruct

type ConfigStruct struct {
	WebAppConfig
	DataBaseConfig
	CacheConfig
}

type WebAppConfig struct {
	AppIP   string `envconfig:"APP_IP" default:"localhost"`
	AppPort string `envconfig:"APP_PORT" default:"8000"`
	URL     string `envconfig:"URL" required:"true"`
}

type DataBaseConfig struct {
	UserName     string `envconfig:"DBUSER" required:"true"`
	Password     string `envconfig:"DBPASS" required:"true"`
	Host         string `envconfig:"DBHOST" required:"true"`
	Port         string `envconfig:"DBPORT" required:"false"`
	DataBaseName string `envconfig:"DBNAME" required:"true"`
}

type CacheConfig struct {
	TTL              int `envconfig:"TTL" required:"true"`
	HardMaxCacheSize int `envconfig:"HARD_MAX_CACHE_SIZE" required:"true"`
	MaxEntrySizes    int `envconfig:"MAX_ENTRY_SIZE" required:"true"`
	Shards           int `envconfig:"SHARDS" required:"true"`
}

func InitConfig() error {
	// Загрузка файла .env
	if err := godotenv.Load(); err != nil {
		logger.Log.Error("Загрузка приложения без env файла")
	}

	Config = &ConfigStruct{}
	err := envconfig.Process("", Config)
	if err != nil {
		return err
	}
	logger.Log.Info("Загруженые параметры: \n", Config)
	return nil
}
