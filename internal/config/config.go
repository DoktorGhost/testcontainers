package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	DB_HOST  string `env:"DB_HOST"`
	DB_PORT  string `env:"DB_PORT"`
	DB_NAME  string `env:"DB_NAME"`
	DB_LOGIN string `env:"DB_LOGIN"`
	DB_PASS  string `env:"DB_PASS"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	//считываем все переменны окружения в cfg
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("Ошибка загрузки конфигурации: %v", err)
	}

	return config, nil
}
