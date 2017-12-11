package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/badoux/checkmail"
	"strconv"
	"github.com/marove2000/hack-and-pay/user/v1"
	config "github.com/marove2000/hack-and-pay/config/v1"
	payment "github.com/marove2000/hack-and-pay/payment/v1"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func PublicUserIndex(w http.ResponseWriter, r *http.Request) {

	//TODO: Check if Authenticated request is needed!

	// get all user data
	users := v1.GetPublicUserIndex()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(strconv.Itoa(http.StatusInternalServerError)))
		return
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {

	conf := config.ReadConfig()

	if conf.LDAPActive == false {
		var	user v1.User

		// get body data
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			log.Println(err)
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(strconv.Itoa(http.StatusInternalServerError)))
			return
		}
		defer r.Body.Close()

		// check user data
		// check if all needed data is set
		if user.UserName == "" {
			log.Println("No user name")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(strconv.Itoa(http.StatusBadRequest) + " - no user name set"))
			return
		} else if user.UserPassword == "" {
			log.Println("No password set")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(strconv.Itoa(http.StatusBadRequest) + " - no password set"))
			return
		} else if user.UserEmail != "" {
			// check if email
			err = checkmail.ValidateFormat(user.UserEmail)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(strconv.Itoa(http.StatusBadRequest) + " - wrong email format"))
				return
			}
		}

		user.UserID, err = v1.AddUser(user, "password")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error " + strconv.Itoa(http.StatusInternalServerError)))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user.UserID); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(strconv.Itoa(http.StatusInternalServerError)))
			return
		}
	} else {
		log.Println("LDAP active, creation of local user not allowed")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(strconv.Itoa(http.StatusBadRequest) + " - Creation of local user not allowed"))
		return
	}

}

func GetPublicUserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var user v1.User

	// Careful vars["id"] can be email, username or id
	id := vars["id"]

	if _, err := strconv.Atoi(id); err == nil {

		id64, err:= strconv.ParseInt(vars["id"],10,64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error " + strconv.Itoa(http.StatusInternalServerError)))
			return
		}
		user, err = v1.GetPublicUserDataById(id64)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error " + strconv.Itoa(http.StatusInternalServerError)))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error " + strconv.Itoa(http.StatusInternalServerError)))
			return
		}

	} else {

		user, err = v1.GetPublicUserDataByUserName(id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error " + strconv.Itoa(http.StatusInternalServerError)))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error " + strconv.Itoa(http.StatusInternalServerError)))
			return
		}
	}

	return
}

func GetAuthentication(w http.ResponseWriter, r *http.Request) {

	var user, dbUser v1.User

	conf := config.ReadConfig()

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	// check if user is in Database and get public user data
	dbUser, err = v1.GetPublicUserDataByUserName(user.UserName)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if dbUser.UserID == 0 {
		// user does not exist, check ldap first
		err = v1.GetLDAPAuthentication(user)
		if err != nil {
			// user does not exist in ldap too or credentials are wrong
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			// user is existing in ldap
			// TODO get LDAP Email Address from LDAP-field mail
			v1.AddUser(user, "ldap")
			user.UserJWT, err = v1.JWTGet(user)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(user); err != nil {
					log.Println(err)
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

		// if ldap also does not exist
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if conf.LDAPActive == false {
		// LDAP not active, check password in database

		err = v1.CheckPassword(dbUser, []byte(user.UserPassword))
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			dbUser.UserJWT, err = v1.JWTGet(dbUser)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		// LDAP active, check with LDAP
		err = v1.GetLDAPAuthentication(user)
		if err != nil {
			// user does not exist in ldap too or credentials are wrong
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			dbUser.UserJWT, err = v1.JWTGet(dbUser)
			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbUser); err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func ChangeBalance (w http.ResponseWriter, r *http.Request) {

	var user, dbUser v1.User

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(strconv.Itoa(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	if user.UserID == 0 {
		dbUser, err = v1.GetPublicUserDataByUserName(user.UserName)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if user.UserName == "" {
		dbUser, err = v1.GetPublicUserDataById(user.UserID)
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else if user.UserID == 0 && user.UserName == "" {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check JWT
	err = v1.JWTValidate(user.UserJWT, dbUser.UserID, false)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// get get data
	vars := mux.Vars(r)
	Change, err := strconv.ParseFloat(vars["change"], 64)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// only do something if change is not 0
	if Change != 0 {
		println("oiawjedfoi")
		err = payment.PaymentTransfer(dbUser.UserID, Change, "")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(strconv.Itoa(http.StatusInternalServerError)))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		return

	} else {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}