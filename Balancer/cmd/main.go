package main

import (
	"cloud/Balancer/config"
	"cloud/Balancer/internal/handler/middlewares"
	"cloud/Balancer/internal/models"
	"cloud/Balancer/internal/service"
	"fmt"
	"log"
	"net/http"
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

	balancer := models.NewSimpleBalancer(cfg.Backends)

	fmt.Println(balancer)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		Handler: middlewares.LoggerMiddleware(),
	}

	if server.ListenAndServe() != nil {
		service.ErrorLogger.Fatalf("Ошибка сервера : %v", err)
	}

	fmt.Printf("%+v\n", cfg)

}
