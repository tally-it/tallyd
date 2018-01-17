package contract

import "github.com/jmoiron/sqlx/types"


// swagger:model
type UserSlice []User

// swagger:model
type User struct {
	UserID   int           `json:"userID" db:"user_id"`
	Name     string        `json:"name" db:"name"`
	Email    string        `json:"email" db:"email"`
	IsActive bool          `json:"active"`
	IsAdmin  types.BitBool `json:"isAdmin" db:"is_admin"`
	Balance  float64       `json:"balance" db:"balance"`
}

// swagger:model
type AddUserRequestBody struct {
	// required: true
	Name     string `json:"name"`
	// required: true
	Email    string `json:"email"`
	// required: true
	Password string `json:"password"`
}

// swagger:model
type LoginRequestBody struct {
	Name     string `json:"name" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

type AddUserResponseBody struct {
	UserID int `json:"userID"`
}
