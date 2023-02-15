package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
	"log"
)

type KKTService struct {
	gw *gateways.Gateway
}

func NewKKTService(gw *gateways.Gateway) *KKTService {
	return &KKTService{gw: gw}
}

func (s *KKTService) MakeSession(info entities.Info) string {
	if err := s.gw.Open(); err != nil {
		log.Println(err.Error())
		return "Ошибка при установлении связи с кассой"
	}

	if err := s.gw.NewCashierRegister(info.Session); err != nil {
		log.Println(err.Error())
		return "Ошибка при регистрации кассира"
	}

	if err := s.gw.OpenShift(); err != nil {
		if s.gw.ShiftIsExpired() {
			err = s.gw.CloseShift()
			if err != nil {
				return "Критическая ошибка работы приложения"
			}
			err = s.gw.OpenShift()
			if err != nil {
				return "Критическая ошибка работы приложения"
			}
		}
	}
	return ""
}

func (s *KKTService) CloseShift() string {
	if err := s.gw.KKT.CloseShift(); err != nil {
		log.Println(err.Error())
		return "Ошибка закрытия смены"
	}
	return ""
}
