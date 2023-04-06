package repository

import (
	"database/sql"

	"github.com/ngfenglong/ikou-backend/api/config"
)

func NewDBModel(cfg config.Config) (*DBModel, error) {
	db, err := sql.Open("mysql", cfg.DbSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DBModel{DB: db}, nil
}

type DBModel struct {
	DB *sql.DB
}
