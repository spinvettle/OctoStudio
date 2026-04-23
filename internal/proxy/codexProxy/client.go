package codexProxy

import (
	"net/http"
	"time"
)

var RelayClient *http.Client
var FetchClient *http.Client

func InitHttpClient() {
	RelayClient = &http.Client{
		// Timeout: time.Minute * 1,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
		}}

	FetchClient = &http.Client{}

}
