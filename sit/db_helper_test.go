package sit

import (
	"net"
	"testing"

	"github.com/tally-it/tallyd/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

type testDb struct {
	db *sqlx.DB
}

func NewDB(t *testing.T) *testDb {
	host := "mysql"
	_, err := net.LookupHost(host)
	if err != nil {
		host = "127.0.0.1"
	}

	db, err := sqlx.Connect("mysql", buildConnectionString(&config.Mysql{
		Password: "hackandpay",
		User:     "hackandpay",
		Port:     "3306",
		Host:     "127.0.0.1",
		Database: "hackandpay",
	}))
	require.NoError(t, err)

	return &testDb{
		db: db,
	}
}

func buildConnectionString(conf *config.Mysql) string {
	return conf.User + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Database + "?parseTime=true&multiStatements=true"
}

func (d *testDb) Clear() {

	d.db.MustExec(`
		SET FOREIGN_KEY_CHECKS = 0;
		TRUNCATE categories;
		TRUNCATE category_parent_map;
		TRUNCATE product_category_map;
		TRUNCATE products;
		TRUNCATE stock;
		TRUNCATE transactions;
		TRUNCATE user_auths;
		TRUNCATE users;
		SET FOREIGN_KEY_CHECKS = 1;
	`)
}
