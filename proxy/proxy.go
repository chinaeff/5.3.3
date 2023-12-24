package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxyConfig struct {
	BackendURLs []string
}

func NewReverseProxy(config ReverseProxyConfig) http.Handler {
	director := func(req *http.Request) {
		targetURL := config.BackendURLs[0]
		target, err := url.Parse(targetURL)
		if err != nil {
			log.Fatal("Error parsing target URL:", err)
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}

	return &httputil.ReverseProxy{Director: director}
}
