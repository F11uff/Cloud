package service

import (
	"cloud/Balancer/internal/core"
	"cloud/Balancer/internal/models"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

var counter uint64

func NewSimpleBalancer(URLs []string) (*models.SimpleBalancer, error) {
	var backends []*url.URL

	for _, u := range URLs {
		parsed, err := url.Parse(u)
		if err != nil {
			ErrorLogger.Printf("Ошибка парсинга URL %s: %v", u, err)
			continue
		}
		backends = append(backends, parsed)
	}

	if len(backends) == 0 {
		ErrorLogger.Println("Нет валидных URL бэкендов")
		return nil, fmt.Errorf("нет валидных бэкендов")
	}

	lb := &models.SimpleBalancer{
		Proxy: &httputil.ReverseProxy{
			Director:     director,
			ErrorHandler: ErrorHandler,
		},
		Backends: backends,
	}

	core.LiveBackends = backends

	return lb, nil
}

func director(req *http.Request) {
	if len(core.LiveBackends) == 0 {
		AppLogger.Println("Нет доступных бэкендов для перенаправления")
		return
	}

	// Round Robin алгоритм
	n := atomic.AddUint64(&counter, 1)
	backend := core.LiveBackends[n%uint64(len(core.LiveBackends))]

	req.URL.Scheme = backend.Scheme
	req.URL.Host = backend.Host
	req.Header.Set("Host", req.Host)
	req.Host = backend.Host

	AppLogger.Printf("Перенаправление %s %s → %s",
		req.Method,
		req.URL.Path,
		backend.String(),
	)
}
