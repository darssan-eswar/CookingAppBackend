package utils

import (
	"database/sql"
	"fmt"
	"os"
)

func DBConnect() (*sql.DB, error) {
	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))
	if err != nil {
		return nil, err
	}

	return db, nil
}
