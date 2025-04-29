package main

import (
	"cloud/Balancer/config"
	"cloud/Balancer/internal/service"
	"fmt"
	"log"
)

func main() {
	err := service.InitLogger()
	defer service.Close()

	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.InitConfig()

	if err != nil {
		service.ErrorLogger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	service.AppLogger.Println("Hello")
	service.ErrorLogger.Println("Crit")

	fmt.Printf("%+v\n", cfg)

}
