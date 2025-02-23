package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *sql.DB

func InitDB() error {
	if db != nil {
		return nil
	}

	dbUrl := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")

	if dbUrl == "" || dbToken == "" {
		return fmt.Errorf("DB_URL and DB_TOKEN must be set")
	}

	var err error
	db, err = sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", dbUrl, dbToken))
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to database")
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
