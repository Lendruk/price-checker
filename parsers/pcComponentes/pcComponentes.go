package pccomponentes

import (
	"fmt"
	"price-tracker/models"
	"price-tracker/utils"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

func UpdateProduct(product models.VendorEntry, browser *rod.Browser) (bool, models.VendorEntry) {
	url := product.Url
	fmt.Println(url)

	// data, _ := os.ReadFile("./pcComponentesProductPageUpdate.html")
	// html := string(data)

	page := stealth.MustPage(browser)
	page.MustNavigate(url)
	// Wait stable being funky for some reason
	time.Sleep(3 * time.Second)

	updated, entry := CheckProductPageForUpdates(page.MustHTML(), product.SKU, product.Vendor)

	return updated, entry
}

func CreateFromProductPage(url string, browser *rod.Browser) (models.Product, error) {
	fmt.Println(url)

	// data, _ := os.ReadFile("./pcComponentesProductPage.html")
	// html := string(data)
	page := stealth.MustPage(browser)
	page.MustNavigate(url)
	// Wait stable being funky for some reason
	time.Sleep(3 * time.Second)

	// data, _ := os.ReadFile("./globalDataProductPage.html")
	// html := string(data)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(page.MustHTML()))

	if err != nil {
		panic(err)
	}

	productElement := document.Find("#pdp-section-buybox")
	productFullName := document.Find("#pdp-title").Text()
	productPrice := productElement.Find("#pdp-price-current-integer").Text()
	productPriceFormatted, _ := utils.FormatPrice(strings.TrimSpace(productPrice))
	productImage, _ := document.Find(".sc-iBbrVh.jhaxZx").Attr("src")
	productSKU := document.Find("#pdp-mpn").Text()
	productSKU = strings.Split(productSKU, " ")[1]
	productSKU = strings.ReplaceAll(productSKU, "Cod.", "")

	vendorProduct := models.NewVendorProduct(productFullName, productPriceFormatted, url, models.PcComponentes, productSKU, "https:"+productImage, models.InStock)

	return models.InsertProduct(vendorProduct), nil

}

func CheckProductPageForUpdates(html string, sku string, vendor models.Vendor) (bool, models.VendorEntry) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		panic(err)
	}

	rawPrice := document.Find("#pdp-price-current-integer").Text()
	productAvailability := models.InStock

	if rawPrice == "" {
		return false, models.VendorEntry{}
	}
	productPrice, _ := utils.FormatPrice(rawPrice)

	updated, entry, err := models.UpdateVendorEntry(productPrice, productAvailability, sku, vendor)
	return updated, entry
}
