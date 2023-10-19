package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// connects and initializes the database
func init() {
	sqlDb, err := sql.Open("sqlite3", "./price-tracker.db")

	sqlDb.Exec("CREATE TABLE IF NOT EXISTS vendorProducts (id INTEGER PRIMARY KEY, fullName TEXT, price REAL, url TEXT, vendor TEXT, sku TEXT, availability INTEGER, lastUpdated INTEGER)")

	if err != nil {
		panic(err)
	}

	db = sqlDb
}

func GetDb() *sql.DB {
	return db
}
