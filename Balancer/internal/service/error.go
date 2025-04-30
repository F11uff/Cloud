package service

import (
	"net/http"
	"sync/atomic"
)

var (
	regCount uint32
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	atomic.AddUint32(&regCount, 1) // увеличиваем атомарно счетчик ошибок
	AppLogger.Printf("!!!ERROR APP: %d: %v", regCount, err)
	w.WriteHeader(http.StatusBadGateway)

	if _, err := w.Write([]byte("502 Bad Gateway")); err != nil {
		ErrorLogger.Printf("Ошибка при записи ответа: %v", err)
	}
}
