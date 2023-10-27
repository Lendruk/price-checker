package models

import (
	"price-tracker/db"
)

type User struct {
	Id        int64 `json:"id"`
	Watchlist []Product
}

func NewUser() User {
	result, err := db.GetDb().Exec("INSERT INTO users DEFAULT VALUES")

	if err != nil {
		panic(err)
	}

	newId, _ := result.LastInsertId()

	return User{
		Id:        newId,
		Watchlist: make([]Product, 0),
	}
}

func GetUserById(userId int64) (User, error) {
	userResult := db.GetDb().QueryRow("SELECT * FROM users WHERE id = ?", userId)

	var user User
	err := userResult.Scan(&user.Id)

	if err != nil {
		return User{}, err
	}

	user.Watchlist = FetchUserWatchList(int(userId))

	return user, nil
}

func FetchUserWatchList(userId int) []Product {
	db := db.GetDb()

	rows, err := db.Query("SELECT product from watchlists WHERE user = ?", userId)

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	products := make([]Product, 0)
	for rows.Next() {
		var productId int

		err := rows.Scan(&productId)
		if err != nil {
			panic(err)
		}

		products = append(products, GetProductById(productId))
	}

	return products
}
