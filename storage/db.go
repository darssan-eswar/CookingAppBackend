package storage

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func InitDB() error {
	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	var err error
	db, err = sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))
	if err != nil {
		return err
	}

	return nil
}

func GetDB() (*sql.DB, error) {
	if db == nil {
		err := InitDB()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
