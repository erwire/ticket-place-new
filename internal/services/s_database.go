package services

import (
	"fmt"
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	"fptr/internal/gateways"
	"github.com/google/logger"
)

type DatabaseService struct {
	*logger.Logger
	gw *gateways.Gateway
}

func NewDatabaseService(logger *logger.Logger, gw *gateways.Gateway) *DatabaseService {
	return &DatabaseService{Logger: logger, gw: gw}
}

func (s *DatabaseService) UploadSellsNote(dto entities.SellsDTO) error {
	_, err := s.gw.UploadSellsNote(dto)
	if err != nil {
		message := "Ошибка в процессе создания записи в истории операции продаж"
		logger.Errorf(message+" -> %v", err)
		return apperr.NewDatabaseError("Ошибка в процессе создания записи в истории операции продаж", err)
	}
	return nil
}
func (s *DatabaseService) GetAllSellsNote() ([]entities.SellsDTO, error) {
	notes, err := s.gw.GetAllSellsNote()
	if err != nil {
		message := "Ошибка в процессе получения всех записей истории операций продаж"
		logger.Errorf(message+" -> %v", err)
		return notes, apperr.NewDatabaseError(message, err)
	}
	return notes, nil
}
func (s *DatabaseService) DeleteSellsNote(sellID string) error {
	err := s.gw.DeleteSellsNote(sellID)
	if err != nil {
		message := "Ошибка в процессе удаления записей истории операций продаж"
		logger.Errorf(message+" -> %v", err)
		return apperr.NewDatabaseError(message, err)
	}
	return nil
}
func (s *DatabaseService) GetSellNoteByID(sellID string) (entities.SellsDTO, error) {
	note, err := s.gw.GetSellNoteByID(sellID)
	if err != nil {
		message := fmt.Sprintf("Ошибка в процессе получения записи из таблицы истории операций продаж с номером %s", sellID)
		logger.Errorf(message+" -> %v", err)
		return note, apperr.NewDatabaseError(message, err)
	}
	return note, nil
}
func (s *DatabaseService) DeleteAllSellsNote() error {
	err := s.gw.DeleteAllSellsNote()
	if err != nil {
		message := "Ошибка в процессе удаления всех записей в таблице истории операций продаж"
		logger.Errorf(message+" -> %v", err)
		return apperr.NewDatabaseError(message, err)
	}
	return nil
}
func (s *DatabaseService) GetUnfinishedSellsNote(status string) ([]entities.SellsDTO, error) {
	notes, err := s.gw.GetUnfinishedSellsNote(status)
	if err != nil {
		message := "Ошибка в процессе получения записей (по статусу) истории операций продаж"
		logger.Errorf(message+" -> %v", err)
		return notes, apperr.NewDatabaseError(message, err)
	}
	return notes, nil
}

func (s *DatabaseService) UploadUsers(dto entities.Users) error {
	_, err := s.gw.UploadUsers(dto)
	if err != nil {
		message := "Ошибка в процессе создания записи в таблице пользователей"
		logger.Errorf(message+" -> %v", err)
		return apperr.NewDatabaseError(message, err)
	}
	logger.Infof(fmt.Sprintf("Успешное добавление пользователя %s в базу", dto.Login))
	return nil
}
func (s *DatabaseService) GetAllUsers() ([]entities.Users, error) {
	users, err := s.gw.GetAllUsers()
	if err != nil {
		message := "Ошибка в процессе получения записей в таблице пользователей"
		logger.Errorf(message+" -> %v", err)
		return users, apperr.NewDatabaseError(message, err)
	}
	return users, nil
}
func (s *DatabaseService) GetUser(login string) (entities.Users, error) {
	user, err := s.gw.GetUser(login)
	if err != nil {
		message := "Ошибка в процессе получения записи в таблице пользователей"
		logger.Errorf(message+" -> %v", err)
		return user, apperr.NewDatabaseError(message, err)
	}
	return user, nil
}
