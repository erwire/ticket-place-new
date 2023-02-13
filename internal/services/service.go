package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
)

type Services struct {
	Listener
}

func NewServices(g *gateways.Gateway) *Services {
	return &Services{
		Listener: NewClientService(g),
	}
}

type Listener interface {
	GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, string)
	Login(config entities.AppConfig) (*entities.SessionInfo, string)
}
