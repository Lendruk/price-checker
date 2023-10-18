package models

type Product struct {
	ID            int
	UniversalName string
	SKU           string
	Vendors       []VendorProduct
}

type VendorProduct struct {
	Price    float64
	Url      string
	Vendor   string
	SKU      string
	FullName string
}

func NewVendorProduct(fullName string, price float64, url string, vendor string, sku string) VendorProduct {
	return VendorProduct{
		FullName: fullName,
		Price:    price,
		Url:      url,
		Vendor:   vendor,
		SKU:      sku,
	}
}
