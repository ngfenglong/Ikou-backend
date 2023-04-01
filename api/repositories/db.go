package repository

import (
	"database/sql"
	"fmt"
	"ikou/api/config"
)

func NewDBModel(cfg config.Config) (*DBModel, error) {
	db, err := sql.Open("mysql", cfg.DbSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &DBModel{DB: db}, nil
}

type DBModel struct {
	DB *sql.DB
}
