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
	PcComponentes
)

type VendorEntry struct {
	Id              int              `json:"id"`
	UniversalId     int              `json:"universalId"`
	Price           float64          `json:"price"`
	Url             string           `json:"url"`
	ProductImageUrl string           `json:"productImageUrl"`
	Vendor          Vendor           `json:"vendor"`
	SKU             string           `json:"sku"`
	FullName        string           `json:"fullName"`
	LastUpdated     int64            `json:"lastUpdated"`
	Availability    Availability     `json:"availability"`
	History         []ProductHistory `json:"history"`
}

type ProductHistory struct {
	Id            int          `json:"id"`
	VendorEntryId int          `json:"vendorEntryId"`
	Price         float64      `json:"price"`
	Availability  Availability `json:"availability"`
	UpdatedAt     int64        `json:"updatedAt"`
}

func NewVendorProduct(fullName string, price float64, url string, vendor Vendor, sku string, productImageUrl string, availability Availability) VendorEntry {
	return VendorEntry{
		FullName:        fullName,
		Price:           price,
		Url:             url,
		Vendor:          vendor,
		SKU:             sku,
		ProductImageUrl: productImageUrl,
		Availability:    availability,
		LastUpdated:     time.Now().Unix(),
	}
}

func UpdateVendorEntry(newPrice float64, newAvailability Availability, sku string, vendor Vendor) (bool, VendorEntry, error) {
	product, err := GetVendorProductEntry(sku, vendor)

	if product.Availability != newAvailability || product.Price != newPrice {
		statement, _ := db.GetDb().Prepare("INSERT INTO productHistory (vendorEntryId, price, availability, updatedAt) VALUES (?, ?, ?, ?)")
		defer statement.Close()
		insertionResult, insertionError := statement.Exec(product.Id, product.Price, product.Availability, product.LastUpdated)

		if insertionError != nil {
			fmt.Println(insertionError)
		}

		_, updateError := db.GetDb().Exec("UPDATE vendorEntries SET price=?,availability=?,lastUpdated=? WHERE sku=? AND vendor=?", newPrice, newAvailability, time.Now().Unix(), sku, vendor)
		newId, _ := insertionResult.LastInsertId()
		newHistory, _ := GetVendorHistoryEntryById(int(newId))
		product.History = append(product.History, newHistory)
		updatedProduct, _ := GetVendorProductEntry(sku, vendor)

		if updateError != nil {
			return false, VendorEntry{}, updateError
		} else {
			return true, updatedProduct, nil
		}

	}

	return false, VendorEntry{}, err
}

func GetVendorProductEntry(sku string, vendor Vendor) (VendorEntry, error) {
	row := db.GetDb().QueryRow("SELECT id, universalId, fullName, price, url, vendor, sku, availability, lastUpdated FROM vendorEntries WHERE sku = ? AND vendor = ?", sku, vendor)
	var result VendorEntry
	err := row.Scan(&result.Id, &result.UniversalId, &result.FullName, &result.Price, &result.Url, &result.Vendor, &result.SKU, &result.Availability, &result.LastUpdated)

	result.History = GetVendorEntryHistory(result.Id)
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

		product.History = GetVendorEntryHistory(product.Id)
		if err != nil {
			panic(err)
		}

		products = append(products, product)
	}

	return products
}

func GetVendorEntriesByUniversalId(universalId int) []VendorEntry {
	rows, err := db.GetDb().Query("SELECT id, fullName, price, url, vendor, sku, availability, productImageUrl, lastUpdated FROM vendorEntries WHERE universalId = ?", universalId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	products := make([]VendorEntry, 0)
	for rows.Next() {
		var product VendorEntry

		err := rows.Scan(&product.Id, &product.FullName, &product.Price, &product.Url, &product.Vendor, &product.SKU, &product.Availability, &product.ProductImageUrl, &product.LastUpdated)

		product.History = GetVendorEntryHistory(product.Id)
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

func GetVendorEntryHistory(productId int) []ProductHistory {
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

func GetVendorHistoryEntryById(historyId int) (ProductHistory, error) {
	row := db.GetDb().QueryRow("SELECT id, vendorEntryId, price, availability, updatedAt FROM productHistory WHERE id = ?", historyId)

	var productHistory ProductHistory

	err := row.Scan(&productHistory.Id, &productHistory.VendorEntryId, &productHistory.Price, &productHistory.Availability, &productHistory.UpdatedAt)

	if err != nil {
		return ProductHistory{}, err
	} else {
		return productHistory, nil
	}
}
