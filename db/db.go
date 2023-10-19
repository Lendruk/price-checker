package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// connects and initializes the database
func init() {
	sqlDb, err := sql.Open("sqlite3", "./price-tracker.db")

	sqlDb.Exec("CREATE TABLE IF NOT EXISTS vendorEntries (id INTEGER PRIMARY KEY, fullName TEXT, price REAL, url TEXT, vendor TEXT, sku TEXT, availability INTEGER, lastUpdated INTEGER, universalId INTEGER, FOREIGN KEY(universalId) REFERENCES products(id))")
	sqlDb.Exec("CREATE TABLE IF NOT EXISTS productHistory (id INTEGER PRIMARY KEY, productId INTEGER, price REAL, availability INTEGER, updatedAt INTEGER, FOREIGN KEY(productId) REFERENCES vendorEntries(id))")
	sqlDb.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, sku TEXT)")
	if err != nil {
		panic(err)
	}

	db = sqlDb
}

func GetDb() *sql.DB {
	return db
}
