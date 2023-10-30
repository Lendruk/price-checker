package models

import (
	"fmt"
	"price-tracker/db"
)

type Availability int

const (
	InStock Availability = iota
	OutOfStock
	ByOrder
	PreOrder
)

type Product struct {
	Id            int           `json:"id"`
	SKU           string        `json:"sku"`
	VendorEntries []VendorEntry `json:"vendorEntries"`
}

func GetProducts() []Product {
	rows, err := db.GetDb().Query("SELECT id, sku FROM products")
	defer rows.Close()

	products := make([]Product, 0)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var product Product

		err := rows.Scan(&product.Id, &product.SKU)

		product.VendorEntries = GetVendorEntriesByUniversalId(product.Id)
		if err != nil {
			panic(err)
		}

		products = append(products, product)
	}

	return products
}

func GetProductById(productId int) Product {
	row := db.GetDb().QueryRow("SELECT * FROM products WHERE id = ? ", productId)

	var product Product
	err := row.Scan(&product.Id, &product.SKU)

	if err != nil {
		panic(err)
	}

	product.VendorEntries = GetVendorEntries(product.SKU)

	return product
}

func GetProductId(sku string) (int, error) {
	row := db.GetDb().QueryRow("SELECT id FROM products WHERE sku = ?", sku)

	var id int
	err := row.Scan(&id)

	return id, err
}

func getOrCreateProduct(sku string) Product {
	row := db.GetDb().QueryRow("SELECT id FROM products WHERE sku = ?", sku)
	var result int
	err := row.Scan(&result)

	if err != nil {
		statement, _ := db.GetDb().Prepare("INSERT INTO products (sku) VALUES (?)")
		defer statement.Close()
		_, err := statement.Exec(sku)

		if err != nil {
			panic(err)
		}

		return getOrCreateProduct(sku)
	}

	return Product{Id: result, SKU: sku, VendorEntries: GetVendorEntriesByUniversalId(result)}
}

func InsertProduct(product VendorEntry) Product {
	fmt.Println("Inserting product ", product.FullName, product.SKU)
	// Universal Product
	universalProduct := getOrCreateProduct(product.SKU)

	// Vendor Product
	if DoesVendorProductExist(product.SKU, product.Vendor) == true {
		updated, updatedVendorEntry, _ := UpdateVendorEntry(product.Price, product.Availability, product.SKU, product.Vendor)

		if updated {
			updatedEntries := make([]VendorEntry, 0)

			for _, entry := range universalProduct.VendorEntries {
				if entry.SKU == updatedVendorEntry.SKU && entry.Vendor == updatedVendorEntry.Vendor {
					updatedEntries = append(updatedEntries, updatedVendorEntry)
				} else {
					updatedEntries = append(updatedEntries, entry)
				}
			}
		}
	} else {
		statement, _ := db.GetDb().Prepare("INSERT INTO vendorEntries (fullName, price, url, vendor, sku, lastUpdated, availability, universalId, productImageUrl) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer statement.Close()
		_, err := statement.Exec(product.FullName, product.Price, product.Url, product.Vendor, product.SKU, product.LastUpdated, product.Availability, universalProduct.Id, product.ProductImageUrl)

		if err != nil {
			panic(err)
		}

		universalProduct.VendorEntries = append(universalProduct.VendorEntries, product)
	}

	return universalProduct
}
