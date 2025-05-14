package repository

import (
	_ "embed"

	"github.com/nleof/goyesql"

	"megumin/pkg/lib/sql"
)

//go:embed user.sql
var userQueries []byte

type User interface {
	Create(id string) error
	Exists(id string) (bool, error)
}

type user struct {
	db      sql.Connection
	queries goyesql.Queries
}

func NewUser(db sql.Connection) User {
	return &user{
		db:      db,
		queries: goyesql.MustParseBytes(userQueries),
	}
}

func (u user) Create(id string) error {
	_, err := u.db.Exec(u.queries["create"], id)
	return err
}

func (u user) Exists(id string) (bool, error) {
	var exists bool
	err := u.db.Get(&exists, u.queries["exists"], id)
	if err != nil {
		return false, err
	}

	return exists, nil
}
