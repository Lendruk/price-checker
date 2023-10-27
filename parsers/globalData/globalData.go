package globalData

import (
	"fmt"
	"os"
	"price-tracker/models"
	"price-tracker/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

const GlobalDataUrl = "https://www.globaldata.pt"

func ParseQueryPage(html string) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		panic(err)
	}

	products := make([]models.VendorEntry, 0)
	document.Find("article.ck-product-box").Each(func(i int, selection *goquery.Selection) {
		productLinkElement := selection.Find(".js-gtm-product-link-algolia")
		productLink := GlobalDataUrl + strings.TrimSpace(productLinkElement.AttrOr("href", ""))
		fullProductName := strings.TrimSpace(productLinkElement.Text())
		availabilityText := strings.TrimSpace(selection.Find(".availability-product").Find(".small").Text())
		productAvailability := models.InStock
		switch availabilityText {
		case "Esgotado":
			productAvailability = models.OutOfStock
		case "Por encomenda":
			productAvailability = models.ByOrder
		case "Poucas unidades":
			productAvailability = models.InStock
		}
		productSku := selection.Find(".ck-product-box-sku").Text()
		productPrice, _ := utils.FormatPrice(selection.Find(".price__amount").Text())
		products = append(products, models.NewVendorProduct(fullProductName, productPrice, productLink, models.GlobalData, productSku, productAvailability))
	})

	for _, v := range products {
		if models.DoesVendorProductExist(v.SKU, v.Vendor) == false {
			models.InsertProduct(v)
		}
		// TODO: Update product
	}
}

func CheckProductPageForUpdates(html string, sku string, vendor models.Vendor) (bool, models.VendorEntry) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		panic(err)
	}

	productElement := document.Find(".ck-product-cta-box-inner")
	availabilityText := strings.TrimSpace(productElement.Find(".custom-element.availability-product").Find(".small").Text())
	availabilityText = strings.TrimSpace(strings.Split(availabilityText, "-")[1])
	productAvailability := models.InStock

	switch availabilityText {
	case "Esgotado":
		productAvailability = models.OutOfStock
	case "Por encomenda":
		productAvailability = models.ByOrder
	case "Poucas unidades":
		productAvailability = models.InStock
	}

	productPrice, _ := utils.FormatPrice(strings.TrimSpace(productElement.Find(".price__amount").Text()))

	updated, entry, err := models.UpdateVendorEntry(productPrice, productAvailability, sku, vendor)

	return updated, entry
}

// Queries for products according to the provided name
func QueryProduct(productName string, browser *rod.Browser) {
	url := GlobalDataUrl + "/?query=" + strings.ReplaceAll(productName, " ", "%2520")
	fmt.Println(url)

	data, _ := os.ReadFile("./globalDataSearchPage.html")
	html := string(data)

	// page := browser.MustPage(url)
	// page.MustWaitStable()

	// html, err := page.HTML()

	// if err != nil {
	// 	panic(err)
	// }

	ParseQueryPage(html)
}

func UpdateProduct(product models.VendorEntry, browser *rod.Browser) (bool, models.VendorEntry) {
	url := product.Url
	fmt.Println(url)

	data, _ := os.ReadFile("./globalDataProductPageUpdate.html")
	html := string(data)

	updated, entry := CheckProductPageForUpdates(html, product.SKU, product.Vendor)

	// page := browser.MustPage(url)
	// // Wait stable being funky for some reason
	// time.Sleep(3 * time.Second)

	// html, err := page.HTML()
	// os.WriteFile("./globalDataProductPage.html", []byte(html), 0644)
	// if err != nil {
	// 	panic(err)
	// }

	return updated, entry
}

func CreateFromProductPage(url string, browser *rod.Browser) (models.Product, error) {
	fmt.Println(url)

	// TODO replace with real browser call
	data, _ := os.ReadFile("./globalDataProductPage.html")
	html := string(data)

	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return models.Product{}, err
	}

	productElement := document.Find(".ck-product-cta-box-inner")
	availabilityText := strings.TrimSpace(productElement.Find(".custom-element.availability-product").Find(".small").Text())
	productFullName := strings.TrimSpace(productElement.Find("h1[class='mb-0 h4 h2-md']").Text())
	availabilityText = strings.TrimSpace(strings.Split(availabilityText, "-")[1])
	productSKU := strings.TrimSpace(document.Find("div[class='ck-product-sku-ean-warranty-info__item small d-inline-block my-1']").First().Text())
	productSKU = strings.Split(productSKU, " ")[1]
	productAvailability := models.InStock

	switch availabilityText {
	case "Esgotado":
		productAvailability = models.OutOfStock
	case "Por encomenda":
		productAvailability = models.ByOrder
	case "Poucas unidades":
		productAvailability = models.InStock
	}

	productPrice, _ := utils.FormatPrice(strings.TrimSpace(productElement.Find(".price__amount").Text()))

	vendorProduct := models.NewVendorProduct(productFullName, productPrice, url, models.GlobalData, productSKU, productAvailability)

	return models.InsertProduct(vendorProduct), nil
}
