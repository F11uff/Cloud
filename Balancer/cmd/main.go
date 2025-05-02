package main

import (
	"cloud/Balancer/config"
	"cloud/Balancer/internal/API"
	"cloud/Balancer/internal/handler/middlewares"
	"cloud/Balancer/internal/service"
	"cloud/Balancer/internal/storage/PSQL"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func ensureDatabaseExists(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/postgres?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PostgreSQL: %v", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		cfg.DB.Name,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("ошибка проверки существования БД: %v", err)
	}

	if !exists {
		// Создаем новую БД
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DB.Name))
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
				return nil
			}
			return fmt.Errorf("не удалось создать БД: %v", err)
		}
		service.AppLogger.Printf("База данных %s успешно создана", cfg.DB.Name)
	}

	return nil
}

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

	if err := ensureDatabaseExists(cfg); err != nil {
		service.ErrorLogger.Fatalf("Ошибка инициализации БД: %v", err)
	}

	pgStorage, err := PSQL.NewStorage(fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Name,
	))
	if err != nil {
		service.ErrorLogger.Fatalf("Ошибка инициализации хранилища: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pgStorage.InitSchema(ctx); err != nil {
		service.ErrorLogger.Fatalf("Ошибка создания таблицы client_limits: %v", err)
	}

	rl, err := service.NewRateLimiterWithStorage(
		cfg.RateLimit.DefaultCapacity,
		cfg.RateLimit.DefaultRate,
		pgStorage,
	)
	if err != nil {
		service.ErrorLogger.Fatalf("Ошибка инициализации rate limiter: %v", err)
	}
	defer rl.Stop()

	balancer, err := service.NewSimpleBalancer(cfg.Backends)
	if err != nil {
		service.ErrorLogger.Fatalf("Ошибка инициализации балансировщика: %v", err)
	}

	service.StartBackends()
	go service.LiveCheck()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		balancer.Proxy.ServeHTTP(w, r)
	})

	apiHandler := API.NewRateLimitHandler(rl)
	mux.HandleFunc("/api/set", apiHandler.SetLimit)
	mux.HandleFunc("/api/get", apiHandler.GetLimits)

	handlerChain := middlewares.LoggerMiddleware(
		middlewares.RateLimitMiddleware(rl, mux),
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPServer.Port),
		Handler: handlerChain,
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		service.AppLogger.Printf("Сервер запущен на порту :%d", cfg.HTTPServer.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			service.ErrorLogger.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	<-shutdownChan
	service.AppLogger.Println("Получен сигнал завершения")

	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		service.ErrorLogger.Printf("Ошибка при завершении работы сервера: %v", err)
	}

	service.AppLogger.Println("Сервер успешно остановлен")
}
