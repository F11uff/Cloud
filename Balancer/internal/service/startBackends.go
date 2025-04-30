package service

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func StartBackends() {
	backends := []struct {
		port     string
		path     string
		response string
	}{
		{":8081", "/health", "Backend 1 OK"},
		{":8082", "/health", "Backend 2 OK"},
		{":8083", "/health", "Backend 3 OK"},
		{":8084", "/health", "Backend 4 OK"},
		{":8085", "/health", "Backend 5 OK"},
	}

	var wg sync.WaitGroup
	wg.Add(len(backends))

	for _, b := range backends {
		go func(port, path, response string) {
			defer wg.Done()
			mux := http.NewServeMux()
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(response))
			})
			AppLogger.Printf("Тестовый бэкенд запущен на %s", port)
			log.Fatal(http.ListenAndServe(port, mux)) // Сопоставляем URL-пути с обработчиками
		}(b.port, b.path, b.response)
	}

	time.Sleep(500 * time.Millisecond)
}
