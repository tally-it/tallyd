package contract

type ChangeStockRequestBody struct {
	UserID    int `json:"userID" db:"user_id"`
	ProductID int `json:"productID" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity" validate:"nonzero"`
}
