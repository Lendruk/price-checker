package parsers

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-rod/rod"
)

const GlobalDataUrl = "https://www.globaldata.pt/?query="

func ParseQueryPage(page *rod.Page) {

}

func QueryProduct(productName string) string {
	url := GlobalDataUrl + strings.ReplaceAll(productName, " ", "%2520")
	browser := rod.New().MustConnect()

	defer browser.Close()

	fmt.Println(url)
	page := browser.MustPage(url)
	page.MustWaitStable()

	html, err := page.HTML()

	if err != nil {
		panic(err)
	}

	fmt.Println(html)

	file, _ := os.Create(("test.html"))
	defer file.Close()

	file.WriteString(html)

	return "test"
}
