package gateways

import "fptr/internal/entities"

type Gateway struct {
	Listener
}

type Listener interface {
	Listen() error
	MakeRequest(url string, method string, structure interface{}, data ...interface{}) error
	Authorization(config entities.AppConfig) error
}
