package utils

import (
	"strconv"
	"strings"
)

func FormatPrice(price string) (float64, error) {
	price = strings.ReplaceAll(price, " ", "")
	price = strings.ReplaceAll(price, "€", "")
	price = strings.ReplaceAll(price, "\u00a0", "")
	price = strings.ReplaceAll(price, ",", ".")

	return strconv.ParseFloat(price, 64)
}
