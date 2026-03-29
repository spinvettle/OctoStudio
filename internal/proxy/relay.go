package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/config"
	"github.com/spinvettle/OctoStudio/internal/utils"
)

const CodexResponsesURL = "https://chatgpt.com/backend-api/codex/responses"

type RelayHandler struct {
	client *http.Client
	pool   *AccountPool
}

func NewRelayHandler() *RelayHandler {
	client := &http.Client{Timeout: time.Second * 120,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
		}}
	pool := &AccountPool{}

	return &RelayHandler{
		client: client,
		pool:   pool,
	}

}

func (r *RelayHandler) Relay(c *gin.Context) {

	originBodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.ERROR(c, http.StatusBadRequest, "body格式错误")
	}

	if err != nil {
		utils.ERROR(c, http.StatusInternalServerError, "代理请求构造失败:"+err.Error())
	}

	account := r.pool.GetAccount()
	var resp *http.Response
	for range config.GlobalConfig.RetryCount {
		proxyReq, err := http.NewRequest("POST", CodexResponsesURL, bytes.NewBuffer(originBodyBytes))
		if err != nil {
			utils.ERROR(c, http.StatusInternalServerError, "代理失败:"+err.Error())
		}
		for k, v := range c.Request.Header {
			proxyReq.Header[k] = v
		}
		proxyReq.Header["Authorization"] = []string{"Bearer " + account.AccessToken}

		resp, err = r.client.Do(proxyReq)
		if err != nil {
			utils.ERROR(c, http.StatusInternalServerError, "代理失败:"+err.Error())
		}

		if resp.StatusCode == http.StatusOK {
			break
		} else if resp.StatusCode == http.StatusTooManyRequests {
			account = r.pool.GetAccount()
		} else if resp.StatusCode == http.StatusUnauthorized {
			account = r.pool.GetAccount()
		}
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				log.Printf("close body error: %v", err)
			}
		}

	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("close body error: %v", err)
		}
	}()
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.WriteHeader(http.StatusOK)
	buffer := make([]byte, 1024*4)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, _ = c.Writer.Write(buffer[:n])
			c.Writer.Flush()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Stream read error: %v", err)
			break
		}

	}

}
