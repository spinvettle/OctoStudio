package proxy

import (
	"bufio"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/utils"
)

func Relay(c *gin.Context) {

	originBodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.FAIL(c, http.StatusBadRequest, "json format error", nil)
		return
	}
	resp, err := proxyService.DoProxyRequest(originBodyBytes, c.Request.Header)
	if err != nil {
		utils.FAIL(c, http.StatusServiceUnavailable, "server proxy error:"+err.Error(), nil)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
		}
	}

	c.Writer.WriteHeader(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		doSSEResponse(c, resp)
	} else {
		doErrResponse(c, resp)
	}

}

func doSSEResponse(c *gin.Context, resp *http.Response) {

	scanner := bufio.NewScanner(resp.Body)
	// reader := bufio.NewReader(resp.Body)
	for scanner.Scan() {
		select {
		case <-c.Request.Context().Done():
			return
		default:
			line := scanner.Text()
			if line != "" {
				line := scanner.Text()
				_, _ = c.Writer.Write([]byte(line + "\n"))
				c.Writer.Flush()

			}

		}

	}
	if err := scanner.Err(); err != nil {
		log.Panicln("sse parse error:" + err.Error())
	}
}

func doErrResponse(c *gin.Context, resp *http.Response) {
	upstreamBody, _ := io.ReadAll(resp.Body)

	_, _ = c.Writer.Write(upstreamBody)
	// io.Copy(c.Writer, resp.Body)
}
