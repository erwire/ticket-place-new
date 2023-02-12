package services

import (
	"errors"
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

func (s *ClientService) Listen() error {
	return nil
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
