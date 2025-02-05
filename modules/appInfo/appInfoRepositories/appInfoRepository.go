package appInfoRepositories

import (
	"github.com/jmoiron/sqlx"
)

type IAppInfoRepositories interface {
}

type appInfoRepository struct {
	db *sqlx.DB
}

func NewAppInfoRepository(db *sqlx.DB) IAppInfoRepositories {
	return &appInfoRepository{
		db: db,
	}
}
