package services

import (
	"errors"
	"fmt"
	"fptr/internal/entities"
	"fptr/internal/gateways"
	errorlog "fptr/pkg/error_logs"
	"log"
)

type ClientService struct {
	gw *gateways.Gateway
}

func NewClientService(gw *gateways.Gateway) *ClientService {
	return &ClientService{gw: gw}
}

func (s *ClientService) GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, string) {
	click, err := s.gw.GetLastReceipt(connectionURL, session)
	if err != nil {
		log.Println(err.Error())
		switch errors.Unwrap(err) {
		default:
			return nil, "Произошла ошибка получения последнего запроса на печать"
		}
	}

	return click, ""
}

func (s *ClientService) PrintSell(info entities.Info, id string) string {
	sell, err := s.gw.Listener.GetSell(info, id)
	if err != nil {
		log.Println(err.Error())
		return "Ошибка во время выполнения запроса"
	}
	err = s.gw.KKT.PrintSell(*sell)
	if err != nil {
		log.Println(err.Error())
		return fmt.Sprintf("Ошибка во время печати заказа с номером %s", id)
	}
	return ""
}

func (s *ClientService) PrintRefound(info entities.Info, id string) string {
	refound, err := s.gw.Listener.GetRefound(info, id)
	if err != nil {
		log.Println(err.Error())
		return "Ошибка во время выполнения запроса"
	}
	err = s.gw.KKT.PrintRefound(*refound)
	if err != nil {
		log.Println(err.Error())
		return fmt.Sprintf("Ошибка во время печати возврата заказа с номером %s", id)
	}
	return ""
}

func (s *ClientService) Login(config entities.AppConfig) (*entities.SessionInfo, string) {
	session, err := s.gw.Login(config)
	if err != nil {
		log.Println(err.Error())
		switch errors.Unwrap(err) {
		case errorlog.EmptyURLDataError:
			return nil, "Нет данных по хосту. Пожалуйста, добавьте адрес в настройках"
		case errorlog.InvalidLoginOrPassword:
			return nil, "Логин или пароль не заполнены или заполнены некорректно"
		case errorlog.AuthorizationError:
			return nil, "Произошла ошибка авторизации"
		case errorlog.JsonUnmarshalError:
			return nil, "Неправильный логин или пароль"
		default:
			return nil, "Произошла непредвиденная ошибка"
		}
	}
	return session, ""
}
