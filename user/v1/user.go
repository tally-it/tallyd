package v1

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/marove2000/hack-and-pay/config/v1"
)

type User struct {
	UserID int64 `json:"userID"`
	UserName string `json:"userName"`
	UserEmail string `json:"userEmail"`
	UserPassword string `json:"UserPassword"`
	UserActive bool	`json:"userActive"`
}

func GetUserIndex() (result []User ) {

	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		log.Fatal(err)
	}

	// receive stuff
	rows, err := db.Query("select userID, userName, userTimeBlocked from user WHERE userTimeDeleted IS NULL")
	if err != nil {
		log.Fatal(err)
	}

	var (
		userID int64
		userName string
		userTimeBlocked sql.NullString
		userTimeBlockedBool bool
	)

	// write rows to struct slice
	for rows.Next() {
		err := rows.Scan(&userID, &userName, &userTimeBlocked)
		if err != nil {
			log.Fatal(err)
		}

		// check if user is blocked
		if userTimeBlocked.Valid == true {
			userTimeBlockedBool = false
		} else {
			userTimeBlockedBool = true
		}

		result = append(result, User{UserID:userID,UserName: userName,UserActive: userTimeBlockedBool})
	}

	// close DB connection
	defer db.Close()

	return result
}

func AddUser(user User) (userID int64, err error) {

	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return 0, err
	}

	// check if data is complete
	// is data null?
	// is email a email?

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	// prepare statements
	stmt, err := db.Prepare("INSERT INTO user(userName, userEmail) VALUES(?,?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(user.UserName, user.UserEmail)
	if err != nil {
		return
	}

	// assign id
	lastId, err := res.LastInsertId()
	if err != nil {
		return
	}
	user.UserID = lastId

	stmt, err = db.Prepare("INSERT INTO userAuth(userID, userAuthMethod, userAuthValue) VALUES(?,?,?)")
	if err != nil {
		return
	}

	// Create Password-Hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	res, err = stmt.Exec(user.UserID, "password", hashedPassword)
	if err != nil {
		return
	}

	tx.Commit()

	return user.UserID, nil
}