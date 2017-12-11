package v1

import (
	"database/sql"
	"github.com/marove2000/hack-and-pay/config/v1"
)

func PaymentTransfer(userID int64, value float64, paymentTag string) (err error) {
	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// prepare statements
	stmt, err := db.Prepare("INSERT INTO payment(userID, paymentValue, paymentTag) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	_ , err = stmt.Exec(userID, value, paymentTag)
	if err != nil {
		return err
	}
	tx.Commit()

	defer db.Close()

	return nil
}