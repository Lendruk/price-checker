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

	products := make([]models.VendorProduct, 0)
	document.Find("article.ck-product-box").Each(func(i int, selection *goquery.Selection) {
		productLinkElement := selection.Find(".js-gtm-product-link-algolia")
		productLink := strings.TrimSpace(productLinkElement.AttrOr("href", ""))
		fullProductName := strings.TrimSpace(productLinkElement.Text())
		productSku := selection.Find(".ck-product-box-sku").Text()
		productPrice, _ := utils.FormatPrice(selection.Find(".price__amount").Text())
		products = append(products, models.NewVendorProduct(fullProductName, productPrice, productLink, "https://www.globaldata.pt/", productSku, models.InStock))
	})

	for _, v := range products {
		if models.DoesProductExist(v.SKU, v.Vendor) == false {
			models.InsertProduct(v)
		}
	}
}

func QueryProduct(productName string, browser *rod.Browser) string {
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

	return "test"
}
