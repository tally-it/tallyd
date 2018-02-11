package contract

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/shopspring/decimal"
)

// swagger:model
type UserSlice []User

// swagger:model
type User struct {
	UserID    int             `json:"userID" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Email     string          `json:"email" db:"email"`
	IsBlocked types.BitBool   `json:"isBlocked" db:"is_blocked"`
	IsAdmin   types.BitBool   `json:"isAdmin" db:"is_admin"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
}

// swagger:model
type EditUserRequestBody struct {
	// required: true
	Name string `json:"name" validate:"nonzero"`
	// required: true
	Email string `json:"email"`
	// required: true
	IsBlocked types.BitBool `json:"isBlocked"`
	// required: true
	IsAdmin types.BitBool `json:"isAdmin"`
}

// swagger:model
type AddUserRequestBody struct {
	// required: true
	Name string `json:"name"`
	// required: true
	Email string `json:"email"`
	// required: true
	Password string `json:"password"`
}

// swagger:model
type LoginRequestBody struct {
	Name     string `json:"name" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

type LoginResponse struct {
	UserID int    `json:"userID"`
	JWT    string `json:"jwt"`
}

type AddUserResponseBody struct {
	UserID int `json:"userID"`
}
