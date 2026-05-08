package adapter

type Adapter interface {
	DoRequest() error
	DoResponse() error
}
