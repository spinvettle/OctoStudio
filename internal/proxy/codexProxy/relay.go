package codexProxy

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/config"
	"github.com/spinvettle/OctoStudio/internal/utils"
)

func CodexRelay(c *gin.Context) {

	originBodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to read request body", slog.Any("error", err))
		utils.FAIL(c, http.StatusBadRequest, "failed to read request body"+err.Error(), nil)
		return
	}
	ctx1 := context.WithoutCancel(c.Request.Context()) //只传Value
	ctx, cancel := context.WithTimeout(ctx1, time.Duration(config.CodexRelayTimeOut)*time.Second)
	defer cancel()
	resp, err := ProxyService.DoProxyRequest(ctx, originBodyBytes, c.Request.Header)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "upstream failed relay", slog.Any("error", err))
		utils.FAIL(c, http.StatusServiceUnavailable, "upstream failed relay"+err.Error(), nil)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
		}
	}

	c.Writer.WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to response client", slog.Any("error", err))
		return
	}
	// if resp.StatusCode == http.StatusOK {
	// 	doSSEResponse(c, resp)
	// } else {
	// 	doErrResponse(c, resp)
	// }

}

// func doSSEResponse(c *gin.Context, resp *http.Response) {

// 	reader := bufio.NewReader(resp.Body)
// 	// reader := bufio.NewReader(resp.Body)
// 	for {
// 		select {
// 		case <-c.Request.Context().Done():
// 			return
// 		default:
// 			line, readeErr := reader.ReadString('\n')
// 			if line != "" {
// 				_, err := c.Writer.Write([]byte(line + "\n"))
// 				fmt.Println(line)
// 				c.Writer.Flush()
// 				if err != nil {
// 					log.Println("sse client write error:", err)
// 					return
// 				}

// 			}
// 			if readeErr != nil {
// 				if readeErr == io.EOF {
// 					return
// 				}
// 				log.Println("sse parse error:" + readeErr.Error())
// 				return
// 			}

// 		}

// 	}

// }

// func doErrResponse(c *gin.Context, resp *http.Response) {
// 	upstreamBody, _ := io.ReadAll(resp.Body)

// 	_, _ = c.Writer.Write(upstreamBody)
// 	// io.Copy(c.Writer, resp.Body)
// }
