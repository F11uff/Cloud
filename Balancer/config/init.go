package config

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig() (*Config, error) {
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Ошибка чтения конфига: %v. Будут выставленны значения по умолчанию", err)
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
		log.Printf("Произошла ошибка в парсинге файла %v", err)
		return nil, err
	}

	log.Printf("Конфигурация: порт - %v, backends - %v", cfg.HTTPServer.Port, cfg.Backends)

	return &cfg, nil

}
