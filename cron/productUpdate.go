package cron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"price-tracker/models"
	"price-tracker/parsers"
	"time"

	"github.com/go-co-op/gocron"
)

func RegisterProductUpdateCronJob() {

	// Registers cron job
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(30).Second().WaitForSchedule().Do(func() {
		fmt.Println("Updating products")
		// Get all vendor products that have watchers
		productIds, _ := models.GetAllProductsInWatchlists()

		vendorEntries := make([]models.VendorEntry, 0)
		for _, id := range productIds {
			product := models.GetProductById(id)

			for _, entry := range product.VendorEntries {
				vendorEntries = append(vendorEntries, entry)
			}
		}
		// Run parsers.UpdateProducts (need to update method to return products that have changes)
		updatedEntries := parsers.UpdateProducts(vendorEntries)

		for _, webhook := range models.GetRegisteredWebhooks() {

			body, err := json.Marshal(updatedEntries)

			if err != nil {
				panic(err)
			}
			// TODO: Separate updated list according to the webhook with users that have them in the watchlist
			http.Post(webhook.Hook, "application/json", bytes.NewBuffer(body))
		}
	})

	scheduler.StartAsync()
}
