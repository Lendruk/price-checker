package parsers

import (
	"errors"
	"price-tracker/models"
	"price-tracker/parsers/globalData"
	"price-tracker/parsers/pcDiga"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func RegisterProductByName(productName string) {
	browser := rod.New().MustConnect()
	defer browser.Close()

	// Search GD
	globalData.QueryProduct(productName, browser)
	// Search PCDiga
	pcDiga.QueryProduct(productName, browser)
}

func RegisterProductFromUrl(url string) (models.Product, error) {
	browser := rod.New().MustConnect()
	defer browser.Close()

	if strings.Contains(url, "globaldata") {
		return globalData.CreateFromProductPage(url, browser)
	} else if strings.Contains(url, "pcdiga") {
		return pcDiga.CreateFromProductPage(url, browser)
	}
	return models.Product{}, errors.New("Sent url does not have a parser")
}

func UpdateProducts(products []models.VendorEntry) []models.VendorEntry {
	browser := rod.New().MustConnect()
	defer browser.Close()

	updatedEntries := make([]models.VendorEntry, 0)
	for _, entry := range products {
		switch entry.Vendor {
		case models.GlobalData:
			updated, entry := globalData.UpdateProduct(entry, browser)
			if updated {
				updatedEntries = append(updatedEntries, entry)
			}
		case models.PCDiga:
			updated, entry := pcDiga.UpdateProduct(entry, browser)
			if updated {
				updatedEntries = append(updatedEntries, entry)
			}
		}

		time.Sleep(2 * time.Second)
	}

	return updatedEntries
}
