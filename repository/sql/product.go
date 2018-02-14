package sql

import (
	"context"
	"database/sql"

	"github.com/tally-it/tallyd/contract"
	"github.com/tally-it/tallyd/errors"
	"github.com/tally-it/tallyd/repository/sql/models"

	"gopkg.in/nullbio/null.v6"
	"github.com/vattle/sqlboiler/queries/qm"
	"time"
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
	logger := pkgLogger.ForFunc(ctx, "AddProduct")
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

func (m *Mysql) EditProduct(ctx context.Context, sku int, r contract.AddProductRequestBody) (error) {
	logger := pkgLogger.ForFunc(ctx, "EditProduct")
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

	// check if sku is existing
	c, err := models.Products(tx, qm.Where("SKU_id=?", sku)).Exists()
	if err != nil {
		logger.WithError(err).Error("failed to fetch product from db")
		return errors.InternalServerError("db error", err)
	}
	if c == false {
		logger.WithField("sku", sku).Warn("SKU does not exist")
		return errors.BadRequest("SKU does not exists")
	}

	// TODO: Check if nothing has changed
	// insert product
	err = product.Insert(tx)
	if err != nil {
		logger.WithError(err).Error("failed to insert product")
		return errors.InternalServerError("db error", err)
	}

	// update old SKUs if SKU is set in body and not equal old
	if sku != r.SKU && r.SKU != 0 {
		// check if new sku is already taken
		c, err = models.Products(tx, qm.Where("SKU_id=?", r.SKU)).Exists()
		if err != nil {
			logger.WithError(err).Error("failed to fetch product from db")
			return errors.InternalServerError("db error", err)
		}

		if c == true {
			logger.WithField("sku", r.SKU).Warn("SKU is already taken")
			return errors.BadRequest("SKU is already taken")
		}
		// get skus to update
		skuUpdateProducts, err := models.Products(tx, qm.Where("SKU_id=?", sku)).All()

		err = skuUpdateProducts.UpdateAll(tx, models.M{"SKU_id": r.SKU})
		//err = models.Products(tx, qm.Where("SKU_id=?", sku)).UpdateAll(models.M{"SKU_id": r.SKU})
		if err != nil {
			logger.WithError(err).Error("failed to update old skus")
			return errors.InternalServerError("db error", err)
		}
	}

	//TODO add category map

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit")
		return errors.InternalServerError("db error", err)
	}

	return nil
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

func (m *Mysql) DeleteProduct(ctx context.Context, sku int) (error) {
	logger := pkgLogger.ForFunc(ctx, "DeleteProduct")
	logger.Debug("enter repo")

	// check if sku is existing
	dbProduct, err := m.GetProductBySKU(ctx, sku)
	if err != nil {
		logger.WithField("sku", sku).WithError(err).Warn("failed to find product")
		return errors.BadRequest("failed to find product")
	}

	product := models.Product{
		SKUID:        dbProduct.SKU,
		Name:         dbProduct.Name,
		GTIN:         null.StringFrom(dbProduct.GTIN),
		Price:        dbProduct.Price.String(),
		Quantity:     null.StringFrom(dbProduct.Quantity.String()),
		QuantityUnit: null.StringFrom(dbProduct.QuantityUnit),
		DeletedAt: null.TimeFrom(time.Now()),
	}

	if dbProduct.Visibility == true {
		product.IsVisible = boolToString(true)
	} else {
		product.IsVisible = boolToString(false)
	}

	// delete product
	err = product.Insert(m.db)
	if err != nil {
		logger.WithError(err).Error("failed to delete product")
		return errors.InternalServerError("db error", err)
	}

	return nil
}

func (m *Mysql) ChangeStock(ctx context.Context, r contract.ChangeStockRequestBody) (error) {
	logger := pkgLogger.ForFunc(ctx, "ChangeProduct")
	logger.Debug("enter repo")


	// check if sku is existing
	_, err := m.GetProductBySKU(ctx, r.SKU)
	if err != nil {
		logger.WithField("sku", r.SKU).WithError(err).Warn("failed to find product")
		return errors.BadRequest("failed to find product")
	}

	stock := models.Stock{
		SKUID: r.SKU,
		UserID: null.IntFrom(r.UserID),
		Quantity: r.Quantity,
	}

	// update stock
	err = stock.Insert(m.db)
	if err != nil {
		logger.WithField("sku", r.SKU).WithError(err).Error("failed to update stock")
		return errors.InternalServerError("db error", err)
	}

	return nil
}