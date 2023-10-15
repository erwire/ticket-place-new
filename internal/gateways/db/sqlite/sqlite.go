package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"fptr/internal/entities"
	"github.com/jmoiron/sqlx"
	"strconv"
)

const (
	dbSource = "./db/sqlite.db"
)

const (
	dataTable    = "data"
	historyTable = "history"
	sellTable    = "sells"
)

type SqliteDB struct {
	db *sqlx.DB
}

func NewSqliteDB(db *sqlx.DB) *SqliteDB {
	return &SqliteDB{db: db}
}

// Для работы с таблицой Sells

func (s *SqliteDB) UploadSellsNote(dto entities.SellsDTO) (string, error) {
	if _, err := s.GetSellNoteByID(dto.SellID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id, err := s.CreateSellsNote(dto)
			if err != nil {
				return "", err
			}
			return id, nil
		} else {
			return "", err
		}
	}

	id, err := s.UpdateSellsNote(dto)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *SqliteDB) CreateSellsNote(dto entities.SellsDTO) (string, error) {
	query := fmt.Sprintf("INSERT INTO %s(sell_id, date, amount, status, error, event) VALUES(:sell_id, :date, :amount, :status, :error, :event)", sellTable)
	exec, err := s.db.NamedExec(query, &dto)
	if err != nil {
		return "", err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(id)), nil
}

func (s *SqliteDB) UpdateSellsNote(dto entities.SellsDTO) (string, error) {
	query := fmt.Sprintf("UPDATE %s SET date=:date, amount=:amount, status=:status, error=:error, event=:event WHERE sell_id = :sell_id", sellTable)
	exec, err := s.db.NamedExec(query, &dto)

	if err != nil {
		return "", err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(id)), nil
}

func (s *SqliteDB) GetAllSellsNote() ([]entities.SellsDTO, error) {
	var dto []entities.SellsDTO
	query := fmt.Sprintf("SELECT * FROM %s", sellTable)
	err := s.db.Select(&dto, query)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func (s *SqliteDB) GetUnfinishedSellsNote(status string) ([]entities.SellsDTO, error) {
	var dto []entities.SellsDTO
	query := fmt.Sprintf("SELECT * FROM %s WHERE status=$1", sellTable)
	err := s.db.Select(&dto, query, status)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func (s *SqliteDB) DeleteSellsNote(sellID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE sell_id=$1", sellTable)
	_, err := s.db.Exec(query, sellID)
	return err
}

func (s *SqliteDB) GetSellNoteByID(sellID string) (entities.SellsDTO, error) {
	var dto entities.SellsDTO
	query := fmt.Sprintf("SELECT * FROM %s WHERE sell_id = $1", sellTable)
	err := s.db.Get(&dto, query, sellID)
	if err != nil {
		return dto, err
	}
	return dto, nil
}

func (s *SqliteDB) DeleteAllSellsNote() error {
	query := fmt.Sprintf("DELETE FROM %s", sellTable)
	_, err := s.db.Exec(query)
	return err
}
