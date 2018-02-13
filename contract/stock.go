package contract

type ChangeStockRequestBody struct {
	UserID   int    `json:"userID" db:"user_id"`
	SKU      int    `json:"sku" db:"SKU_id"`
	Quantity int    `json:"quantity" db:"quantity" validate:"nonzero"`
}
