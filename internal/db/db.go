package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// InitDB init db connection
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
