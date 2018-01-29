package sql

import (
	"context"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
)

func (m *Mysql) GetProductsWithStock(ctx context.Context) ([]*contract.Product, error) {
	logger := pkgLogger.ForFunc(ctx, "GetProductsWithStock")
	logger.Debug("enter repo")

	products := []*contract.Product{}
	err := m.db.Select(&products, `
			SELECT p1.product_id, p1.SKU, p1.name, p1.GTIN, p1.quantity, p1.quantity_unit, p1.price, COALESCE(SUM(stock.quantity), 0) AS 'stock' 
			FROM products p1 LEFT JOIN stock ON p1.product_id = stock.product_id 
			WHERE p1.deleted_at IS NULL
			AND p1.is_visible=1 
			AND p1.product_id = (SELECT p2.product_id FROM products p2 WHERE p2.SKU = p1.SKU ORDER BY p2.product_id DESC LIMIT 1) 
			GROUP BY p1.product_id`)
	if err != nil {
		logger.WithError(err).Error("failed to fetch products from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return products, nil
}
