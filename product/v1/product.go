package v1

import (
	"database/sql"
	"github.com/marove2000/hack-and-pay/config/v1"
)


type Product struct {
	ProductID				int64		`json:"productID,string"`
	ProductName				string		`json:"productName"`
	ProductGTIN				int64		`json:"productGTIN"`
	ProductPrice			float64		`json:"productPrice"`
	ProductVisibility		bool		`json:"productVisibility"`
	ProductCategory			[]int64		`json:"productCategory"`
	ProductQuantity			float64		`json:"productQuantity"`
	ProductQuantityUnit		string		`json:"productQuantityUnit"`
	ProductTimeAdded		string		`json:"productTimeAdded"`
	ProductTimeChanged		string		`json:"productTimeChange"`
	ProductTimeDeleted		string		`json:"productTimeDeleted"`
}

func AddProduct(product Product) (productID int64, err error) {


	var productVisible int


	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return 0, err
	}

	// TODO check if data is complete
	// is data null?
	// is email a email?
	// is username string? Numbers should be forbidden

	if product.ProductVisibility == true {
		productVisible = 1
	} else {
		productVisible = 0
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// prepare statements
	stmt, err := db.Prepare("INSERT INTO product(productName, productGTIN, productPrice, productVisible, ProductQuantity, ProductQuantityUnit) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(product.ProductName, product.ProductGTIN, product.ProductPrice, productVisible, product.ProductQuantity, product.ProductQuantityUnit)
	if err != nil {
		return 0, err

	}

	// assign id
	productID, err = res.LastInsertId()
	if err != nil {
		return 0, err

	}

	// add parent ids
	if len(product.ProductCategory) != 0 {
		for _, categoryID := range product.ProductCategory {
			stmt, err = db.Prepare("INSERT INTO productCategoryMapping(productID, categoryID) VALUES(?,?)")
			if err != nil {
				return 0, err

			}

			_ , err = stmt.Exec(product.ProductID, categoryID)
			if err != nil {
				return 0, err
			}
		}
	}

	tx.Commit()

	defer db.Close()

	return productID, nil
}

func GetProductIndex(public bool) (products []Product, err error) {

	var product Product

	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return nil, err
	}

	if public == false {
		// return all data
		rows, err := db.Query("select productID, productName, productGTIN, productPrice, productTimeAdded, productTimeChanged, productTimeDeleted, productVisible, ProductQuantity, ProductQuantityUnit from product")
		if err != nil {
			return nil, err
		}

		// write rows to struct slice
		for rows.Next() {
			err := rows.Scan(&product.ProductID, &product.ProductName, &product.ProductGTIN, product.ProductPrice, product.ProductTimeAdded, product.ProductTimeChanged, product.ProductTimeDeleted, product.ProductVisibility, product.ProductQuantity, product.ProductQuantityUnit)
			if err != nil {
				return nil, err
			}
			products = append(products, product)
		}

	} else {
		// return only public data
		// receive stuff
		rows, err := db.Query("select productID, productName, productGTIN, productPrice, ProductQuantity, ProductQuantityUnit from product WHERE productTimeDeleted IS NULL AND productVisible = 1")
		if err != nil {
			return nil, err
		}

		// write rows to struct slice
		for rows.Next() {
			err := rows.Scan(&product.ProductID, &product.ProductName, &product.ProductGTIN, product.ProductPrice, product.ProductQuantity, product.ProductQuantityUnit)
			if err != nil {
				return nil, err
			}
			products = append(products, product)
		}
	}


	// close DB connection
	defer db.Close()

	return products, nil
}