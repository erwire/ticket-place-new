package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
)

type Services struct {
	Listener
	KKT
	*LoggerService
}

func NewServices(g *gateways.Gateway, logger *LoggerService) *Services {
	return &Services{
		Listener:      NewClientService(g, logger.Logger),
		KKT:           NewKKTService(g, logger.Logger),
		LoggerService: logger,
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
	ShiftIsOpened() bool
	ShiftIsClosed() bool
	ShiftIsExpired() bool
	CurrentShiftStatus() uint
}
