package handler

import (
	"context"
	"net/http"

	"github.com/marove2000/hack-and-pay/errors"
	"encoding/json"
	"github.com/go-validator/validator"
	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/ctxutil"
)

func (h *Handler) addCategory(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "addCategory")
	logger.Debug("enter handler")

	category := &contract.AddCategoryBody{}

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(category)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	// validate data
	if err = validator.Validate(category); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	var productID int
	if ctxutil.GetAdminStatus(ctx) == true {
		productID, err = h.repo.AddCategory(ctx, *category)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.Unauthorized()
	}

	return productID, nil
}
