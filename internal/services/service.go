package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
)

type Services struct {
	Listener
	KKT
}

func NewServices(g *gateways.Gateway) *Services {
	return &Services{
		Listener: NewClientService(g),
		KKT:      NewKKTService(g),
	}
}

type Listener interface {
	GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, string)
	PrintSell(info entities.Info, id string) string
	PrintRefound(info entities.Info, id string) string
	Login(config entities.AppConfig) (*entities.SessionInfo, string)
}

type KKT interface {
	MakeSession(info entities.Info) string
	CloseShift() string
}
