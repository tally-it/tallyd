package v1

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/marove2000/hack-and-pay/config/v1"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"time"
	"crypto/tls"
	ldap2 "gopkg.in/ldap.v2"
)

type User struct {
	UserID int64 `json:"userID"`
	UserName string `json:"userName"`
	UserEmail string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
	UserActive bool	`json:"userActive"`
	UserIsAdmin bool `json:"userIsAdmin"`
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
	rows, err := db.Query("select userID, userName, userTimeBlocked, userIsAdmin from user WHERE userTimeDeleted IS NULL")
	if err != nil {
		log.Fatal(err)
	}

	var (
		userID int64
		userName string
		userTimeBlocked sql.NullString
		userTimeBlockedBool bool
		userIsAdmin int
		userIsAdminBool bool
	)

	// write rows to struct slice
	for rows.Next() {
		err := rows.Scan(&userID, &userName, &userTimeBlocked, &userIsAdmin)
		if err != nil {
			log.Fatal(err)
		}

		// check if user is blocked
		if userTimeBlocked.Valid == true {
			userTimeBlockedBool = false
		} else {
			userTimeBlockedBool = true
		}

		// check if user is admin
		if userIsAdmin == 1 {
			userIsAdminBool = true
		} else {
			userIsAdminBool = false
		}

		result = append(result, User{UserID:userID,UserName: userName,UserActive: userTimeBlockedBool, UserIsAdmin: userIsAdminBool})
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
	rows, err := db.Query("SELECT userID, userName, userTimeBlocked, userIsAdmin FROM user WHERE userID='" + strconv.FormatInt(id,10) + "'")
	if err != nil {
		return user, err
	}

	var (
		userID int64
		userName string
		userTimeBlocked sql.NullString
		userTimeBlockedBool bool
		userIsAdmin int
		userIsAdminBool bool
	)

	// write rows to struct slice
	for rows.Next() {
		err := rows.Scan(&userID, &userName, &userTimeBlocked, &userIsAdmin)
		if err != nil {
			return user, err
		}

		// check if user is blocked
		if userTimeBlocked.Valid == true {
			userTimeBlockedBool = false
		} else {
			userTimeBlockedBool = true
		}

		// check if user is admin
		if userIsAdmin == 1 {
			userIsAdminBool = true
		} else {
			userIsAdminBool = false
		}

	}

	// close DB connection
	defer db.Close()


	return User{UserID:userID,UserName: userName,UserActive: userTimeBlockedBool, UserIsAdmin: userIsAdminBool}, nil
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
	rows, err := db.Query("SELECT userID, userName, userTimeBlocked, userIsAdmin FROM user WHERE userName LIKE '" + username + "'")
	if err != nil {
		return user, err
	}

	var (
		userID int64
		userName string
		userTimeBlocked sql.NullString
		userTimeBlockedBool bool
		userIsAdmin int
		userIsAdminBool bool
	)

	// write rows to struct slice
	for rows.Next() {
		err := rows.Scan(&userID, &userName, &userTimeBlocked, &userIsAdmin)
		if err != nil {
			return user, err
		}

		// check if user is blocked
		if userTimeBlocked.Valid == true {
			userTimeBlockedBool = false
		} else {
			userTimeBlockedBool = true
		}

		// check if user is admin

		if userIsAdmin == 1 {
			userIsAdminBool = true
		} else {
			userIsAdminBool = false
		}

	}

	// close DB connection
	defer db.Close()


	return User{UserID:userID,UserName: userName,UserActive: userTimeBlockedBool, UserIsAdmin: userIsAdminBool}, nil
}

func AddUser(user User, authMethod string) (userID int64, err error) {

	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return 0, err
	}

	if conf.NewUserOnlyByAdmin == true {
		// TODO check user
	}

	// TODO
	// check if data is complete
	// is data null?
	// is email a email?
	// is username string? Numbers should be forbidden

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

	res, err = stmt.Exec(user.UserID, authMethod, hashedPassword)
	if err != nil {
		return
	}

	tx.Commit()

	return user.UserID, nil
}

func CheckPassword(user User, password []byte)(err error) {

	// get password-hash from database
	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// connect to DB
	db, err := sql.Open("mysql", conf.DBUser+":"+conf.DBPassword+"@tcp("+conf.DBServer+":"+conf.DBPort+")/"+conf.DBDatabase)
	if err != nil {
		return  err
	}

	// receive stuff
	rows, err := db.Query("SELECT userAuthValue FROM userAuth WHERE userID='" + strconv.FormatInt(user.UserID, 10)  + "' AND userAuthMethod LIKE 'password'")
	if err != nil {
		return err
	}

	var (
		userAuthValue string
	)

	// write rows to struct slice
	for rows.Next() {
		err := rows.Scan(&userAuthValue)
		if err != nil {
			return err
		}
	}

	defer db.Close()

	// compare hash
	err = bcrypt.CompareHashAndPassword([]byte(userAuthValue), password)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetJWT(user User) (tokenString string, err error) {

	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.UserID,
		"userName": user.UserName,
		"userActive": user.UserActive,
		"userIsAdmin": user.UserIsAdmin,
		"ExpiresAt": time.Now().Add(time.Second * time.Duration(conf.JWTValidTime)),
	})

	tokenString, err = token.SignedString([]byte(conf.JWTSecret))
	if err != nil {
		return "", err
	} else {
		return tokenString, err
	}
}

func GetLDAPAuthentication(user User) (err error) {

	// read config
	var conf v1.Config
	conf = v1.ReadConfig()

	// TODO Check Certificate
	tlsConfig := &tls.Config{InsecureSkipVerify:true}
	// Connect to LDAP
	ldap, err := ldap2.DialTLS(conf.LDAPProtocol, conf.LDAPHost + ":" + strconv.Itoa(conf.LDAPPort), tlsConfig)
	if err != nil {
		return err
	}
	err = ldap.Bind("cn=" + user.UserName + ",ou=people,dc=binary-kitchen,dc=de",user.UserPassword)
	if err != nil {
		//no user found
		return err
	}

	// user found and credentials are fine
	return nil
}