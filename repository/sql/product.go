package sql

import (
	"context"
	"database/sql"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/models"

	"gopkg.in/nullbio/null.v6"
	"github.com/vattle/sqlboiler/queries/qm"
)

func (m *Mysql) GetProductsWithStock(ctx context.Context) ([]*contract.Product, error) {
	logger := pkgLogger.ForFunc(ctx, "GetProductsWithStock")
	logger.Debug("enter repo")

	products := []*contract.Product{}
	err := m.db.Select(&products, `
			SELECT 
				p1.product_id, 
				p1.SKU_id, 
				p1.name, 
				p1.GTIN, 
				p1.quantity, 
				p1.quantity_unit,
				p1.price, 
				p1.is_visible,
				COALESCE(SUM(stock.quantity), 0) AS 'stock' 
			FROM products p1 LEFT JOIN stock ON p1.SKU_id = stock.SKU_id 
			WHERE p1.deleted_at IS NULL
			AND p1.is_visible=1 
			AND p1.product_id = (SELECT p2.product_id FROM products p2 WHERE p2.SKU_id = p1.SKU_id ORDER BY p2.product_id DESC LIMIT 1) 
			GROUP BY p1.product_id`)
	if err != nil {
		logger.WithError(err).Error("failed to fetch products from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return products, nil
}

func (m *Mysql) AddProduct(ctx context.Context, r contract.AddProductRequestBody) (int, error) {
	logger := pkgLogger.ForFunc(ctx, "GetProductsWithStock")
	logger.Debug("enter repo")

	product := models.Product{
		SKUID:        r.SKU,
		Name:         r.Name,
		GTIN:         null.StringFrom(r.GTIN),
		Price:        r.Price.String(),
		Quantity:     null.StringFrom(r.Quantity.String()),
		QuantityUnit: null.StringFrom(r.QuantityUnit),
	}

	if r.Visibility == true {
		product.IsVisible = boolToString(true)
	} else {
		product.IsVisible = boolToString(false)
	}

	// start transaction
	tx, err := m.db.Beginx()
	defer func() {
		errr := tx.Rollback()
		if errr != nil && errr != sql.ErrTxDone {
			logger.WithError(errr).Error("failed to roll back tx")
			err = errors.InternalServerError("db error", errr)
		}
	}()

	// check if product already exists
	c, err := models.Products(tx, qm.Where("SKU_id=?", r.SKU)).Exists()
	if err != nil {
		logger.WithError(err).Error("failed to fetch product from db")
		return 0, errors.InternalServerError("db error", err)
	}

	if c == true {
		logger.WithField("sku", r.SKU).Warn("product already exists")
		return 0, errors.BadRequest("SKU already exists")
	}

	err = product.Insert(tx)
	if err != nil {
		logger.WithError(err).Error("failed to insert product")
		return 0, errors.InternalServerError("db error", err)
	}

	//TODO add category map

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit")
		return 0, errors.InternalServerError("db error", err)
	}

	return product.ProductID, nil
}

func (m *Mysql) GetProductBySKU(ctx context.Context, SKU int) (*contract.Product, error) {
	logger := pkgLogger.ForFunc(ctx, "GetProductBySKU")
	logger.Debug("enter repo")

	product := new(contract.Product)
	err := m.db.Get(product, `
			SELECT 
				p1.product_id,
				p1.SKU_id, 
				p1.name, 
				p1.GTIN, 
				p1.quantity, 
				p1.quantity_unit, 
				p1.price, 
				p1.is_visible,
				COALESCE(SUM(stock.quantity), 0) AS 'stock' 
			FROM products p1 LEFT JOIN stock ON p1.SKU_id = stock.SKU_id
			WHERE p1.deleted_at IS NULL
			AND p1.is_visible=1 
			AND p1.product_id = (SELECT p2.product_id FROM products p2 WHERE p2.SKU_id = p1.SKU_id ORDER BY p2.product_id DESC LIMIT 1)
			AND p1.SKU_id = ?
			GROUP BY p1.product_id`, SKU)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("failed to find product")
			return nil, errors.NotFound("product not found")
		}

		logger.WithError(err).Error("failed to fetch product from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return product, nil
}
