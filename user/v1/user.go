package v1

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/marove2000/hack-and-pay/config/v1"
	"strconv"
)

type User struct {
	UserID int64 `json:"userID"`
	UserName string `json:"userName"`
	UserEmail string `json:"userEmail"`
	UserPassword string `json:"UserPassword"`
	UserActive bool	`json:"userActive"`
}

func GetPublicUserIndex() (result []User ) {

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

func GetPublicUserDataById(id int64) (user User, err error) {
	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return user, err
	}

	// receive stuff
	rows, err := db.Query("SELECT userID, userName, userTimeBlocked FROM user WHERE userID=" + strconv.FormatInt(id,10))
	if err != nil {
		return user, err
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
			return user, err
		}

		// check if user is blocked
		if userTimeBlocked.Valid == true {
			userTimeBlockedBool = false
		} else {
			userTimeBlockedBool = true
		}

	}

	// close DB connection
	defer db.Close()


	return User{UserID:userID,UserName: userName,UserActive: userTimeBlockedBool}, nil
}

func GetPublicUserDataByUserName(username string) (user User, err error) {
	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return user, err
	}

	// receive stuff
	rows, err := db.Query("SELECT userID, userName, userTimeBlocked FROM user WHERE userName=" + username)
	if err != nil {
		return user, err
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
			return user, err
		}

		// check if user is blocked
		if userTimeBlocked.Valid == true {
			userTimeBlockedBool = false
		} else {
			userTimeBlockedBool = true
		}

	}

	// close DB connection
	defer db.Close()


	return User{UserID:userID,UserName: userName,UserActive: userTimeBlockedBool}, nil
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

	// TODO
	// check if data is complete
	// is data null?
	// is email a email?
	// is username string? Numbers should be forbidden
	// check if username is already gone (groß und kleinschreibung beachten) muss das noch beachtet werden? datenbank sagt nein?

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