package models

import (
	"cloud/Balancer/internal/core"
	//"cloud/Balancer/internal/service"
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

func (sl *SimpleBalancer) NextBackend() *url.URL {
	n := atomic.AddUint64(&sl.Counter, 1)
	return core.LiveBackends[(n-1)%uint64(len(core.LiveBackends))]
}

func (sl *SimpleBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sl.Proxy.ServeHTTP(w, r)
}
