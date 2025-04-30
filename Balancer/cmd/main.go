package main

import (
	"cloud/Balancer/config"
	"cloud/Balancer/internal/handler/middlewares"
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

	balancer, err := service.NewSimpleBalancer(cfg.Backends)

	service.StartBackends()

	go service.LiveCheck()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		balancer.Proxy.ServeHTTP(w, r)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		Handler: middlewares.LoggerMiddleware(handler),
	}

	if server.ListenAndServe() != nil {
		service.ErrorLogger.Fatalf("Ошибка сервера : %v", err)
	}
}
