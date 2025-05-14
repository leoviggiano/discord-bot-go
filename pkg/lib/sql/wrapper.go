package sql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type Wrapper struct {
	DB     *sqlx.DB
	Logger Logger
}

func (w Wrapper) DriverName() string {
	return w.DB.DriverName()
}

func (w Wrapper) Rebind(s string) string {
	return w.DB.Rebind(s)
}

func (w Wrapper) BindNamed(s string, i interface{}) (string, []interface{}, error) {
	return w.DB.BindNamed(s, i)
}

func (w Wrapper) log(query string, args ...interface{}) {
	if w.Logger != nil {
		w.Logger.Info("[POSTGRES] %s %v", space.ReplaceAllString(query, " "), args)
	}
}

func (w Wrapper) Close() error {
	w.log("closed")
	return w.DB.Close()
}

func (w Wrapper) Select(dest interface{}, query string, args ...interface{}) error {
	w.log(query, args...)
	return w.DB.Select(dest, query, args...)
}

func (w Wrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	w.log(query, args...)
	return w.DB.Query(query, args...)
}

func (w Wrapper) QueryRow(query string, args ...interface{}) *sql.Row {
	w.log(query, args...)
	return w.DB.QueryRow(query, args...)
}

func (w Wrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	w.log(query, args...)
	return w.DB.Exec(query, args...)
}

func (w Wrapper) Get(dest interface{}, query string, args ...interface{}) error {
	w.log(query, args...)
	return w.DB.Get(dest, query, args...)
}

func (w Wrapper) In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}

func (w Wrapper) NamedExec(query string, arg interface{}) (sql.Result, error) {
	w.log(query, arg)
	return w.DB.NamedExec(query, arg)
}
