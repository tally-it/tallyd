package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/tally-it/tallyd/errors"
	"encoding/json"
	"github.com/go-validator/validator"
	"github.com/tally-it/tallyd/contract"
	"github.com/tally-it/tallyd/ctxutil"
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
	product, err := h.repo.GetProductByID(ctx, SKU)
	if err != nil {
		logger.WithError(err).Error("failed to get product data")
		return nil, err
	}

	return product, nil
}

func (h *Handler) addProduct(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "addProduct")
	logger.Debug("enter handler")

	//TODO check if exact the same product is already existing

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

func (h *Handler) editProduct(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "changeProduct")
	logger.Debug("enter handler")

	// TODO check if anything has changed

	// read id
	productIDStr := pathParams["productID"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		logger.WithError(err).Error("failed to parse product sku")
		return nil, errors.BadRequest(err.Error())
	}

	product := &contract.AddProductRequestBody{}

	// get body data
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(product)
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

	if ctxutil.GetAdminStatus(ctx) == true {
		err = h.repo.EditProduct(ctx, productID,product)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.Unauthorized()
	}

	return nil, nil
}

func (h *Handler) deleteProduct(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "deleteProduct")
	logger.Debug("enter handler")

	// read id
	productIDString := pathParams["productID"]
	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		logger.WithError(err).Error("failed to parse product productID")
		return nil, errors.BadRequest(err.Error())
	}

	if ctxutil.GetAdminStatus(ctx) == true {

		// update delete status
		err = h.repo.DeleteProduct(ctx, productID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.Unauthorized()
	}

	return nil, nil
}

func (h *Handler) changeStock(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "changeStock")
	logger.Debug("enter handler")

	// read id
	productIDStr := pathParams["productID"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		logger.WithError(err).Error("failed to parse product productID")
		return nil, errors.BadRequest(err.Error())
	}

	// get body data
	stock := &contract.ChangeStockRequestBody{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(stock)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	// insert user id
	stock.UserID = ctxutil.GetUserID(ctx)

	// insert sku
	stock.ProductID = productID

	// validate data
	if err = validator.Validate(stock); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	if ctxutil.GetAdminStatus(ctx) == true {
		// change stock
		err = h.repo.ChangeStock(ctx, *stock)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.Unauthorized()
	}

	return nil, nil
}