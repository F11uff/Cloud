package service

import (
	"net/http"
	"sync/atomic"
)

var (
	regCount uint32
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Проверяем, не отменён ли уже контекст
	select {
	case <-r.Context().Done():
		ErrorLogger.Printf("Запрос отменён: %v", r.Context().Err())
		return
	default:
	}

	atomic.AddUint32(&regCount, 1)
	AppLogger.Printf("!!!ERROR APP: %d: %v", regCount, err)

	if _, ok := w.(http.CloseNotifier); ok {
		if cn, ok := w.(http.CloseNotifier); ok {
			select {
			case <-cn.CloseNotify():
				ErrorLogger.Println("Соединение закрыто клиентом")
				return
			default:
			}
		}
	}

	w.WriteHeader(http.StatusBadGateway)
	if _, err := w.Write([]byte("502 Bad Gateway")); err != nil {
		ErrorLogger.Printf("Ошибка при записи ответа: %v", err)
	}
}
