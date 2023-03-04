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
	GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error)
	PrintSell(info entities.Info, id string) error
	PrintRefound(info entities.Info, id string) error
	Login(config entities.AppConfig) (*entities.SessionInfo, error)
	PrintRefoundFromSell(info entities.Info, id string) error
}

type KKT interface {
	PrintXReport() error
	MakeSession(info entities.Info) error
	CloseShift() error
	ShiftIsOpened() bool
	ShiftIsClosed() bool
	ShiftIsExpired() bool
	CurrentShiftStatus() uint
	CashIncome(income float64) error
	CurrentError() error
	PrintLastCheckPressedFromKKT() error
	Beep(beepType string)
}
