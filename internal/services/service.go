package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
	"fyne.io/fyne/v2/data/binding"
	"time"
)

type Services struct {
	Listener
	KKT
	PrinterInterface
	*LoggerService
	DS *DatabaseService
}

func NewServices(g *gateways.Gateway, logger *LoggerService) *Services {
	ds := NewDatabaseService(logger.Logger, g)
	return &Services{
		Listener:         NewClientService(g, logger.Logger, ds),
		KKT:              NewKKTService(g, logger.Logger),
		PrinterInterface: NewPrinterService(logger.Logger),
		LoggerService:    logger,
		DS:               ds,
	}
}

type Listener interface {
	GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error)
	PrintSell(info entities.Info, id string, uuid *string, dto entities.PageParamsDTO, box entities.PrintCheckBox) error
	PrintRefound(info entities.Info, id string, uuid *string) error
	Login(config entities.AppConfig) (*entities.SessionInfo, error)
	PrintRefoundFromSell(info entities.Info, id string) error
	SetTimeout(timeout time.Duration)
	SetProgressData(pc binding.Float, st binding.String)
}

type KKT interface {
	Destroy()
	Configurate() error
	PrintXReport() error
	MakeSession(fullName string) error
	CloseShift() error
	ShiftIsOpened() bool
	ShiftIsClosed() bool
	ShiftIsExpired() bool
	CurrentShiftStatus() uint
	CashIncome(income float64) error
	CurrentError() error
	PrintLastCheckPressedFromKKT() error
	Beep(beepType string)
	Open() error
}
