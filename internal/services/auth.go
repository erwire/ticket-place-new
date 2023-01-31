package services

import "fptr/internal/gateways"

type AuthService struct {
	gateways.Gateway
}

func (s *AuthService) Auth() {

}
