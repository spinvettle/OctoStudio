package providers

import (
	"github.com/spinvettle/OctoStudio/internal/consts"
	"github.com/spinvettle/OctoStudio/internal/schema"
)

type CodexProvider struct {
	BaseUrl string
}

func NewCodexProvider() *CodexProvider {
	return &CodexProvider{
		BaseUrl: consts.CodexResponsesURL,
	}
}

func (c *CodexProvider) DoRequest(request schema.CodexResponseRequest) {

}
