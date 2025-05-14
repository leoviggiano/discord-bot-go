//go:generate go run -mod=mod github.com/golang/mock/mockgen -package=mock -source=$GOFILE -destination=../../../test/mock/tx.go

package sql

import (
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Tx interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	Commit() error
	Rollback() error
}

type Transactioner interface {
	Begin() (Tx, error)
}

func NewTransactioner(db Connection, logger *slog.Logger) Transactioner {
	return sqlxTransactioner{
		logger: logger,
		db:     db,
	}
}
