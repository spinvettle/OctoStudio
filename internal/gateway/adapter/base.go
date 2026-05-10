package adapter

import (
	"net/http"

	"github.com/spinvettle/OctoStudio/internal/gateway"
	"github.com/spinvettle/OctoStudio/internal/gateway/channel"
)

type Adapter interface {
	DoRequest(body []byte, channelKey channel.ChannelKey, baseUrl string) (*http.Response, error)
	DoResponse() error
}

func GetAdapterByRelayMode(mode gateway.RelayMode) Adapter {
	switch mode {
	case gateway.CodexResponse:
		return &CodexAdopter{}
	}
	return nil
}
