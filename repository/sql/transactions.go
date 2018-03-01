package sql

import (
	"context"
	"database/sql"

	"github.com/tally-it/tallyd/contract"
	"github.com/tally-it/tallyd/errors"
	"github.com/tally-it/tallyd/repository/sql/models"

	"github.com/go-sql-driver/mysql"
	sqlerror "github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"
)

func (m *Mysql) AddTransaction(ctx context.Context, r *contract.ChangeBalanceRequestBody) error {
	logger := pkgLogger.ForFunc(ctx, "AddTransaction")
	logger.Debug("enter repository")

	var productID null.Int
	if r.ProductID != 0 {
		productID = null.IntFrom(r.ProductID)
	}

	transaction := models.Transaction{
		UserID:    null.IntFrom(r.UserID),
		ProductID: productID,
		Value:     r.Value.String(),
		Tag:       null.StringFrom(r.Tag),
	}

	if r.ProductID != 0 {
		// check if ProductID is existing
		productVersion, err := models.ProductVersions(m.db,
			qm.Select(models.ProductVersionColumns.Price),
			qm.Where(models.ProductVersionColumns.ProductID+"=?", r.ProductID),
			qm.OrderBy(models.ProductVersionColumns.ProductVersionID+" DESC")).One()
		switch err{
		case sql.ErrNoRows:
			logger.WithField("productID", r.ProductID).Warn("failed to find product")
			return errors.NotFound("failed to find product")
		case nil:

		default:
			logger.WithError(err).Error("failed to get product")
			return errors.InternalServerError("db error", err)
		}

		transaction.Value = productVersion.Price

	}

	err := transaction.Insert(m.db)
	if err != nil {
		sqlerr, ok := sqlerror.Cause(err).(*mysql.MySQLError)
		if !ok {
			logger.WithError(err).Error("failed to insert transaction")
			return errors.InternalServerError("db error", err)
		}

		switch sqlerr.Number {
		case 1452:
			logger.WithField("productID", r.ProductID).Warn("productID not found")
			return errors.NotFound("productID not found")
		default:
			logger.WithError(err).Error("failed to insert transaction")
			return errors.InternalServerError("db error", err)
		}

		logger.WithError(err).Error("failed to insert transaction")
		return errors.InternalServerError("db error", err)
	}

	return nil
}
