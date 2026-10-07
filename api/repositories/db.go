package repository

import (
	"database/sql"
	"time"

	"github.com/ngfenglong/ikou-backend/api/config"
)

func NewDBModel(cfg config.Config) (*DBModel, error) {
	db, err := sql.Open("postgres", cfg.DbSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DBModel{DB: db}, nil
}

type DBModel struct {
	DB *sql.DB
}
