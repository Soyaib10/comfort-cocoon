package dbrepo

import (
	"database/sql"

	"github.com/Soyaib10/comfort-cocoon/internal/config"
)

type mySqlDBRepo struct {
	DB  *sql.DB
	App *config.AppConfig
}

func NewMySqlDBRepo(db *sql.DB, app *config.AppConfig) *mySqlDBRepo {
	return &mySqlDBRepo{
		DB:  db,
		App: app,
	}
}
