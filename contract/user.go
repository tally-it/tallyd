package contract

import "github.com/jmoiron/sqlx/types"

type User struct {
	UserID   int           `json:"userID" db:"user_id"`
	Name     string        `json:"name" db:"name"`
	Email    string        `json:"email" db:"email"`
	IsActive bool          `json:"active"`
	IsAdmin  types.BitBool `json:"isAdmin" db:"is_admin"`
	Balance  float64       `json:"balance" db:"balance"`
}

type AddUserRequestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestBody struct {
	Name     string `json:"name" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

type AddUserResponseBody struct {
	UserID int `json:"userID"`
}
