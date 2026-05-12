package httpclient

import (
	"net/http"
	"time"
)

func NewRelayClient(timeout time.Duration) *http.Client {
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxConnsPerHost:     50,
			MaxIdleConnsPerHost: 20,
			MaxIdleConns:        200,
		},
	}
	return &client
}
func NewFetchClient(timeout time.Duration) *http.Client {
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	return &client
}
