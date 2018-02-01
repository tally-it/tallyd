package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/marove2000/hack-and-pay/errors"
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

func (h *Handler) productDetail(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "productDetail")
	logger.Debug("enter handler")

	// read id
	sku := pathParams["sku"]
	SKU, err := strconv.Atoi(sku)
	if err != nil {
		logger.WithError(err).Error("failed to parse product sku")
		return nil, errors.BadRequest(err.Error())
	}

	// get all product data
	product, err := h.repo.GetProductBySKU(ctx, SKU)
	if err != nil {
		logger.WithError(err).Error("failed to get product data")
		return nil, err
	}

	return product, nil
}