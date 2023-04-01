package store

import (
	"ikou/api/config"
	repository "ikou/api/repositories"
)

type Store struct {
	Config config.Config
	DB     repository.DBModel
}

func NewStore(cfg config.Config) (*Store, error) {
	db, err := repository.NewDBModel(cfg)
	if err != nil {
		return nil, err
	}

	store := &Store{
		Config: cfg,
		DB:     *db,
	}

	return store, nil
}
