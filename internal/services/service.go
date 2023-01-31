package services

type Service struct {
	Auth
	Listener
}

type Auth interface {
}

type Listener interface {
	Listen()
	MakeRequest(url string, method string, structure interface{}, data ...interface{}) error
}
