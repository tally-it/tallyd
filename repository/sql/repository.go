package sql

import (
	"context"

	"github.com/marove2000/hack-and-pay/config"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/repository/sql/migration"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/marove2000/hack-and-pay/log"
	sqlmigrate "github.com/rubenv/sql-migrate"
)

var pkgLogger = log.New("sql")

type Mysql struct {
	db *sqlx.DB
}

func New(conf *config.Mysql) (*Mysql, error) {
	logger := pkgLogger.ForFunc(context.Background(), "New")

	db, err := sqlx.Connect("mysql", buildConnectionString(conf))
	if err != nil {
		logger.WithError(err).Error("db connection failed")
		return nil, errors.InternalServerError("db connection failed", err)
	}

	n, err := sqlmigrate.Exec(db.DB, "mysql", &sqlmigrate.AssetMigrationSource{
		Asset:    migration.Asset,
		AssetDir: migration.AssetDir,
		Dir:      "resources",
	}, sqlmigrate.Up)
	if err != nil {
		logger.WithError(err).Error("db migration failed")
		return nil, errors.InternalServerError("db migration failed", err)
	}

	if n > 0 {
		logger.Infof("applied %d migrations", n)
	} else {
		logger.Info("db is up to date")
	}

	return &Mysql{
		db: db,
	}, nil
}

func buildConnectionString(conf *config.Mysql) string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?parseTime=true"
}
