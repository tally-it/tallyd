package sql

import (
	"context"
	"database/sql"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/models"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx/types"
	sqlerror "github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"
)

type user struct {
	UserID    int             `json:"userID" db:"user_id"`
	Name      string          `json:"name" db:"name"`
	Email     sql.NullString  `json:"email" db:"email"`
	IsBlocked types.BitBool   `json:"isBlocked" db:"is_blocked"`
	IsAdmin   types.BitBool   `json:"isAdmin" db:"is_admin"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
}

func (m *Mysql) GetUserWithBalance(ctx context.Context, userID int) (*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetUsersWithBalance")
	logger.Debug("enter repo")

	user := &user{}
	err := m.db.Get(user, `
		SELECT users.user_id, 
			users.email, 
			users.name, 
			users.is_blocked,
			users.is_admin, 
			COALESCE(SUM(transactions.value), 0.00) AS 'balance' 
		FROM users 
		LEFT JOIN transactions 
		ON users.user_id = transactions.user_id 
		WHERE users.user_id = ?
		GROUP BY users.user_id`, userID)
	switch err {
	case nil: // ok

	case sql.ErrNoRows:
		logger.Warn("failed to find user")
		return nil, errors.NotFound("user not found")

	default:
		logger.WithError(err).Error("failed to fetch users from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return &contract.User{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email.String,
		IsBlocked: user.IsBlocked,
		IsAdmin:   user.IsAdmin,
		Balance:   user.Balance,
	}, nil
}

func (m *Mysql) GetUsersWithBalance(ctx context.Context) ([]*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetUsersWithBalance")
	logger.Debug("enter repo")

	var users []*contract.User
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
		GROUP BY users.user_id`)
	if err != nil {
		logger.WithError(err).Error("failed to fetch users from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return users, nil
}

func (m *Mysql) AddTransaction(ctx context.Context, r *contract.ChangeBalanceRequestBody) error {
	logger := pkgLogger.ForFunc(ctx, "AddTransaction")
	logger.Debug("enter repository")

	var sku null.Int
	if r.SKU != 0 {
		sku = null.IntFrom(r.SKU)
	}

	transaction := models.Transaction{
		UserID: null.IntFrom(r.UserID),
		SKUID:  sku,
		Value:  r.Value.String(),
		Tag:    null.StringFrom(r.Tag),
	}

	var err error
	if r.SKU != 0 {
		// check if SKU ID is existing
		_, err = models.Products(m.db, qm.Where("SKU_id=?", r.SKU)).One()
		if err != nil && err != sql.ErrNoRows {
			logger.WithError(err).Error("failed to get product")
			return errors.InternalServerError("db error", err)
		}

		if err == sql.ErrNoRows {
			logger.WithError(err).Error(err)
			logger.Warn("failed to find product with sku ", r.SKU)
			r.SKU = 0
		}

		product, err := m.GetProductBySKU(ctx, r.SKU)
		if err != nil {
			return err
		}
		transaction.Value = product.Price.String()

	}

	err = transaction.Insert(m.db)

	if err != nil {
		sqlerr, ok := sqlerror.Cause(err).(*mysql.MySQLError)
		if !ok {
			logger.WithError(err).Error("failed to insert transaction")
			return errors.InternalServerError("db error", err)
		}

		switch sqlerr.Number {
		case 1452:
			logger.WithField("sku", r.SKU).Warn("sku not found")
			return errors.NotFound("sku not found")
		default:
			logger.WithError(err).Error("failed to insert transaction")
			return errors.InternalServerError("db error", err)
		}

		logger.WithError(err).Error("failed to insert transaction")
		return errors.InternalServerError("db error", err)
	}

	return nil
}
