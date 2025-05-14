package service

import (
	"log/slog"

	"megumin/pkg/lib/sql"
)

type All struct {
	User     User
	CustomID CustomID
}

func GetAll(db sql.Connection, logger *slog.Logger) All {
	return All{
		User:     NewUser(db, logger),
		CustomID: NewCustomID(logger),
	}
}
