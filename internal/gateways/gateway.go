package gateways

import "fptr/internal/gateways/db"

type Gateway struct {
	db.SQLite
}
