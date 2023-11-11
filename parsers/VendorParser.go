package parsers

import (
	"price-tracker/models"

	"github.com/go-rod/rod"
)

type VendorParser interface {
	ParseQueryPage(html string)

	CheckProductPageForUpdates(html string, sku string, vendor models.Vendor) (bool, models.VendorEntry)

	QueryProduct(productName string, browser *rod.Browser)

	UpdateProduct(product models.Product, browser *rod.Browser) (bool, models.VendorEntry)

	CreateFromProductPage(url string, browser *rod.Browser) (models.Product, error)
}
