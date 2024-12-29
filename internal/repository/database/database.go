package database

import (
	"database/sql"

	"github.com/mayloo89/bamos/internal/config"
	"github.com/mayloo89/bamos/internal/repository"
)

type postgreDatabase struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDatabase struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgreDatabase(conn *sql.DB, app *config.AppConfig) repository.DBRepository {
	return &postgreDatabase{
		App: app,
		DB:  conn,
	}
}

func NewTestingDatabase(app *config.AppConfig) repository.DBRepository {
	return &testDatabase{
		App: app,
	}
}

// Implement the methods of repository.DBRepository interface
func (p *postgreDatabase) SomeMethod() error {
	// Implementation here
	return nil
}
