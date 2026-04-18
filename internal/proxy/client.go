package proxy

import (
	"net/http"
	"time"
)

var RelayClient *http.Client
var FetchClient *http.Client

func Init() {
	RelayClient = &http.Client{Timeout: time.Second * 120,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
		}}
}
