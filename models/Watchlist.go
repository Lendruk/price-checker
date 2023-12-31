package models

import (
	"price-tracker/db"
)

type WatchlistEntry struct {
	product int64
	user    int64
}

func AddProductToWatchlist(userId int, productId int) error {
	db := db.GetDb()

	statement, err := db.Prepare("INSERT INTO watchlists (product, user) values (?, ?)")

	defer statement.Close()

	_, insertionError := statement.Exec(productId, userId)

	if insertionError != nil {
		panic(insertionError)
	}

	return err
}

func RemoveProductFromWatchlist(userId int, productSKU string) error {
	db := db.GetDb()

	productId, err := GetProductId(productSKU)
	statement, err := db.Prepare("DELETE FROM watchlists WHERE user = ? AND product = ?")

	statement.Exec(userId, productId)

	return err
}

func IsProductInWatchlist(userId int, productId int) bool {
	db := db.GetDb()
	row := db.QueryRow("SELECT * FROM watchlists WHERE product = ? AND user = ?", productId, userId)

	var result WatchlistEntry

	err := row.Scan(&result.user, &result.product)

	if err != nil {
		return false
	}
	return true
}

func GetAllProductsInWatchlists() ([]int, error) {
	db := db.GetDb()
	var products []int

	rows, err := db.Query("SELECT product FROM watchlists")

	for rows.Next() {
		var product int
		rows.Scan(&product)

		products = append(products, product)
	}
	return products, err
}
