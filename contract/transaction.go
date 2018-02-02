package contract

type ChangeBalanceRequestBody struct {
	UserID    int64   `json:"userID" db:"user_id"`
	ProductID int64   `json:"productID" db:"product_id"`
	Value     float64 `json:"value" db:"value"`
	Tag       string  `json:"tag" db:"tag"`
}
