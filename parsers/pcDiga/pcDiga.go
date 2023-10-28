package pcDiga

import (
	"fmt"
	"os"
	"price-tracker/models"
	"price-tracker/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
)

const PcDigaUrl = "https://www.pcdiga.com"

func mapAvailability(availabilityText string) models.Availability {
	productAvailability := models.InStock
	switch availabilityText {
	case "Esgotado":
		productAvailability = models.OutOfStock
	case "Por Encomenda":
		productAvailability = models.ByOrder
	case "Poucas Unidades":
		productAvailability = models.InStock
	}
	return productAvailability
}

func ParseQueryPage(html string) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		panic(err)
	}

	products := make([]models.VendorEntry, 0)
	document.Find("div[class='grid justify-between gap-x-2 gap-y-4 grid-cols-2 sm:grid-cols-3 lg:grid-cols-prod-list 2xl:grid-cols-prod-list-lg']").Children().Each(func(i int, s *goquery.Selection) {
		productLinkElement := s.Find("a")
		productLink := PcDigaUrl + strings.TrimSpace(productLinkElement.AttrOr("href", ""))
		productName := strings.TrimSpace(productLinkElement.Find("span").Text())
		skuElement := productLinkElement.Next().Next()
		productSku := skuElement.Text()
		availabilityElement := skuElement.Next()
		availabilityText := availabilityElement.Find(".stock_availability").Text()

		productAvailability := mapAvailability(availabilityText)

		priceElement := availabilityElement.Next()
		productPrice, _ := utils.FormatPrice(priceElement.Text())
		products = append(products, models.NewVendorProduct(productName, productPrice, productLink, models.PCDiga, productSku, "", productAvailability))
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

	productElement := document.Find("div[class='md:p-4 lg:bg-background-off lg:rounded-md hidden lg:grid gap-y-4']")

	productPrice, _ := utils.FormatPrice(productElement.Find("div[class='text-primary text-2xl md:text-3xl font-black']").Text())
	availabilityText := strings.TrimSpace(productElement.Find(".stock_availability").Text())
	productAvailability := mapAvailability(availabilityText)

	updated, entry, err := models.UpdateVendorEntry(productPrice, productAvailability, sku, vendor)

	return updated, entry
}

func QueryProduct(productName string, browser *rod.Browser) {
	url := PcDigaUrl + "/search?query=" + productName
	fmt.Println(url)

	data, _ := os.ReadFile("./pcDigaSearchPage.html")
	html := string(data)
	ParseQueryPage(html)

	// page := browser.MustPage(url)
	// page.MustWaitStable()
	// html, err := page.HTML()

	// if err != nil {
	// 	panic(err)
	// }

	// os.WriteFile("./pcDigaSearchPage.html", []byte(html), 0644)
}

func UpdateProduct(product models.VendorEntry, browser *rod.Browser) (bool, models.VendorEntry) {
	url := product.Url
	fmt.Println(url)

	data, _ := os.ReadFile("./pcDigaProductPage.html")
	html := string(data)

	updated, entry := CheckProductPageForUpdates(html, product.SKU, product.Vendor)

	// page := browser.MustPage(url)
	// // Wait stable being funky for some reason
	// time.Sleep(3 * time.Second)

	// html, err := page.HTML()
	// os.WriteFile("./pcDigaProductPage.html", []byte(html), 0644)
	// if err != nil {
	// 	panic(err)
	// }

	return updated, entry
}

func CreateFromProductPage(url string, browser *rod.Browser) (models.Product, error) {
	fmt.Println(url)

	// TODO replace with real browser call
	data, _ := os.ReadFile("./pcDigaProductPage.html")
	html := string(data)

	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return models.Product{}, err
	}

	productElement := document.Find("div[class='md:p-4 lg:bg-background-off lg:rounded-md hidden lg:grid gap-y-4']")

	productPrice, _ := utils.FormatPrice(productElement.Find("div[class='text-primary text-2xl md:text-3xl font-black']").Text())
	srcAttr, _ := document.Find("img").First().Attr("src")
	productImageUrl := strings.TrimSpace(srcAttr)
	availabilityText := strings.TrimSpace(productElement.Find(".stock_availability").Text())
	productAvailability := mapAvailability(availabilityText)
	productFullName := strings.TrimSpace(productElement.Find("h1[class='font-bold text-2xl']").Text())
	productSKU := strings.TrimSpace(document.Find("div[class='flex flex-col lg:block text-xs']").First().Text())
	productSKU = strings.Split(productSKU, " ")[1]

	vendorProduct := models.NewVendorProduct(productFullName, productPrice, url, models.PCDiga, productSKU, productImageUrl, productAvailability)

	return models.InsertProduct(vendorProduct), nil
}
