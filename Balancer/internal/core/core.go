package core

import (
	"net/http"
	"net/url"
)

var (
	LiveBackends []*url.URL
)

type Balancer interface {
	NextBackend() *url.URL
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
