package database

import (
	"log"

	"github.com/chakornpat-tn/go-rest-api/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDBConfig) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatal("connect to database error : ", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxConnections())
	return db
}
