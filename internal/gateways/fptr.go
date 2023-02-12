package gateways

import fptr10 "github.com/EuginKostomarov/ftpr10"

type KKTGateway struct {
	IFptr *fptr10.IFptr
}

func NewKKTGateway(IFptr *fptr10.IFptr) *KKTGateway {
	return &KKTGateway{IFptr: IFptr}
}

func (g *KKTGateway) Open() error {
	return nil
}
