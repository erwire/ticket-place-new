package services

import (
	"errors"
	"fmt"
	"fptr/internal/entities"
	"fptr/internal/gateways"
	errorlog "fptr/pkg/error_logs"
	"github.com/google/logger"
	"log"
)

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

func (s *ClientService) GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, string) {
	click, err := s.gw.GetLastReceipt(connectionURL, session)
	if err != nil {
		switch errors.Unwrap(err) {
		default:
			s.Errorf("Произошла ошибка получения последнего запроса на печать: %v\n", err)
			return nil, "Произошла ошибка получения последнего запроса на печать"
		}
	}

	return click, ""
}

func (s *ClientService) PrintSell(info entities.Info, id string) string {
	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		s.Errorf("Ошибка во время выполнения запроса: %v\n", err)
		return "Ошибка во время выполнения запроса"
	}
	err = s.gw.KKT.PrintSell(*sell)
	if err != nil {
		switch errors.Unwrap(err) {
		case errorlog.ShiftIsExpired:
			s.gw.CloseShift()
			return "Смена истекла. Пожалуйста, переавторизуйтесь."
		}

		errorMessage := fmt.Sprintf("Ошибка во время печати заказа с номером %s", id)
		s.Errorf("%s: %v\n", errorMessage, err)

		return errorMessage
	}
	s.Infof("Выполнена печать чека заказа с номером: %s\n", id)
	return ""
}

func (s *ClientService) PrintRefoundFromSell(info entities.Info, id string) string {
	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		s.Errorf("Ошибка во время выполнения запроса: %v\n", err)
		return "Ошибка во время выполнения запроса"
	}
	err = s.gw.KKT.PrintRefoundFromCheck(*sell)
	if err != nil {
		switch errors.Unwrap(err) {
		case errorlog.ShiftIsExpired:
			s.gw.CloseShift()
			return "Смена истекла. Пожалуйста, переавторизуйтесь."
		}

		errorMessage := fmt.Sprintf("Ошибка во время печати возврата с номером %s", id)
		s.Errorf("%s: %v\n", errorMessage, err)

		return errorMessage
	}
	s.Infof("Выполнена печать чека возврата заказа с номером: %s\n", id)
	return ""
}

func (s *ClientService) PrintRefound(info entities.Info, id string) string {
	refound, err := s.gw.Listener.GetRefound(info, id)
	if err != nil {
		s.Errorf("Ошибка во время выполнения запроса: %v\n", err)
		return "Ошибка во время выполнения запроса"
	}
	err = s.gw.KKT.PrintRefound(*refound)
	if err != nil {
		switch errors.Unwrap(err) {
		case errorlog.ShiftIsExpired:
			s.gw.CloseShift()
			return "Смена истекла. Пожалуйста, переавторизуйтесь."
		}
		errorMessage := fmt.Sprintf("Ошибка во время печати заказа с номером %s", id)
		s.Errorf("%s: %v\n", errorMessage, err)
		return errorMessage
	}
	s.Infof("Выполнена печать чека возврата заказа с номером: %s\n", id)
	return ""
}

func (s *ClientService) Login(config entities.AppConfig) (*entities.SessionInfo, string) {
	session, err := s.gw.Login(config)
	if err != nil {
		log.Println(err.Error())
		switch errors.Unwrap(err) {
		case errorlog.EmptyURLDataError:
			s.Errorf("Нет данных по хосту. Пожалуйста, добавьте адрес в настройках: %v\n", err)
			return nil, "Нет данных по хосту. Пожалуйста, добавьте адрес в настройках"
		case errorlog.InvalidLoginOrPassword:
			s.Errorf("Логин или пароль не заполнены или заполнены некорректно: %v\n", err)
			return nil, "Логин или пароль не заполнены или заполнены некорректно"
		case errorlog.AuthorizationError:
			s.Errorf("Произошла ошибка авторизации: %v\n", err)
			return nil, "Произошла ошибка авторизации"
		case errorlog.JsonUnmarshalError:
			s.Errorf("Неправильный логин или пароль: %v\n", err)
			return nil, "Неправильный логин или пароль"
		default:
			s.Errorf("Произошла непредвиденная ошибка: %v\n", err)
			return nil, "Произошла непредвиденная ошибка"
		}
	}
	return session, ""
}
