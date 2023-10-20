package models

import (
	"fmt"
	"price-tracker/db"
	"time"
)

type Vendor int

const (
	GlobalData Vendor = iota
	PCDiga
)

type VendorEntry struct {
	Id           int
	UniversalId  int
	Price        float64
	Url          string
	Vendor       Vendor
	SKU          string
	FullName     string
	LastUpdated  int64
	Availability Availability
	History      []ProductHistory
}

type ProductHistory struct {
	Id            int
	VendorEntryId int
	Price         float64
	Availability  Availability
	UpdatedAt     int64
}

func NewVendorProduct(fullName string, price float64, url string, vendor Vendor, sku string, availability Availability) VendorEntry {
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

func UpdateVendorEntry(newPrice float64, newAvailability Availability, sku string, vendor Vendor) error {
	product, err := GetVendorProductEntry(sku, vendor)

	if product.Availability != newAvailability || product.Price != newPrice {
		statement, _ := db.GetDb().Prepare("INSERT INTO productHistory (vendorEntryId, price, availability, updatedAt) VALUES (?, ?, ?, ?)")
		defer statement.Close()
		_, insertionError := statement.Exec(product.Id, product.Price, product.Availability, product.LastUpdated)

		if insertionError != nil {
			fmt.Println(insertionError)
		}

		_, updateError := db.GetDb().Exec("UPDATE vendorEntries SET price=?,availability=?,lastUpdated=? WHERE sku=? AND vendor=?", newPrice, newAvailability, time.Now().Unix(), sku, vendor)

		if updateError != nil {
			return updateError
		}

	}

	return err
}

func GetVendorProductEntry(sku string, vendor Vendor) (VendorEntry, error) {
	row := db.GetDb().QueryRow("SELECT id, universalId, fullName, price, url, vendor, sku, availability, lastUpdated FROM vendorEntries WHERE sku = ? AND vendor = ?", sku, vendor)
	var result VendorEntry
	err := row.Scan(&result.Id, &result.UniversalId, &result.FullName, &result.Price, &result.Url, &result.Vendor, &result.SKU, &result.Availability, &result.LastUpdated)

	if err != nil {
		fmt.Println("Error fetching product ")
		return result, err
	}

	return result, nil
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

func DoesVendorProductExist(sku string, vendor Vendor) bool {
	row := db.GetDb().QueryRow("SELECT id FROM vendorEntries WHERE sku = ? AND vendor = ?", sku, vendor)
	var result int
	err := row.Scan(&result)

	if err != nil {
		return false
	}

	return true
}

func GetProductHistory(productId int) []ProductHistory {
	rows, err := db.GetDb().Query("SELECT id, vendorEntryId, price, availability, updatedAt FROM productHistory WHERE vendorEntryId = ?", productId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	history := make([]ProductHistory, 0)

	for rows.Next() {
		var productHistory ProductHistory

		err := rows.Scan(&productHistory.Id, &productHistory.VendorEntryId, &productHistory.Price, &productHistory.Availability, &productHistory.UpdatedAt)

		if err != nil {
			panic(err)
		}

		history = append(history, productHistory)
	}

	return history
}