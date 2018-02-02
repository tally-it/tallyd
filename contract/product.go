package contract

import (
	"github.com/jmoiron/sqlx/types"
)

type Product struct {
	ProductID    int           `json:"productID" db:"product_id"`
	SKU          int           `json:"SKU" db:"SKU_id"`
	Name         string        `json:"Name" db:"name"`
	GTIN         int64         `json:"GTIN" db:"GTIN"`
	Price        float64       `json:"price" db:"price"`
	Visibility   types.BitBool `json:"visibility" db:"is_visible"`
	Category     []int64       `json:"category"`
	Quantity     float64       `json:"quantity" db:"quantity"`
	QuantityUnit string        `json:"quantityUnit" db:"quantity_unit"`
	Stock        int           `json:"stock" db:"stock"`
}

type ProductReturnBody struct {
	*Product
	TimeAdded   string `json:"TimeAdded" db:"added_ad"`
	TimeChanged string `json:"TimeChange" db:"updated_ad"`
	TimeDeleted string `json:"TimeDeleted" db:"deleted_ad"`
}
