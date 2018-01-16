package sql

import (
	"context"
	"database/sql"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/models"

	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"
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

func (m *Mysql) AddLDAPUser(ctx context.Context, body *contract.AddUserRequestBody) (userID int, err error) {
	logger := pkgLogger.ForFunc(ctx, "AddLDAPUser")
	logger.Debug("enter repository")

	tx, err := m.db.Beginx()
	defer func() {
		errr := tx.Rollback()
		if errr != nil && errr != sql.ErrTxDone {
			logger.WithError(errr).Error("failed to roll back tx")
			err = errors.InternalServerError("db error", errr)
		}
	}()

	usr := models.User{
		Name:  body.Name,
		Email: null.StringFrom(body.Email),
	}

	// add user
	err = usr.Insert(m.db, models.UserColumns.Name, models.UserColumns.Email)
	if err != nil {
		logger.WithError(err).Error("failed to insert user")
		return 0, errors.InternalServerError("db error", err)
	}

	auth := models.UserAuth{
		Method: string(contract.AuthTypeLDAP),
		UserID: usr.UserID,
	}

	err = auth.Insert(m.db, models.UserAuthColumns.UserID, models.UserAuthColumns.Method)
	if err != nil {
		logger.WithError(err).Error("failed to insert user auth")
		return 0, errors.InternalServerError("db error", err)
	}

	err = tx.Commit()
	if err != nil {
		logger.WithError(err).Error("failed to commit")
		return 0, errors.InternalServerError("db error", err)
	}

	return usr.UserID, err
}

func (m *Mysql) GetPublicUserDataByUserName(ctx context.Context, name string) (*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetPublicUserDataByUserName")
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
