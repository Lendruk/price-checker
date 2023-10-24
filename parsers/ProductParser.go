package parsers

import (
	"errors"
	"price-tracker/models"
	"price-tracker/parsers/globalData"
	"price-tracker/parsers/pcDiga"
	"strings"

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
