package sql

import (
	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/models"

	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/vattle/sqlboiler/queries/qm"
)

//func (m *Mysql) GetAllUsers() ([]*contract.User, error) {
//	logger := logrus.WithField("func", pkg+"Mysql.GetAllUsers")
//	logger.Debug("enter repository")
//
//	users, err := models.Users(m.db).All()
//	if err != nil {
//		logger.WithError(err).Error("failed to fetch users from db")
//		return nil, errors.InternalServerError("db error", err)
//	}
//
//	out := make([]*contract.User, len(users))
//	for i := range users {
//		out[i] = &contract.User{
//			IsAdmin: tinyToBool(users[i].IsAdmin),
//			Email:   users[i].Email.String,
//			Name:    users[i].Name,
//			UserID:  int(users[i].UserID),
//		}
//	}
//
//
//
//	return out, nil
//}

func (m *Mysql) GetPublicUserDataByUserName(name string) (*contract.User, error) {
	logger := logrus.WithField("func", pkg+"Mysql.GetAllUsers")
	logger.Debug("enter repository")

	user, err := models.Users(m.db, qm.Where("name=?", name)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("failed to find user")
			return nil, errors.NotFound("user not found")
		}

		logger.WithError(err).Error("failed to fetch users from db")
		return nil, errors.InternalServerError("db error", err)
	}

	return &contract.User{
		UserID: user.UserID,
		Name:   user.Name,
		Email:  user.Email.String,
	}, nil
}
