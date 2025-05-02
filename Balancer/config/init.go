package config

import (
	"cloud/Balancer/internal/service"
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() (*Config, error) {
	viper.AddConfigPath("../config") //Ищет в папке нужный файл
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	setDefaults() // Чтение переменных

	if err := viper.ReadInConfig(); err != nil {
		service.AppLogger.Printf("Ошибка чтения конфигурации: %v. Используются значения по умолчанию", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		service.ErrorLogger.Printf("Ошибка парсинга конфигурации: %v", err)
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	service.AppLogger.Printf("Загружена конфигурация: порт - %d, backends - %v",
		cfg.HTTPServer.Port, cfg.Backends)

	return &cfg, nil
}
