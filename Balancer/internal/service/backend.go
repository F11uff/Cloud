package service

import (
	"cloud/Balancer/internal/core"
	"context"
	"net/http"
	"net/url"
	"time"
)

func LiveCheck() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var newHealthy []*url.URL

		for _, backend := range core.LiveBackends {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			req, err := http.NewRequestWithContext(ctx, "GET", backend.String()+"/health", nil)
			if err != nil {
				AppLogger.Printf("Ошибка создания запроса к %s: %v", backend.String(), err)
				cancel()
				continue
			}

			resp, err := http.DefaultClient.Do(req) //отправляем и принимаем http запрос
			cancel()

			if err != nil || resp.StatusCode != http.StatusOK {
				AppLogger.Printf("Бэкенд %s недоступен: %v", backend.String(), err)
			} else {
				newHealthy = append(newHealthy, backend)
				AppLogger.Printf("Бэкенд %s здоров", backend.String())
			}

			if resp != nil {
				resp.Body.Close()
			}
		}

		core.LiveBackends = newHealthy
		AppLogger.Printf("Доступные бекенды: %v", core.LiveBackends)
	}
}
