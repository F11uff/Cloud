package models

import (
	"cloud/Balancer/internal/core"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type SimpleBalancer struct {
	Proxy    *httputil.ReverseProxy
	Backends []*url.URL
	Counter  uint64
}

func (sb *SimpleBalancer) NextBackend() *url.URL {
	if len(sb.Backends) == 0 {
		return nil
	}

	n := atomic.AddUint64(&sb.Counter, 1)
	return core.LiveBackends[(n-1)%uint64(len(core.LiveBackends))]
}

func (sb *SimpleBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sb.Proxy.ServeHTTP(w, r)
}

func (sb *SimpleBalancer) GetBackends() []*url.URL {
	return sb.Backends
}

func (sb *SimpleBalancer) SetDirector(director func(*http.Request)) {
	sb.Proxy.Director = director
}
