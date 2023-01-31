package db

import "database/sql"

type SQLiteHistory struct {
	db sql.DB
}

func (d *SQLiteHistory) CreateHistory() {

}
