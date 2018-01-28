package sql

import (
	"context"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
)

func (m *Mysql) GetUsersWithBalance(ctx context.Context) ([]*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetUsersWithBalance")
	logger.Debug("enter repo")

	users := []*contract.User{}
	err := m.db.Select(&users, "SELECT users.user_id, users.email, users.name, users.is_admin, COALESCE(SUM(transactions.value), 0.00) AS 'balance' FROM users LEFT JOIN transactions ON users.user_id = transactions.user_id WHERE users.deleted_at IS NULL GROUP BY users.user_id")
	if err != nil {
		logger.WithError(err).Error("failed to fetch users from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return users, nil
}
