package handler

import (
	"context"
	"net/http"
)

func (h *Handler) productIndex(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "productIndex")
	logger.Debug("enter handler")

	// get all user data
	products, err := h.repo.GetProductsWithStock(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}
