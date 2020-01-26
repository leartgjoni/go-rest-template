package postgres

import (
	"database/sql"
)

type DB struct {
	*sql.DB
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// check db is available
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
