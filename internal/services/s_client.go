package services

import "fptr/internal/gateways"

type ClientService struct {
	gateways.Gateway
}

func (s *ClientService) Listen() {
	//err := s.Gateway.Listen()

}
