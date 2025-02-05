package appinfoRepositories

import (
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepositories interface {
}

type appinfoRepository struct {
	db *sqlx.DB
}

func NewappinfoRepository(db *sqlx.DB) IAppinfoRepositories {
	return &appinfoRepository{
		db: db,
	}
}
