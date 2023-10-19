package models

import (
	"fmt"
	"price-tracker/db"
	"time"
)

type Product struct {
	ID            int
	UniversalName string
	SKU           string
	Vendors       []VendorProduct
}

type VendorProduct struct {
	Price       float64
	Url         string
	Vendor      string
	SKU         string
	FullName    string
	LastUpdated int64
}

func NewVendorProduct(fullName string, price float64, url string, vendor string, sku string) VendorProduct {
	return VendorProduct{
		FullName:    fullName,
		Price:       price,
		Url:         url,
		Vendor:      vendor,
		SKU:         sku,
		LastUpdated: time.Now().Unix(),
	}
}

func DoesProductExist(sku string) bool {
	row := db.GetDb().QueryRow("SELECT id FROM vendorProducts WHERE sku = ?", sku)
	var result int
	err := row.Scan(&result)

	if err != nil {
		return false
	}

	return true
}

func InsertProduct(product VendorProduct) {
	fmt.Println("Inserting product ", product.FullName, product.SKU)
	statement, _ := db.GetDb().Prepare("INSERT INTO vendorProducts (fullName, price, url, vendor, sku, lastUpdated) VALUES (?, ?, ?, ?, ?, ?)")
	defer statement.Close()
	_, err := statement.Exec(product.FullName, product.Price, product.Url, product.Vendor, product.SKU, product.LastUpdated)

	if err != nil {
		panic(err)
	}
}
