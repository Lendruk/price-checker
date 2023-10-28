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
	type ProductWithWatchlist struct {
		models.Product
		Watchers []int `json:"watchers"`
	}

	type ProductMap map[int]ProductWithWatchlist
	// Registers cron job
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(15).Second().WaitForSchedule().Do(func() {
		fmt.Println("Updating products")
		// Get all vendor products that have watchers
		productIds, _ := models.GetAllProductsInWatchlists()

		fmt.Println(productIds)

		vendorEntries := make([]models.VendorEntry, 0)
		for _, id := range productIds {
			product := models.GetProductById(id)

			for _, entry := range product.VendorEntries {
				vendorEntries = append(vendorEntries, entry)
			}
		}

		// Run parsers.UpdateProducts (need to update method to return products that have changes)
		updatedEntries := parsers.UpdateProducts(vendorEntries)

		// TODO: Refactor this in the future...
		if len(updatedEntries) > 0 {
			for _, webhook := range models.GetRegisteredWebhooks() {
				watchedEntries := make(ProductMap)
				users := models.GetWebbhookUsers(webhook.Id)

				for _, user := range users {
					watchlist := models.FetchUserWatchList(user)

					for _, watchedProduct := range watchlist {
						for _, updatedEntry := range updatedEntries {
							if updatedEntry.UniversalId == watchedProduct.Id {
								_, ok := watchedEntries[updatedEntry.UniversalId]
								if ok == false {
									watchedEntries[updatedEntry.UniversalId] = ProductWithWatchlist{
										Watchers: make([]int, 0),
										Product:  watchedProduct,
									}
								}
								entry := watchedEntries[updatedEntry.UniversalId]

								alreadyInList := false
								for _, w := range entry.Watchers {
									if w == user {
										alreadyInList = true
										break
									}
								}

								if !alreadyInList {
									entry.Watchers = append(entry.Watchers, user)
									watchedEntries[updatedEntry.UniversalId] = entry
								}
							}
						}
					}
				}
				body, err := json.Marshal(watchedEntries)

				if err != nil {
					panic(err)
				}
				http.Post(webhook.Hook, "application/json", bytes.NewBuffer(body))
			}
		}

		fmt.Println("Product update done!")
	})

	scheduler.StartAsync()
}
