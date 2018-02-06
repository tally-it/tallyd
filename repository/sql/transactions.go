package sql

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/models"
	"gopkg.in/nullbio/null.v6"
	"github.com/vattle/sqlboiler/queries/qm"
)

func (m *Mysql) GetUsersWithBalance(ctx context.Context) ([]*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetUsersWithBalance")
	logger.Debug("enter repo")

	users := []*contract.User{}
	err := m.db.Select(&users, `
		SELECT users.user_id, 
			users.email, 
			users.name, 
			users.is_blocked,
			users.is_admin, 
			COALESCE(SUM(transactions.value), 0.00) AS 'balance' 
		FROM users 
		LEFT JOIN transactions 
		ON users.user_id = transactions.user_id 
		WHERE users.deleted_at IS NULL 
		GROUP BY users.user_id`)
	if err != nil {
		logger.WithError(err).Error("failed to fetch users from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return users, nil
}

func (m *Mysql) AddTransaction(ctx context.Context, body contract.ChangeBalanceRequestBody) (error) {
	logger := pkgLogger.ForFunc(ctx, "AddTransaction")
	logger.Debug("enter repository")

	var err error
	if body.SKU != 0 {
		// check if SKU ID is existing
		_, err = models.Products(m.db, qm.Where("SKU_id=?", body.SKU)).One()
		if err != nil && err != sql.ErrNoRows {
			logger.WithError(err).Error("failed to get product")
			return errors.InternalServerError("db error", err)
		}

		if err == sql.ErrNoRows {
			logger.WithError(err).Error(err)
			logger.Warn("failed to find product with sku ", body.SKU)
			body.SKU = 0
		}
	}

	transaction := models.Transaction{
		UserID: body.UserID,
		SKUID:  null.IntFrom(body.SKU),
		Value:  decimal.NewFromFloat(body.Value).String(),
		Tag:    null.StringFrom(body.Tag),
	}

	err = transaction.Insert(m.db)

	if err != nil {
		logger.WithError(err).Error("failed to insert transaction")
		return errors.InternalServerError("db error", err)
	}

	return nil
}
