package config

import (
	"cloud/Balancer/internal/service"
	"github.com/spf13/viper"
)

func InitConfig() (*Config, error) {
	viper.AddConfigPath("../config") // Ищет в папке нужный yaml файл
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // Чтение переменных

	if err := viper.ReadInConfig(); err != nil {
		service.AppLogger.Printf("Ошибка чтения конфигурации: %v. Будут выставленны значения по умолчанию", err)
		viper.SetDefault("http_server.port", 8080)
		viper.SetDefault("backends", []string{
			"http://localhost:8081",
			"http://localhost:8082",
			"http://localhost:8083",
			"http://localhost:8084",
			"http://localhost:8085",
		})
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		service.ErrorLogger.Printf("Произошла ошибка в парсинге файла %v", err)
		return nil, err
	}

	service.AppLogger.Println("Конфигурация: порт - %v, backends - %v", cfg.HTTPServer.Port, cfg.Backends)

	return &cfg, nil

}
