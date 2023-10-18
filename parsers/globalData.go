package parsers

import (
	"fmt"
	"os"
	"price-tracker/models"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const GlobalDataUrl = "https://www.globaldata.pt"

func ParseQueryPage(html string) {
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		panic(err)
	}

	products := make([]models.VendorProduct, 0)
	document.Find("article.ck-product-box").Each(func(i int, selection *goquery.Selection) {
		productLinkElement := selection.Find(".js-gtm-product-link-algolia")
		productLink := strings.TrimSpace(productLinkElement.AttrOr("href", ""))
		fullProductName := strings.TrimSpace(productLinkElement.Text())
		productSku := selection.Find(".ck-product-box-sku").Text()
		productPrice := selection.Find(".price__amount").Text()
		productPrice = strings.ReplaceAll(productPrice, " ", "")
		productPrice = strings.ReplaceAll(productPrice, "€", "")
		productPrice = strings.ReplaceAll(productPrice, "\u00a0", "")
		productPrice = strings.ReplaceAll(productPrice, ",", ".")

		parsedPrice, _ := strconv.ParseFloat(productPrice, 64)
		products = append(products, models.NewVendorProduct(fullProductName, parsedPrice, productLink, "https://www.globaldata.pt/", productSku))
	})
}

func QueryProduct(productName string) string {
	url := GlobalDataUrl + "/?query=" + strings.ReplaceAll(productName, " ", "%2520")
	fmt.Println(url)

	data, _ := os.ReadFile("./globalDataSearchPage.html")
	html := string(data)

	// browser := rod.New().MustConnect()
	// defer browser.Close()

	// page := browser.MustPage(url)
	// page.MustWaitStable()

	// html, err := page.HTML()

	// if err != nil {
	// 	panic(err)
	// }

	ParseQueryPage(html)

	return "test"
}
