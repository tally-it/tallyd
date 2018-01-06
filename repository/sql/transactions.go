package sql

import (
	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"

	"github.com/sirupsen/logrus"
)

func (m *Mysql) GetUsersWithBalance() ([]*contract.User, error) {
	logger := logrus.WithField("func", pkg+"Mysql.GetAllUsers")
	logger.Debug("enter repo")

	users := []*contract.User{}
	err := m.db.Select(&users, "SELECT users.user_id, users.email, users.name, users.is_admin, SUM(transactions.value) AS 'balance' FROM users LEFT JOIN transactions ON users.user_id = transactions.user_id WHERE users.deleted_at IS NULL GROUP BY users.user_id")
	if err != nil {
		logger.WithError(err).Error("failed to fetch users from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return users, nil
}
