package gateways

import (
	"fptr/internal/entities"
	fptr10 "github.com/EuginKostomarov/ftpr10"
	"net/http"
)

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
	GetSell(connectionURL string, sellID string) (*entities.Sell, error)
	GetRefound(connectionURL string, refoundID string) (*entities.Refound, error)
}

type KKT interface {
	Open() error
}
