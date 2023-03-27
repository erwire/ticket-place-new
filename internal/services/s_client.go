package services

import (
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	"fptr/internal/gateways"
	"github.com/google/logger"
	"net/http"
	"time"
)

const AttemptDurationInSeconds = 2
const AttemptCount = 5

type ClientService struct {
	gw *gateways.Gateway
	*logger.Logger
}

func NewClientService(gw *gateways.Gateway, logg *logger.Logger) *ClientService {
	return &ClientService{
		gw:     gw,
		Logger: logg,
	}
}

func (s *ClientService) SetTimeout(timeout time.Duration) {
	s.gw.SetTimeout(timeout)
	s.Infof("Установлена длительность попытки запроса на %s", timeout)
}

func (s *ClientService) GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error) {

	var err error
	var click *entities.Click

	for i := 0; i < AttemptCount; i++ {
		click, err = s.gw.GetLastReceipt(connectionURL, session)
		if err != nil {
			if i == 0 {
				s.Errorf("Во время запроса к истории печати произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
			}
			return click, nil
		}
	}

	switch err.(type) {
	case *apperr.ClientError:
		if !(err.(*apperr.ClientError).StatusCode == http.StatusNotFound) {
			s.Errorf("[Попытки завершились неудачно] -> Ошибка при запросе последнего заказа: %v", err)
		}
	}
	return nil, err

}

func (s *ClientService) PrintSell(info entities.Info, id string, uuid *string) error {
	uuidStr := ""

	if uuid == nil {
		uuidStr = "отсутствует (печать с панели)"
	} else {
		uuidStr = *uuid
	}

	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		s.Errorf("Ошибка во время получения заказа с номером %s, uuid: %s, клиент: %v", id, uuidStr, err)
		return err
	}

	if err = s.gw.KKT.PrintSell(*sell); err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати чека продажи заказа с номером %s, uuid: %s, ККТ: %v", id, uuidStr, err)
		default:
			s.Errorf("Ошибка во время печати чека продажи заказа с номером %s, uuid: %s, ККТ: %v", id, uuidStr, err)
		}
		return err
	}
	s.Infof("Выполнена печать чека заказа с номером: %s, uuid: %s\n", id, uuidStr)
	return nil
}

func (s *ClientService) PrintRefoundFromSell(info entities.Info, id string) error {
	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		s.Errorf("Ошибка во время печати возврата заказа с номером %s, клиент: %v", id, err)
		return err
	}
	err = s.gw.KKT.PrintRefoundFromCheck(*sell)
	if err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати возврата заказа с номером %s, ККТ: %v", id, err)
		default:
			s.Errorf("Ошибка во время печати возврата заказа с номером %s, ККТ: %v", id, err)
		}
		return err
	}
	s.Infof("Выполнена печать чека возврата заказа с номером: %s\n", id)
	return nil
}

func (s *ClientService) PrintRefound(info entities.Info, id string, uuid *string) error {
	uuidStr := ""

	if uuid == nil {
		uuidStr = "отсутствует (печать с панели)"
	} else {
		uuidStr = *uuid
	}

	refound, err := s.gw.Listener.GetRefound(info, id)

	if err != nil {
		s.Errorf("Ошибка во время получения возврата заказа с номером %s, uuid: %s, клиент: %v", id, uuidStr, err)
		return err
	}

	err = s.gw.KKT.PrintRefound(*refound)
	if err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати возврата заказа с номером %s, uuid: %s, ККТ: %v", id, uuidStr, err)
		default:
			s.Errorf("Ошибка во время печати возврата заказа с номером %s, uuid: %s, ККТ: %v", id, uuidStr, err)
		}
		return err
	}

	s.Infof("Выполнена печать чека возврата заказа с номером: %s, uuid: %s\n", id, uuid)
	return nil
}

func (s *ClientService) Login(config entities.AppConfig) (*entities.SessionInfo, error) {
	var err error
	var session *entities.SessionInfo
	for i := 0; i < AttemptCount; i++ {
		session, err = s.gw.Login(config)
		if err != nil {
			if i == 0 {
				s.Errorf("Во время авторизации произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Во время авторизации произошла ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
			}
			return session, nil
		}
	}
	s.Errorf("[Попытки завершились неудачно] -> Во время авторизации произошла ошибка: %v", err)
	return nil, err
}
