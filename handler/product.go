package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/marove2000/hack-and-pay/errors"
	"encoding/json"
	"github.com/go-validator/validator"
	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/ctxutil"
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

func (h *Handler) addProduct(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "addProduct")
	logger.Debug("enter handler")

	product := &contract.AddProductRequestBody{}

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(product)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	// validate data
	if err = validator.Validate(product); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	var productID int
	if ctxutil.GetAdminStatus(ctx) == true {
		productID, err = h.repo.AddProduct(ctx, *product)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.Unauthorized()
	}

	return productID, nil
}