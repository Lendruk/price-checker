package models

type Product struct {
	ID            int
	UniversalName string
	Prices        []Price
}

type Price struct {
	Vendor string
	Price  float64
}
