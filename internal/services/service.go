package services

type Service struct {
	Auth
	Listener
}

type Auth interface {
}

type Listener interface {
	Listen() error
	MakeRequest(url string, method string, structure interface{}, data ...interface{}) error
}
