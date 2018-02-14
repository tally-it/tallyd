package main

import (
	"fmt"
	"time"

	"github.com/tally-it/tallyd/repository/sql/migration"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	sqlmigrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &prefixed.TextFormatter{
		ForceColors:     true,
		TimestampFormat: time.RFC1123,
	}

	db, err := sqlx.Connect("mysql", buildConnectionString("hackandpay", "hackandpay", "hackandpay", "127.0.0.1", 3306))
	if err != nil {
		logger.WithError(err).Fatal("db error")
	}

	n, err := sqlmigrate.Exec(db.DB, "mysql", &sqlmigrate.AssetMigrationSource{
		Asset:    migration.Asset,
		AssetDir: migration.AssetDir,
		Dir:      "resources",
	}, sqlmigrate.Up)
	if err != nil {
		logger.WithError(err).Fatal("migration failed")
	}

	if n > 0 {
		logger.Infof("applied %d migrations", n)
	} else {
		logger.Info("db is up to date")
	}
}

func buildConnectionString(user, pass, dbname, host string, port int) string {
	return user + ":" + pass + "@tcp(" + host + ":" + fmt.Sprint(port) + ")/" + dbname + "?parseTime=true"
}
