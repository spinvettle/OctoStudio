package adapter

import (
	"net/http"

	"github.com/spinvettle/OctoStudio/internal/gateway/channel"
)

type CodexAdopter struct {
}

func (a *CodexAdopter) DoRequest(body []byte, channelKey channel.ChannelKey, baseUrl string) (*http.Response, error) {
	return nil, nil
}

func (a *CodexAdopter) DoResponse() error {
	return nil
}
