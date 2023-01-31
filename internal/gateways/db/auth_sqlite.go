package db

import "database/sql"

type SQLiteAuth struct {
	db sql.DB
}

func (d *SQLiteAuth) CreateConnection() {

}
