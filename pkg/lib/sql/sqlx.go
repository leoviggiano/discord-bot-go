package sql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func connect(settings Settings, logger *slog.Logger) (Connection, error) {
	db, err := sqlx.Connect("postgres", settings.Conn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(settings.MaxIdleCons)
	db.SetMaxOpenConns(settings.MaxOpenCons)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime)

	return SQLX{
		DB:     db,
		Logger: logger,
	}, err
}

func NewSQLXConnection(settings Settings, log *slog.Logger) (Connection, error) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	timeoutExceeded := time.After(time.Hour)
	for {
		select {
		case <-timeoutExceeded:
			return nil, errors.New("db connection failed")

		case <-ticker.C:
			db, err := connect(settings, log)
			if err == nil {
				log.Info("[SQLX] Database connected with success")
				return db, nil
			}

			log.Error("Database connection failed", "error", err)
		}
	}
}

type SQLX struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

func (s SQLX) DriverName() string {
	return s.DB.DriverName()
}

func (s SQLX) Rebind(query string) string {
	return s.DB.Rebind(query)
}

func (s SQLX) Ping() error {
	return s.DB.Ping()
}

func (s SQLX) BindNamed(query string, i interface{}) (string, []interface{}, error) {
	return s.DB.BindNamed(query, i)
}

func (s SQLX) log(query string, args ...interface{}) {
	if s.Logger != nil {
		s.Logger.Info(fmt.Sprintf("[POSTGRES] %s %v", space.ReplaceAllString(query, " "), args))
	}
}

func (s SQLX) Close() error {
	if err := s.DB.Close(); err != nil {
		return err
	}

	s.log("closed")
	return nil
}

func (s SQLX) Select(dest interface{}, query string, args ...interface{}) error {
	s.log(query, args...)
	return s.DB.Select(dest, query, args...)
}

func (s SQLX) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s.log(query, args...)
	return s.DB.Query(query, args...)
}

func (s SQLX) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.log(query, args...)
	return s.DB.Exec(query, args...)
}

func (s SQLX) Get(dest interface{}, query string, args ...interface{}) error {
	s.log(query, args...)
	err := s.DB.Get(dest, query, args...)
	if err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (s SQLX) In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}

func (s SQLX) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return s.DB.NamedExec(query, arg)
}

func (s SQLX) begin() (driver.Tx, error) {
	return s.DB.Beginx()
}

type sqlxTransaction struct {
	tx     *sqlx.Tx
	logger *slog.Logger
}

func (t sqlxTransaction) log(query string, args ...interface{}) {
	if t.logger != nil {
		t.logger.Info(fmt.Sprintf("[POSTGRES][TX] %s %v", space.ReplaceAllString(query, " "), args))
	}
}

func (t sqlxTransaction) Get(dest interface{}, query string, args ...interface{}) error {
	t.log(query, args...)
	err := t.tx.Get(dest, query, args...)
	if err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (t sqlxTransaction) Select(dest interface{}, query string, args ...interface{}) error {
	t.log(query, args...)
	return t.tx.Select(dest, query, args...)
}

func (t sqlxTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	t.log(query, args...)
	return t.tx.Exec(query, args...)
}

func (t sqlxTransaction) NamedExec(query string, arg interface{}) (sql.Result, error) {
	t.log(query, arg)
	return t.tx.NamedExec(query, arg)
}

func (t sqlxTransaction) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	t.log(query, arg)
	return t.tx.NamedQuery(query, arg)
}

func (t sqlxTransaction) Commit() error {
	t.log("COMMIT")
	return t.tx.Commit()
}

func (t sqlxTransaction) Rollback() error {
	t.log("ROLLBACK")
	return t.tx.Rollback()
}

func (t sqlxTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	t.log(query, args...)
	return t.tx.QueryRow(query, args...)
}

type sqlxTransactioner struct {
	db     Connection
	logger *slog.Logger
}

func (t sqlxTransactioner) Begin() (Tx, error) {
	tx, err := t.db.begin()

	txr := sqlxTransaction{
		tx:     tx.(*sqlx.Tx),
		logger: t.logger,
	}
	txr.log("BEGIN")
	return txr, err
}

func NewMock(t *testing.T) (Connection, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return SQLX{DB: sqlxDB}, mock
}
