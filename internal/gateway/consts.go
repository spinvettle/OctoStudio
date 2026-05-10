package gateway

type RelayMode int

const (
	CodexResponse RelayMode = iota + 1
	Response
	ChatResponses
	Messages
	GenerateContent

	UnKnown
)
