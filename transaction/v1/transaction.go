package v1

import (
	"database/sql"
	"github.com/marove2000/hack-and-pay/config/v1"
)

type Transaction struct {
	UserID       	int64   `json:"userID"`
	ProductID     	int64	`json:"productID"`
	Value    		float64 `json:"value"`
	Tag 		string  `json:"tag"`
}

func PaymentTransfer(userID int64, transaction Transaction) (err error) {
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
	stmt, err := db.Prepare("INSERT INTO transactions(user_id, value, tag, product_id) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	_ , err = stmt.Exec(userID, transaction.Value, transaction.Tag, transaction.ProductID)
	if err != nil {
		return err
	}
	tx.Commit()

	defer db.Close()

	return nil
}