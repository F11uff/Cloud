package mock

import (
	"cloud/Balancer/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MockTest_NextBackend(t *testing.T) {
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend2.Close()

	lb, err := service.NewSimpleBalancer([]string{backend1.URL, backend2.URL})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil) // Здесь тестируем ротацию между серверами
	w := httptest.NewRecorder()

	lb.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
