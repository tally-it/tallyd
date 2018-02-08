package contract

import "github.com/shopspring/decimal"

type ChangeBalanceRequestBody struct {
	UserID int             `json:"userID" db:"user_id" validate:"nonzero"`
	SKU    int             `json:"sku" db:"sku_id"`
	Value  decimal.Decimal `json:"value" db:"value" validate:"nonzero"`
	Tag    string          `json:"tag" db:"tag"`
}
