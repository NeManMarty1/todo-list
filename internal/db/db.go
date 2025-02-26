package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
