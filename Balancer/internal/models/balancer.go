package models

import (
	"cloud/Balancer/internal/service"
	"net/http/httputil"
	"net/url"
)

type Balancer interface {
}

type SimpleBalancer struct {
	proxy    *httputil.ReverseProxy
	backends []*url.URL
	counter  uint32
}

func NewSimpleBalancer(URLS []string) *SimpleBalancer {
	var sliceURL []*url.URL

	for _, urll := range URLS {
		urlParse, err := url.Parse(urll)
		if err != nil {
			service.ErrorLogger.Printf("Failed to parse url: %s", urlParse)
		}

		sliceURL = append(sliceURL, urlParse)
	}

	return &SimpleBalancer{
		proxy:    nil, // нужно работать с http
		backends: sliceURL,
		counter:  0,
	}
}
