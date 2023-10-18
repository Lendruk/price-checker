package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// connects and initializes the database
func init() {
	sqlDb, err := sql.Open("sqlite3", "./price-tracker.db")

	if err != nil {
		panic(err)
	}

	db = sqlDb
}

func GetDb() *sql.DB {
	return db
}
