package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
	"github.com/google/logger"
)

type KKTService struct {
	gw *gateways.Gateway
	*logger.Logger
}

func NewKKTService(gw *gateways.Gateway, logg *logger.Logger) *KKTService {
	return &KKTService{
		gw:     gw,
		Logger: logg,
	}
}

func (s *KKTService) MakeSession(info entities.Info) string {
	if err := s.gw.Open(); err != nil {
		s.Errorf("Ошибка при установлении связи с кассой: %v\n", err)
		return "Ошибка при установлении связи с кассой"
	}

	if err := s.gw.NewCashierRegister(info.Session); err != nil {
		s.Errorf("Ошибка при регистрации кассира: %v\n", err)
		return "Ошибка при регистрации кассира"
	}

	if err := s.gw.OpenShift(); err != nil {
		if s.gw.ShiftIsExpired() {
			err = s.gw.CloseShift()
			if err != nil {
				s.Errorf("Критическая ошибка работы приложения: %v\n", err)
				return "Критическая ошибка работы приложения"
			}
			err = s.gw.OpenShift()
			if err != nil {
				s.Errorf("Критическая ошибка работы приложения: %v\n", err)
				return "Критическая ошибка работы приложения"
			}
		}
	}
	return ""
}

func (s *KKTService) CloseShift() string {
	if err := s.gw.KKT.CloseShift(); err != nil {
		s.Errorf("Ошибка закрытия смены: %v\n", err)
		return "Ошибка закрытия смены"
	}
	return ""
}

func (s *KKTService) ShiftIsOpened() bool {
	return s.gw.ShiftIsOpened()
}
func (s *KKTService) ShiftIsClosed() bool {
	return s.gw.ShiftIsClosed()
}

func (s *KKTService) ShiftIsExpired() bool {
	return s.gw.ShiftIsExpired()
}

func (s *KKTService) CurrentShiftStatus() uint {
	return s.gw.KKT.CurrentShiftStatus()
}
