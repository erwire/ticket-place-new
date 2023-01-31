package db

type SQLite struct {
	AuthDB
	HistoryDB
}

type AuthDB interface {
	CreateConnection()
}

type HistoryDB interface {
	CreateHistory()
}
