package contract

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/shopspring/decimal"
)

type Product struct {
	ProductID    int             `json:"productID" db:"product_id"`
	SKU          int             `json:"SKU" db:"SKU_id"`
	Name         string          `json:"Name" db:"name"`
	GTIN         int64           `json:"GTIN" db:"GTIN"`
	Price        decimal.Decimal `json:"price" db:"price"`
	Visibility   types.BitBool   `json:"visibility" db:"is_visible"`
	Category     []int64         `json:"category"`
	Quantity     float64         `json:"quantity" db:"quantity"`
	QuantityUnit string          `json:"quantityUnit" db:"quantity_unit"`
	Stock        int             `json:"stock" db:"stock"`
}

type ProductReturnBody struct {
	*Product
	TimeAdded   string `json:"TimeAdded" db:"added_ad"`
	TimeChanged string `json:"TimeChange" db:"updated_ad"`
	TimeDeleted string `json:"TimeDeleted" db:"deleted_ad"`
}

// swagger:model
type AddProductRequestBody struct {
	// required: true
	SKU int `json:"SKU" db:"SKU_id" validate:"nonzero"`
	// required: true
	Name string `json:"Name" db:"name" validate:"nonzero"`
	// required: false
	GTIN string `json:"GTIN" db:"GTIN"`
	// required: true
	Price decimal.Decimal `json:"price" db:"price" validate:"nonzero"`
	// required: false
	Visibility types.BitBool `json:"visibility" db:"is_visible"`
	// required: false
	Category []int64 `json:"category"`
	// required: false
	Quantity decimal.Decimal `json:"quantity" db:"quantity"`
	// required: false
	QuantityUnit string `json:"quantityUnit" db:"quantity_unit"`
}
