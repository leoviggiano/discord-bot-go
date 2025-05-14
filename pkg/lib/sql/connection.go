package sql

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"regexp"
	"time"

	// postgres driver
	_ "github.com/lib/pq"
)

var (
	space     = regexp.MustCompile(`\s+`)
	ErrNoRows = sql.ErrNoRows
)

type Connection interface {
	io.Closer
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	In(query string, args ...interface{}) (string, []interface{}, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	Ping() error
	begin() (driver.Tx, error)
}

type Settings struct {
	Conn            string
	MaxIdleCons     int
	MaxOpenCons     int
	ConnMaxLifetime time.Duration
}
