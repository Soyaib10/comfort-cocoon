package dbrepo

import (
	"database/sql"

	"github.com/Soyaib10/comfort-cocoon/internal/config"
	"github.com/Soyaib10/comfort-cocoon/internal/repository"
)

type mysqlDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}

func NewMySQLRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo{
	return &mysqlDBRepo{
		App: a,
		DB: conn,
	}
}