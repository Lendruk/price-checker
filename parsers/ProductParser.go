package parsers

import (
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
