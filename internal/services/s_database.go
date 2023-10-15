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
	_, err := s.gw.DatabaseInterface.UploadSellsNote(dto)
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
