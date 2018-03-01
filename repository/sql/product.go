package sql

import (
	"context"
	"database/sql"

	"time"

	"github.com/tally-it/tallyd/contract"
	"github.com/tally-it/tallyd/errors"
	"github.com/tally-it/tallyd/repository/sql/models"

	"gopkg.in/nullbio/null.v6"
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

	// start transaction
	tx, err := m.db.Beginx()
	defer func() {
		errr := tx.Rollback()
		if errr != nil && errr != sql.ErrTxDone {
			logger.WithError(errr).Error("failed to roll back tx")
			err = errors.InternalServerError("db error", errr)
		}
	}()

	//add product to generate unique product id
	product := models.Product{}
	err = product.Insert(tx)
	if err != nil {
		logger.WithError(err).Error("failed to insert product")
		return 0, errors.InternalServerError("db error", err)
	}

	productVersion := models.ProductVersion{
		ProductID:    product.ProductID,
		Name:         r.Name,
		GTIN:         null.StringFrom(r.GTIN),
		Price:        r.Price.String(),
		Quantity:     null.StringFrom(r.Quantity.String()),
		QuantityUnit: null.StringFrom(r.QuantityUnit),
		IsVisible:    boolToString(bool(r.Visibility)),
	}

	err = productVersion.Insert(tx)
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

func (m *Mysql) EditProduct(ctx context.Context, productID int, r *contract.AddProductRequestBody) (error) {
	logger := pkgLogger.ForFunc(ctx, "EditProduct")
	logger.Debug("enter repo")

	// start transaction
	tx, err := m.db.Beginx()
	defer func() {
		errr := tx.Rollback()
		if errr != nil && errr != sql.ErrTxDone {
			logger.WithError(errr).Error("failed to roll back tx")
			err = errors.InternalServerError("db error", errr)
		}
	}()

	// check if productID is existing
	c, err := models.ProductExists(tx, productID)
	if err != nil {
		logger.WithError(err).Error("failed to fetch product from db")
		return errors.InternalServerError("db error", err)
	}
	if c == false {
		logger.WithField("productID", productID).Warn("productID does not exist")
		return errors.BadRequest("productID does not exists")
	}

	// TODO: Check if nothing has changed
	// schau hier, dass du die "check if productID is existing" gleich umbauen kannst auf get productversion, wenn da nichts zurück kommt existiert auch das produkt nicht

	productVersion := models.ProductVersion{
		ProductID:    productID,
		Name:         r.Name,
		GTIN:         null.StringFrom(r.GTIN),
		Price:        r.Price.String(),
		Quantity:     null.StringFrom(r.Quantity.String()),
		QuantityUnit: null.StringFrom(r.QuantityUnit),
		IsVisible:    boolToString(bool(r.Visibility)),
	}

	// insert product
	err = productVersion.Insert(tx)
	if err != nil {
		logger.WithError(err).Error("failed to insert product")
		return errors.InternalServerError("db error", err)
	}

	//TODO add category map

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit")
		return errors.InternalServerError("db error", err)
	}

	return nil
}

func (m *Mysql) GetProductByID(ctx context.Context, productID int) (*contract.Product, error) {
	logger := pkgLogger.ForFunc(ctx, "GetProductByID")
	logger.Debug("enter repo")

	product := new(contract.Product)
	err := m.db.Get(product, `
			SELECT 
				p1.product_id,
				p1.name, 
				p1.GTIN, 
				p1.quantity, 
				p1.quantity_unit, 
				p1.price, 
				p1.is_visible,
				COALESCE(SUM(stock.quantity), 0) AS 'stock' 
			FROM product_versions p1 LEFT JOIN stock ON p1.product_id = stock.product_id
			WHERE p1.deleted_at IS NULL
			AND p1.product_id = (SELECT p2.product_id FROM product_versions p2 WHERE p2.product_id = p1.product_id ORDER BY p2.product_version_id DESC LIMIT 1)
			AND p1.product_id = ?
			GROUP BY p1.product_version_id`, productID)

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

func (m *Mysql) DeleteProduct(ctx context.Context, productID int) (error) {
	logger := pkgLogger.ForFunc(ctx, "DeleteProduct")
	logger.Debug("enter repo")

	// check if productID is existing
	dbProduct, err := m.GetProductByID(ctx, productID)
	if err != nil {
		logger.WithField("productID", productID).WithError(err).Warn("failed to find product")
		return errors.BadRequest("failed to find product")
	}

	product := models.ProductVersion{
		ProductID:    dbProduct.ProductID,
		Name:         dbProduct.Name,
		GTIN:         null.StringFrom(dbProduct.GTIN),
		Price:        dbProduct.Price.String(),
		Quantity:     null.StringFrom(dbProduct.Quantity.String()),
		QuantityUnit: null.StringFrom(dbProduct.QuantityUnit),
		DeletedAt:    null.TimeFrom(time.Now()),
		IsVisible:    boolToString(bool(dbProduct.Visibility)),
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
	logger := pkgLogger.ForFunc(ctx, "ChangeStock")
	logger.Debug("enter repo")

	// check if sku is existing
	_, err := m.GetProductByID(ctx, r.ProductID)
	if err != nil {
		logger.WithField("productID", r.ProductID).WithError(err).Warn("failed to find product")
		return errors.BadRequest("failed to find product")
	}

	stock := models.Stock{
		ProductID: r.ProductID,
		UserID:    null.IntFrom(r.UserID),
		Quantity:  r.Quantity,
	}

	// update stock
	err = stock.Insert(m.db)
	if err != nil {
		logger.WithField("productID", r.ProductID).WithError(err).Error("failed to update stock")
		return errors.InternalServerError("db error", err)
	}

	return nil
}
