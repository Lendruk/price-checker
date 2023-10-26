package models

import (
	"price-tracker/db"
)

type Webhook struct {
	Id    int
	Hook  string
	users []int
}

func RegisterWebhook(hook string) error {
	db := db.GetDb()

	statement, _ := db.Prepare("INSERT INTO webhooks (hook) VALUES (?)")

	_, insertionErr := statement.Exec(hook)

	return insertionErr
}

func GetRegisteredWebhooks() []Webhook {
	db := db.GetDb()

	hooks := make([]Webhook, 0)

	rows, err := db.Query("SELECT * FROM webhooks")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var hook Webhook
		rows.Scan(&hook.Id, &hook.Hook)
		hooks = append(hooks, hook)

		hook.users = GetWebbhookUsers(hook.Id)
	}

	return hooks
}

func GetWebbhookUsers(hookId int) []int {
	users := make([]int, 0)
	rows, err := db.GetDb().Query("SELECT user FROM webhook_users WHERE hook = ?", hookId)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var user int
		rows.Scan(&user)

		users = append(users, user)
	}

	return users
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
