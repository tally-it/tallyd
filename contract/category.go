package contract

import (
	"github.com/jmoiron/sqlx/types"
)

// swagger:model
type AddCategoryBody struct {
	// required: true
	Name string `json:"Name" db:"name" validate:"nonzero"`
	// required: false
	IsVisible types.BitBool `json:"isVisible" db:"is_visible"`
	// required: false
	IsActive types.BitBool `json:"isActive" db:"is_active"`
	// required: false
	IsRoot types.BitBool `json:"isRoot" db:"is_root"`
}
