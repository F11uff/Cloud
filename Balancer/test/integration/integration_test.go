package integration_test

import (
	"cloud/Balancer/internal/models"
	"cloud/Balancer/internal/service"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"testing"
)

type BalancerService interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetBackends() []*url.URL
	SetDirector(func(*http.Request))
}

func TestIntegration(t *testing.T) {
	service.AppLogger = log.New(os.Stdout, "[TEST] ", log.LstdFlags)
	service.ErrorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)

	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // мокирование
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer mockBackend.Close()

	backendURL, _ := url.Parse(mockBackend.URL) // Инициализация балансировщика через интерфейс
	var balancer BalancerService = &models.SimpleBalancer{
		Backends: []*url.URL{backendURL},
		Proxy:    httputil.NewSingleHostReverseProxy(backendURL),
	}

	balancer.SetDirector(func(req *http.Request) { // Настройка director через интерфейс
		req.URL.Scheme = backendURL.Scheme
		req.URL.Host = backendURL.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = backendURL.Host
	})

	healthReq := httptest.NewRequest("GET", "/health", nil)
	healthRec := httptest.NewRecorder()
	balancer.ServeHTTP(healthRec, healthReq)

	if healthRec.Code != http.StatusOK {
		t.Fatalf("Health check failed: %d", healthRec.Code)
	}

	testReq := httptest.NewRequest("GET", "/test", nil) // запрос
	testRec := httptest.NewRecorder()

	ctx := context.WithValue(testReq.Context(), "stepsl", balancer)
	testReq = testReq.WithContext(ctx)

	balancer.ServeHTTP(testRec, testReq)

	if testRec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", testRec.Code)
	}
	if testRec.Body.String() != "OK" {
		t.Errorf("Unexpected response body: %s", testRec.Body.String())
	}

	backends := balancer.GetBackends()
	if len(backends) != 1 {
		t.Errorf("Expected 1 backend, got %d", len(backends))
	}
}
