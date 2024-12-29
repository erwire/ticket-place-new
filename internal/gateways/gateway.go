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
}

func NewGateway(client *http.Client, iFptr *fptr10.IFptr) *Gateway {
	return &Gateway{
		Listener: NewClientGateway(client),
		KKT:      NewKKTGateway(iFptr),
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
	OpenShift() error
	CloseShift() error
	PrintSell(sell entities.Sell, taxes entities.TaxesInfo) error
	PrintRefound(refound entities.Refound, taxes entities.TaxesInfo) error
	NewCashierRegister(fullName string, inn uint64) error
	ShiftIsExpired() bool
	ShiftIsOpened() bool
	ShiftIsClosed() bool
	CurrentShiftStatus() uint
	PrintXReport() error
	CashIncome(income float64) error
	CurrentErrorStatusCode() error
	PrintRefoundFromCheck(sell entities.Sell, taxes entities.TaxesInfo) error
	PrintLastCheckPressedFromKKT() error
	WarningBeep()
	ErrorBeep()
}
