package contract

type ChangeBalanceRequestBody struct {
	UserID int     `json:"userID" db:"user_id"`
	SKU    int     `json:"sku" db:"sku_id"`
	Value  float64 `json:"value" db:"value"`
	Tag    string  `json:"tag" db:"tag"`
}
