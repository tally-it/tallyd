package sql

import (
	"context"
	"database/sql"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/models"

	"github.com/vattle/sqlboiler/queries/qm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/nullbio/null.v6"
)

func (m *Mysql) AddLocalUser(ctx context.Context, name, email, password string) (userID int, err error) {
	logger := pkgLogger.ForFunc(ctx, "AddLocalUser")
	logger.Debug("enter repository")

	tx, err := m.db.Beginx()
	defer func() {
		errr := tx.Rollback()
		if errr != nil && errr != sql.ErrTxDone {
			logger.WithError(errr).Error("failed to roll back tx")
			err = errors.InternalServerError("db error", errr)
		}
	}()

	var nullMail null.String
	if email != "" {
		nullMail = null.StringFrom(email)
	}

	usr := models.User{
		Name:  name,
		Email: nullMail,
	}

	// add user
	err = usr.Insert(tx, models.UserColumns.Name, models.UserColumns.Email)
	if err != nil {
		logger.WithError(err).Error("failed to insert user")
		return 0, errors.InternalServerError("db error", err)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithError(err).Error("failed to create password hash")
		return 0, errors.InternalServerError("password hash error", err)
	}

	auth := models.UserAuth{
		Method: string(contract.AuthTypePasswd),
		UserID: usr.UserID,
		Value:  null.Bytes{Bytes: hashedPassword, Valid: true},
	}

	err = auth.Insert(tx, models.UserAuthColumns.UserID, models.UserAuthColumns.Method, models.UserAuthColumns.Value)
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

func (m *Mysql) AddLDAPUser(ctx context.Context, name, email string) (userID int, err error) {
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

	var nullMail null.String
	if email != "" {
		nullMail = null.StringFrom(email)
	}

	usr := models.User{
		Name:  name,
		Email: nullMail,
	}

	// add user
	err = usr.Insert(tx, models.UserColumns.Name, models.UserColumns.Email)
	if err != nil {
		logger.WithError(err).Error("failed to insert user")
		return 0, errors.InternalServerError("db error", err)
	}

	auth := models.UserAuth{
		Method: string(contract.AuthTypeLDAP),
		UserID: usr.UserID,
	}

	err = auth.Insert(tx, models.UserAuthColumns.UserID, models.UserAuthColumns.Method)
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

func (m *Mysql) GetPublicUserDataByUserID(ctx context.Context, userID int) (*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetPublicUserDataByUserID")
	logger.Debug("enter repository")

	user, err := models.Users(m.db, qm.Where("user_id=?", userID)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("userID", userID).Warn("failed to find user")
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

func (m *Mysql) GetPublicUserDataByUserName(ctx context.Context, name string) (*contract.User, error) {
	logger := pkgLogger.ForFunc(ctx, "GetPublicUserDataByUserName")
	logger.Debug("enter repository")

	user, err := models.Users(m.db, qm.Where("name=?", name)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithField("name", name).Warn("failed to find user")
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

func (m *Mysql) GetUserCount(ctx context.Context) (int64, error) {
	logger := pkgLogger.ForFunc(ctx, "GetUserCount")
	logger.Debug("enter repository")

	userCount, err := models.Users(m.db).Count()
	if err != nil {
		logger.WithError(err).Error("failed to count users")
		return 0, errors.InternalServerError("db error", err)
	}

	return userCount, nil
}

func (m *Mysql) Login(ctx context.Context, name, pass string) error {
	logger := pkgLogger.ForFunc(ctx, "Login")
	logger.Debug("enter repository")

	// get user id by name
	user, err := models.Users(m.db, qm.Where("name=?", name)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("failed to find user")
			return errors.NotFound("user not found")
		}

		logger.WithError(err).Error("failed to fetch users from db")
		return errors.InternalServerError("db error", err)
	}

	// get user_auth data by id
	auth, err := models.UserAuths(m.db, qm.Where("user_id=?", user.UserID)).One()
	if err != nil {
		logger.WithError(err).Error("failed to fetch users auth from db")
		return errors.InternalServerError("db error", err)
	}

	// check auth value
	err = bcrypt.CompareHashAndPassword([]byte(auth.Value.Bytes), []byte(pass))
	if err != nil {
		logger.WithError(err).Error("wrong password")
		return errors.Unauthorized()
	}

	return nil
}
