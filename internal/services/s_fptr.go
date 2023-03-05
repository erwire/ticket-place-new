package services

import (
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

func (s *KKTService) PrintXReport() error {
	err := s.gw.PrintXReport()
	if err != nil {
		s.Errorf("%s: %v", "Ошибка при печати X-отчета", err)
		return err
	}
	return nil
}

func (s *KKTService) Open() error {
	if err := s.gw.Open(); err != nil {
		return err
	}
	return nil
}

func (s *KKTService) MakeSession(fullName string) error {
	if err := s.Open(); err != nil {
		s.Errorf("Ошибка при создании соединения с кассой: %v\n", err)
		return err
	}
	if err := s.gw.NewCashierRegister(fullName); err != nil {
		s.Errorf("Ошибка при регистрации кассира: %v\n", err)
		return err
	}

	if err := s.gw.OpenShift(); err != nil {
		if s.gw.ShiftIsExpired() {
			s.Warningf("Не можем открыть смену: %v\n", err)
			err = s.gw.CloseShift()
			if err != nil {
				s.Errorf("Попытка закрытия смены закончилось неудачей: %v\n", err)
				return err
			}
			s.Infof("Попытка закрытия смены закончилось удачно: %v\n", err)
			err = s.gw.OpenShift()
			if err != nil {
				s.Errorf("Попытка открытия смены закончилось неудачей: %v\n", err)
				return err
			}
			s.Infof("Попытка открытия смены закончилось удачно: %v\n", err)
		}
	}
	return nil
}

func (s *KKTService) CloseShift() error {
	if err := s.gw.KKT.CloseShift(); err != nil {
		s.Errorf("Ошибка закрытия смены: %v\n", err)
		return err
	}
	s.Infof("Успешное закрытие смены")
	return nil
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

func (s *KKTService) CashIncome(income float64) error {
	err := s.gw.CashIncome(income)
	if err != nil {
		s.Errorf("%v", err)
		return err
	}
	s.Infof("Успешное внесение в кассу")
	return nil
}

func (s *KKTService) CurrentError() error {
	err := s.gw.CurrentErrorStatusCode()
	if err != nil {
		return err
	}
	return nil
}

func (s *KKTService) PrintLastCheckPressedFromKKT() error {
	if err := s.gw.KKT.PrintLastCheckPressedFromKKT(); err != nil {
		s.Errorf("Ошибка при печати копии последнего чека, напечатанного в ККТ: %v", err)
		return err
	}
	s.Infof("Успешная печать копии последнего чека, напечатанного на ККТ\n")
	return nil
}

func (s *KKTService) Beep(beepType string) {
	switch beepType {
	case "warning_beep":
		s.gw.WarningBeep()
	case "error_beep":
		s.gw.ErrorBeep()
	}

}
