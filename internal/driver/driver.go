package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the database connection pool.
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const (
	// maxOpenConns has the maximum number of open connections to the database.
	maxOpenConns = 10
	// maxIdleConns has the maximum number of connections in the idle connection pool.
	maxIdleConns = 5
	// maxDBLifetime is the maximum amount of time a connection may be reused.
	maxDBLifetime = 5 * time.Minute
)

// ConnectSQL creates a connection pool to our PostgreSQL database.
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenConns)
	d.SetMaxIdleConns(maxIdleConns)
	d.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// NewDatabase creates a new database connection.
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// testDB tries to ping the database to see if it's alive.
func testDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}

	return nil
}
