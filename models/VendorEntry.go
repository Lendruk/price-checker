package models

import (
	"price-tracker/db"
	"time"
)

type VendorEntry struct {
	Id           int
	UniversalId  int
	Price        float64
	Url          string
	Vendor       string
	SKU          string
	FullName     string
	LastUpdated  int64
	Availability Availability
	History      []ProductHistory
}

type ProductHistory struct {
	Id           int
	ProductId    int
	Price        float64
	Availability Availability
	UpdatedAt    int64
}

func NewVendorProduct(fullName string, price float64, url string, vendor string, sku string, availability Availability) VendorEntry {
	return VendorEntry{
		FullName:     fullName,
		Price:        price,
		Url:          url,
		Vendor:       vendor,
		SKU:          sku,
		Availability: availability,
		LastUpdated:  time.Now().Unix(),
	}
}

func GetVendorEntries(sku string) []VendorEntry {
	row := db.GetDb().QueryRow("SELECT id FROM products WHERE sku = ?", sku)
	var result int
	err := row.Scan(&result)

	if err != nil {
		return []VendorEntry{}
	}

	return GetVendorEntriesByUniversalId(result)
}

func GetAllVendorEntries() []VendorEntry {
	rows, err := db.GetDb().Query("SELECT id, fullName, price, url, vendor, sku, availability, lastUpdated FROM vendorEntries")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	products := make([]VendorEntry, 0)
	for rows.Next() {
		var product VendorEntry

		err := rows.Scan(&product.Id, &product.FullName, &product.Price, &product.Url, &product.Vendor, &product.SKU, &product.Availability, &product.LastUpdated)

		product.History = GetProductHistory(product.Id)
		if err != nil {
			panic(err)
		}

		products = append(products, product)
	}

	return products
}

func GetVendorEntriesByUniversalId(universalId int) []VendorEntry {
	rows, err := db.GetDb().Query("SELECT id, fullName, price, url, vendor, sku, availability, lastUpdated FROM vendorEntries WHERE universalId = ?", universalId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	products := make([]VendorEntry, 0)
	for rows.Next() {
		var product VendorEntry

		err := rows.Scan(&product.Id, &product.FullName, &product.Price, &product.Url, &product.Vendor, &product.SKU, &product.Availability, &product.LastUpdated)

		product.History = GetProductHistory(product.Id)
		if err != nil {
			panic(err)
		}

		products = append(products, product)
	}

	return products
}

func DoesVendorProductExist(sku string, vendor string) bool {
	row := db.GetDb().QueryRow("SELECT id FROM vendorEntries WHERE sku = ? AND vendor = ?", sku, vendor)
	var result int
	err := row.Scan(&result)

	if err != nil {
		return false
	}

	return true
}
