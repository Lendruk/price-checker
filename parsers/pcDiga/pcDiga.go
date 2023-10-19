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

func ParseQueryPage(html string) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		panic(err)
	}

	products := make([]models.VendorProduct, 0)
	document.Find("div[class='grid justify-between gap-x-2 gap-y-4 grid-cols-2 sm:grid-cols-3 lg:grid-cols-prod-list 2xl:grid-cols-prod-list-lg']").Children().Each(func(i int, s *goquery.Selection) {
		productLinkElement := s.Find("a")
		productLink := strings.TrimSpace(productLinkElement.AttrOr("href", ""))
		productName := strings.TrimSpace(productLinkElement.Find("span").Text())
		skuElement := productLinkElement.Next().Next()
		productSku := skuElement.Text()
		availabilityElement := skuElement.Next()
		availabilityText := availabilityElement.Find(".stock_availability").Text()

		productAvailability := models.InStock
		switch availabilityText {
		case "Esgotado":
			productAvailability = models.OutOfStock
		case "Por Encomenda":
			productAvailability = models.ByOrder
		case "Poucas Unidades":
			productAvailability = models.InStock
		}

		priceElement := availabilityElement.Next()
		productPrice, _ := utils.FormatPrice(priceElement.Text())
		fmt.Println(productLink)
		fmt.Println(productName)
		fmt.Println(productSku)
		fmt.Println(productPrice)
		fmt.Println(availabilityText)
		fmt.Println(productAvailability)
		products = append(products, models.NewVendorProduct(productName, productPrice, productLink, "https://www.pcdiga.pt/", productSku, models.InStock))
	})

	for _, v := range products {
		if models.DoesProductExist(v.SKU, v.Vendor) == false {
			models.InsertProduct(v)
		}
	}
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
