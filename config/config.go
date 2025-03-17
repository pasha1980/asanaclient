package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

type Config struct {
	AsanaAccessToken string `envconfig:"ASANA_ACCESS_TOKEN"`
	AsanaBaseUrl     string `envconfig:"ASANA_BASE_URL" default:"https://app.asana.com/api/1.0"`
	AsanaLimit       int    `envconfig:"ASANA_LIMIT" default:"5"`

	StorageBasePath string `envconfig:"STORAGE_BASE_PATH" default:"./storage"`
}

var c *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		var cfg Config
		_ = godotenv.Load()

		if err := envconfig.Process("", &cfg); err != nil {
			return
		}

		c = &cfg
	})

	return c
}
