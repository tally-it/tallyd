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
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func UserIndex(w http.ResponseWriter, r *http.Request) {

	//TODO: Check if Authenticated request is needed!

	// get all user data
	users := v1.GetUserIndex()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		panic(err)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {

	var	user v1.User

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// check user data
	// check if all needed data is set
	switch {
	case user.UserName == "", user.UserName == "1":
	}

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
			//TODO Log error
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(strconv.Itoa(http.StatusBadRequest) + " - wrong email format"))
			return
		}
	}

	user.UserID, err = v1.AddUser(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error " + strconv.Itoa(http.StatusBadRequest)))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.UserID); err != nil {
		panic(err)
	}
}