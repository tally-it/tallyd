package sit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/marove2000/hack-and-pay/contract"

	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"io/ioutil"
)

func TestStuff(t *testing.T) {
	db := NewDB(t)
	db.Clear()

	const testuser = "testuser"
	const testpass = "testpass"

	uid := userCreate(t, testuser, testpass)
	require.NotZero(t, uid)

	token := userLogin(t, testuser, testpass)
	require.NotZero(t, token)

	user := userGetById(t, uid)
	require.Zero(t, user.Balance.Cmp(decimal.Zero))

	spew.Dump(user)
}

func userCreate(t *testing.T, name, pass string) int {
	user := &contract.AddUserRequestBody{
		Name:     name,
		Password: pass,
	}

	b, err := json.Marshal(user)
	require.NoError(t, err)

	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user", bytes.NewBuffer(b))
	require.NoError(t, err)

	c := http.Client{
		Timeout: time.Second * 5,
	}

	rsp, err := c.Do(r)
	require.NoError(t, err)
	defer rsp.Body.Close()
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var rData contract.AddUserResponseBody
	err = json.NewDecoder(rsp.Body).Decode(&rData)
	require.NoError(t, err)

	return rData.UserID
}

func userLogin(t *testing.T, name, pass string) string {
	user := &contract.LoginRequestBody{
		Name:     name,
		Password: pass,
	}

	b, err := json.Marshal(user)
	require.NoError(t, err)

	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/login", bytes.NewBuffer(b))
	require.NoError(t, err)

	c := http.Client{
		Timeout: time.Second * 5,
	}

	rsp, err := c.Do(r)
	require.NoError(t, err)
	defer rsp.Body.Close()
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	token, err := ioutil.ReadAll(rsp.Body)
	require.NoError(t, err)

	return string(token)
}

func userGetById(t *testing.T, userID int) *contract.User {
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/v1/user/"+fmt.Sprint(userID), nil)
	require.NoError(t, err)

	c := http.Client{
		Timeout: time.Second * 5,
	}

	rsp, err := c.Do(r)
	require.NoError(t, err)
	defer rsp.Body.Close()
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var rData contract.User
	err = json.NewDecoder(rsp.Body).Decode(&rData)
	require.NoError(t, err)

	return &rData
}
