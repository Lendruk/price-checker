package models

import (
	"price-tracker/db"
)

func RegisterWebhook(hook string) error {
	db := db.GetDb()

	statement, _ := db.Prepare("INSERT INTO webhooks (hook) VALUES (?)")

	_, insertionErr := statement.Exec(hook)

	return insertionErr
}

func IsUserInWebhook(hook int, userId int) bool {
	row := db.GetDb().QueryRow("SELECT hook FROM webhook_users WHERE hook = ? AND user = ?", hook, userId)
	var result int
	err := row.Scan(&result)

	if err != nil {
		return false
	}

	return true
}

func AddUserToWebHook(hook string, userId int) error {
	db := db.GetDb()

	var hookId int
	row := db.QueryRow("SELECT id from webhooks WHERE hook = ?", hook)
	err := row.Scan(&hookId)

	if err != nil {
		return err
	}

	if IsUserInWebhook(hookId, userId) == false {
		_, insertionErr := db.Exec("INSERT INTO webhook_users (hook, user) VALUES (?, ?)", hookId, userId)
		return insertionErr
	}
	return nil
}
