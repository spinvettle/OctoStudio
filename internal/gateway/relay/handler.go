package relay

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spinvettle/OctoStudio/internal/gateway"
	"github.com/spinvettle/OctoStudio/internal/gateway/adapter"
	"github.com/spinvettle/OctoStudio/internal/gateway/channel"
)

type RelayHandler struct {
	channelSvc *channel.ChannelService
	retryTimes int
}

func NewRelayHandler(channel *channel.ChannelService) *RelayHandler {
	return &RelayHandler{
		channelSvc: channel,
	}
}
func (h *RelayHandler) Relay(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1024*2014*10)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to read body", slog.Any("error", err))
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "request body too large limit",
		})
		return
	}
	defer func() {

		err := c.Request.Body.Close()
		if err != nil {
			slog.ErrorContext(c.Request.Context(), "client request body close failed", slog.Any("error", err))
		}

	}()

	var bodyMap map[string]any
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to unmarshal body", slog.Any("error", err))
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "failed to unmarshal body" + err.Error(),
		})
		return
	}
	modelName, err := getModelName(c, bodyMap)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to get modelName from unmarshal body", slog.Any("error", err))
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "need modelName" + err.Error(),
		})
		return
	}
	relayMode := getRelayMode(c)
	channels := h.channelSvc.RandomPickChannelByModel(modelName)
	adopter := adapter.GetAdapterByRelayMode(relayMode)

	resp, err := h.doRelayWithRetry(c, body, channels, adopter)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "gateway relay failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "StatusInternalServerError",
		})
		return
	}
	defer func() {
		if resp != nil {
			err := resp.Body.Close()
			if err != nil {
				slog.ErrorContext(c.Request.Context(), "gateway upstream response close failed", slog.Any("error", err))
			}
		}

	}()

}

func getModelName(c *gin.Context, body map[string]any) (string, error) {
	modelName, ok := body["model"].(string)
	if ok {
		return strings.ToLower(modelName), nil
	}
	modelName = c.Query("model")
	if modelName != "" {
		return strings.ToLower(modelName), nil
	}
	return "", errors.New("not found ModelName")
}

func getRelayMode(c *gin.Context) gateway.RelayMode {
	path := c.Request.URL.Path
	switch {
	case strings.HasSuffix(path, "/backend-api/codex/responses"):
		return gateway.CodexResponse
	}
	return gateway.UnKnown
}

func (h *RelayHandler) doRelayWithRetry(c *gin.Context, requestBody []byte, chs *[]channel.Channel, adapter adapter.Adapter) (*http.Response, error) {
	relayTimes := 0
	for _, ch := range *chs {
		keyIndex := 0
		for _, channelKey := range ch.Keys {
			if channelKey.Status != channel.Enable {
				continue
			}
			keyIndex += 1
			relayTimes += 1
			if relayTimes > h.retryTimes {
				return nil, errors.New("retry too many")
			}
			resp, err := adapter.DoRequest(requestBody, channelKey, ch.BaseURL)
			if err == nil && resp.StatusCode == http.StatusOK {
				return resp, nil
			}
			switch resp.StatusCode {
			case http.StatusUnauthorized, http.StatusForbidden:
				if channelKey.Metadata.CanRefresh {
					go func() {
						err := h.channelSvc.RefreshChannelKey(channelKey)
						if err != nil {
							slog.ErrorContext(c.Request.Context(), "refreshing channelKey failed",
								slog.Int("channelKeyID", channelKey.ID),
								slog.Any("error", err),
							)
						}
					}()
				}
				continue
			case http.StatusTooManyRequests:
				go func() {
					err := h.channelSvc.ColdingChannelKey(channelKey)
					if err != nil {
						slog.ErrorContext(c.Request.Context(), "colding channelKey failed",
							slog.Int("channelKeyID", channelKey.ID),
							slog.Any("error", err))
					}
				}()
				continue

			case http.StatusRequestTimeout,
				http.StatusGatewayTimeout,
				http.StatusServiceUnavailable,
				http.StatusBadGateway:
				//指数退避
				continue

			case http.StatusBadRequest,
				http.StatusNotFound,
				http.StatusUnprocessableEntity:
				return resp, nil

			default:
				return resp, err

			}

		}

	}
	return nil, errors.New("")

}
