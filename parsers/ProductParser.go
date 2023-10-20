package parsers

import (
	"price-tracker/models"
	"price-tracker/parsers/globalData"
	"price-tracker/parsers/pcDiga"

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

func UpdateProducts(products []models.VendorEntry) {
	browser := rod.New().MustConnect()
	defer browser.Close()

	for _, entry := range products {
		switch entry.Vendor {
		case models.GlobalData:
			globalData.UpdateProduct(entry, browser)
		case models.PCDiga:
			pcDiga.UpdateProduct(entry, browser)
		}
	}
}
