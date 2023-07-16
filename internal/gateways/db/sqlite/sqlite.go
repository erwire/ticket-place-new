package sqlite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	dbSource = "./db/sqlite.db"
)

const (
	dataTable    = "data"
	historyTable = "history"
)

type SqliteDB struct {
	db *sqlx.DB
}

func NewSqliteDB(db *sqlx.DB) *SqliteDB {
	return &SqliteDB{db: db}
}

// Для работы с таблицой Data

// CreateHistoryNoteInTX Создание записи в таблице истории
func (s *SqliteDB) CreateHistoryNoteInTX(dto HistoryDataDTO) (uint64, error) {
	queryData := fmt.Sprintf("INSERT INTO %s(data_type, data_content) VALUES(:data_type, :data_content)", dataTable)
	queryHistory := fmt.Sprintf("INSERT INTO %s(id, date, type, status, error, description, data_id) VALUES(:id, :date, :type, :status, :error, :description, :data_id)", historyTable)
	tx := s.db.MustBegin()
	result := tx.MustExec(queryData, &dto.DataDTO)
	id, err := result.LastInsertId()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, err
	}
	dto.HistoryDTO.DataID = uint64(id)
	result = tx.MustExec(queryHistory, &dto.HistoryDTO)
	err = tx.Commit()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, err
	}
	return uint64(id), nil

}

// CreateDataNote Создание записи в таблице данных
func (s *SqliteDB) CreateDataNote(dto DataDTO) (uint64, error) {
	query := fmt.Sprintf("INSERT INTO %s(data_type, data_content) VALUES(:data_type, :data_content)", dataTable)
	exec, err := s.db.NamedExec(query, &dto)
	if err != nil {
		return 0, err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), err
}

// CreateHistoryNote Создание записи в таблице истории
func (s *SqliteDB) CreateHistoryNote(dto HistoryDTO) (uint64, error) {
	query := fmt.Sprintf("INSERT INTO %s(id, date, type, status, error, description, data_id) VALUES(:id, :date, :type, :status, :error, :description, :data_id)", historyTable)
	exec, err := s.db.NamedExec(query, &dto)
	if err != nil {
		return 0, err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), err
}

func (s *SqliteDB) GetHistoryNote() {

}

func (s *SqliteDB) GetAllHistoryNote() {

}
