package errors

import (
	"fmt"
	"net/http"
)

// swagger:model error
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"detail"`
	Cause   error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("[%d] %s", e.Status, e.Details)
	}

	return fmt.Sprintf("[%d] %s: %s", e.Status, e.Details, e.Cause)
}

func NotFound(details string) error {
	return &Error{
		Status:  http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
		Details: details,
	}
}

func Conflict(details string, cause error) error {
	return &Error{
		Status:  http.StatusConflict,
		Message: http.StatusText(http.StatusConflict),
		Details: details,
		Cause:   cause,
	}
}

func InternalServerError(details string, cause error) error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Details: details,
		Cause:   cause,
	}
}

func Unauthorized() error {
	return &Error{
		Status:  http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
		Details: "not authorized",
	}
}

func Unauthenticated() error {
	return &Error{
		Status:  http.StatusForbidden,
		Message: http.StatusText(http.StatusForbidden),
		Details: "not authenticated",
	}
}


func BadRequest(details string) error {
	return &Error{
		Status:  http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
		Details: details,
	}
}

