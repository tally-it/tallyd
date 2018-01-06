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
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AddUserResponseBody struct {
	*User
	JWT string `json:"JWT"`
}
