package services

import (
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	"fptr/internal/gateways"
	"fyne.io/fyne/v2/data/binding"
	"github.com/google/logger"
	"net/http"
	"time"
)

const AttemptDurationInSeconds = 2
const AttemptCount = 5

type ClientService struct {
	gw             *gateways.Gateway
	ProgressCount  binding.Float //+ Внедрение Progress Bar
	ProgressStatus binding.String
	*logger.Logger
}

func NewClientService(gw *gateways.Gateway, logg *logger.Logger) *ClientService {
	return &ClientService{
		gw:     gw,
		Logger: logg,
	}
}

func (s *ClientService) SetProgressData(pc binding.Float, st binding.String) {
	s.ProgressCount = pc
	s.ProgressStatus = st
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

			switch err.(type) {
			case *apperr.ClientError:
				if err.(*apperr.ClientError).StatusCode == http.StatusNotFound {
					return nil, err
				}
			}

			s.ProgressStatus.Set("Текущий статус: во время обращения к истории пользователя возникла ошибка")
			if i == 0 {
				s.Errorf("Во время запроса к истории печати произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			s.ProgressCount.Set(float64(i+1) * 1 / float64(AttemptCount))
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
				s.ProgressCount.Set(1)
				s.ProgressStatus.Set("Текущий статус: операция выполнена успешно")
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

	var sell *entities.Sell
	var err error

	for i := 0; i < AttemptCount; i++ {
		sell, err = s.gw.Listener.GetSell(info, id)
		if err != nil {
			s.ProgressStatus.Set("Текущий статус: во время попытки получить от сервера данные возникла ошибка")
			if i == 0 {
				s.Errorf("Во время запроса к данным заказа произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			s.ProgressCount.Set(float64(i+1) * 1 / float64(AttemptCount))
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
				s.ProgressCount.Set(1)
				s.ProgressStatus.Set("Текущий статус: операция выполнена успешно")
			}
			break
		}
	}

	if err != nil {
		s.Errorf("[Попытки закончены неудачно] -> ID: %s, UUID: %s, ERR: %v", id, uuidStr, err)
		return err
	}

	if err = s.gw.KKT.PrintSell(*sell, info.AppConfig.User.TaxesInfo); err != nil {
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
	var sell *entities.Sell
	var err error
	for i := 0; i < AttemptCount; i++ {
		sell, err = s.gw.Listener.GetSell(info, id)
		if err != nil {
			s.ProgressStatus.Set("Текущий статус: во время попытки получить от сервера данные возникла ошибка")
			if i == 0 {
				s.Errorf("Во время запроса к данным заказа произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			s.ProgressCount.Set(float64(i+1) * 1 / float64(AttemptCount))
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
				s.ProgressCount.Set(1)
				s.ProgressStatus.Set("Текущий статус: операция выполнена успешно")
			}
			break
		}
	}

	if err != nil {
		s.Errorf("Ошибка во время печати возврата заказа с номером %s, клиент: %v", id, err)
		return err
	}

	err = s.gw.KKT.PrintRefoundFromCheck(*sell, info.AppConfig.User.TaxesInfo)
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
	var refound *entities.Refound
	var err error

	if uuid == nil {
		uuidStr = "отсутствует (печать с панели)"
	} else {
		uuidStr = *uuid
	}

	for i := 0; i < AttemptCount; i++ {
		refound, err = s.gw.Listener.GetRefound(info, id)
		if err != nil {
			s.ProgressStatus.Set("Текущий статус: во время попытки получить от сервера данные возникла ошибка")
			if i == 0 {
				s.Errorf("Во время запроса к данным заказа произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			s.ProgressCount.Set(float64(i+1) * 1 / float64(AttemptCount))
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
				s.ProgressCount.Set(1)
				s.ProgressStatus.Set("Текущий статус: операция выполнена успешно")
			}
			break
		}
	}

	if err != nil {
		s.Errorf("Ошибка во время получения возврата заказа с номером %s, uuid: %s, клиент: %v", id, uuidStr, err)
		return err
	}

	err = s.gw.KKT.PrintRefound(*refound, info.AppConfig.User.TaxesInfo)
	if err != nil {
		switch err.(type) {
		case *apperr.BusinessError:
			s.Warningf("Ошибка во время печати возврата заказа с номером %s, uuid: %s, ККТ: %v", id, uuidStr, err)
		default:
			s.Errorf("Ошибка во время печати возврата заказа с номером %s, uuid: %s, ККТ: %v", id, uuidStr, err)
		}
		return err
	}

	s.Infof("Выполнена печать чека возврата заказа с номером: %s, uuid: %s\n", id, uuidStr)
	return nil
}

func (s *ClientService) Login(config entities.AppConfig) (*entities.SessionInfo, error) {
	var err error
	var session *entities.SessionInfo

	for i := 0; i < AttemptCount; i++ {
		session, err = s.gw.Login(config)
		if err != nil {
			s.ProgressStatus.Set("Текущий статус: во время попытки получить от сервера данные возникла ошибка")
			if i == 0 {
				s.Errorf("Во время авторизации произошла ошибка: %v", err)
				s.Warningf("[Операция закончилась неуспешно, запуск повторных попыток подключения]")
			} else {
				s.Errorf("[Попытка номер %d] -> Ошибка: %v", i, err)
			}
			time.Sleep(AttemptDurationInSeconds * time.Second)
			s.ProgressCount.Set(float64(i+1) * 1 / float64(AttemptCount))
			continue
		} else {
			if i != 0 {
				s.Infof("[Попытка номер %d] -> Попытка завершена успешно, операция выполнена", i)
				s.ProgressCount.Set(1)
				s.ProgressStatus.Set("Текущий статус: операция выполнена успешно")
			}
			break
		}
	}
	if err != nil {
		s.Errorf("[Попытки завершились неудачно] -> Во время авторизации произошла ошибка: %v", err)
		return nil, err
	}

	return session, err
}
