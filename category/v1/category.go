package v1

import (
	"database/sql"
	"github.com/marove2000/hack-and-pay/config/v1"
)

type Category struct {
	CategoryID				int64		`json:"categoryID,string"`
	CategoryName			string		`json:"categoryName"`
	CategoryVisible			bool		`json:"categoryVisible"`
	CategoryActive			bool		`json:"categoryActive"`
	CategoryParent 			[]int64		`json:"categoryParent"`
	CategoryIsRoot			bool		`json:"categoryIsRootCateogry"`
}

func AddCateogry(category Category) (categoryID int64, err error) {

	var categoryIsRoot, categoryVisible, categoryActive int


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

	if category.CategoryIsRoot == true {
		categoryIsRoot = 1
	} else {
		categoryIsRoot = 0
	}

	if category.CategoryActive == true {
		categoryActive = 1
	} else {
		categoryActive = 0
	}

	if category.CategoryVisible == true {
		categoryVisible = 1
	} else {
		categoryVisible = 0
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// prepare statements
	stmt, err := db.Prepare("INSERT INTO category(categoryName, categoryVisible, categoryActive, categoryIsRoot) VALUES(?,?,?,?)")
	if err != nil {
		return 0, err

	}

	res, err := stmt.Exec(category.CategoryName, categoryVisible, categoryActive, categoryIsRoot)
	if err != nil {
		return 0, err

	}

	// assign id
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err

	}
	categoryID = lastId

	// add parent ids
	if len(category.CategoryParent) != 0 {
		for _, categoryParent := range category.CategoryParent {
			println(categoryParent)
			stmt, err = db.Prepare("INSERT INTO categoryParentMapping(categoryID, categoryParent) VALUES(?,?)")
			if err != nil {
				return 0, err

			}

			_ , err = stmt.Exec(categoryID, categoryParent)
			if err != nil {
				return 0, err

			}
		}
	}

	tx.Commit()

	defer db.Close()

	return categoryID, nil

}