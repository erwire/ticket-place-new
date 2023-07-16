package gateways

import (
	"fptr/internal/entities"
	"fptr/pkg/fptr10"
	"net/http"
	"time"
)

type PrintType string

type Gateway struct {
	Listener
	KKT
	PrinterInterface
}

func NewGateway(client *http.Client, iFptr *fptr10.IFptr) *Gateway {
	return &Gateway{
		Listener:         NewClientGateway(client),
		KKT:              NewKKTGateway(iFptr),
		PrinterInterface: NewPrinter(),
	}
}

type Listener interface {
	Login(config entities.AppConfig) (*entities.SessionInfo, error)
	GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error)
	GetSell(info entities.Info, sellID string) (*entities.Sell, error)
	GetRefound(info entities.Info, refoundID string) (*entities.Refound, error)
	SetTimeout(timeout time.Duration)
}

type KKT interface {
	Open() error
	Close() error
	Configurate() error
	Destroy()
	OpenShift() error
	CloseShift() error
	PrintSell(sell entities.Sell) error
	PrintRefound(refound entities.Refound) error
	NewCashierRegister(fullName string) error
	ShiftIsExpired() bool
	ShiftIsOpened() bool
	ShiftIsClosed() bool
	CurrentShiftStatus() uint
	PrintXReport() error
	CashIncome(income float64) error
	CurrentErrorStatusCode() error
	PrintRefoundFromCheck(sell entities.Sell) error
	PrintLastCheckPressedFromKKT() error
	WarningBeep()
	ErrorBeep()
}
