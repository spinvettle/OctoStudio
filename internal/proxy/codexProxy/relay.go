package codexProxy

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/utils"
)

func CodexRelay(c *gin.Context) {

	originBodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.FAIL(c, http.StatusBadRequest, "json format error", nil)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	resp, err := ProxyService.DoProxyRequest(ctx, originBodyBytes, c.Request.Header)
	if err != nil {
		utils.FAIL(c, http.StatusServiceUnavailable, "server codexProxy error:"+err.Error(), nil)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
		}
	}

	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)
	// if err != nil {
	// 	slog.Error("response client err:%s" + err.Error())
	// 	return
	// }
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
