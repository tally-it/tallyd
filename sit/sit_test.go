package sit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/tally-it/tallyd/contract"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	db := NewDB(t)
	db.Clear()

	const testuser = "testuser"
	const testpass = "testpass"

	uid := userCreate(t, testuser, testpass)
	require.NotZero(t, uid)

	token := userLogin(t, testuser, testpass)
	require.NotZero(t, token)

	user := userGetById(t, token, uid)
	require.Zero(t, user.Balance.Cmp(decimal.Zero))
	require.False(t, bool(user.IsAdmin))

	transactionAdd(t, token, uid, 15.33)

	user = userGetById(t, token, uid)
	require.Zero(t, user.Balance.Cmp(decimal.NewFromFloat(15.33)))

	user = userGetById(t, "", uid)
	require.Zero(t, user.Balance.Cmp(decimal.Zero))
}

func userCreate(t *testing.T, name, pass string) int {
	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user", bytes.NewBufferString(`
{	
	"name":"`+name+`",
	"password":"`+pass+`"
}
`))
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
	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/login", bytes.NewBufferString(`
{	
	"name":"`+name+`",
	"password":"`+pass+`"
}
`))
	require.NoError(t, err)

	c := http.Client{
		Timeout: time.Second * 5,
	}

	rsp, err := c.Do(r)
	require.NoError(t, err)
	defer rsp.Body.Close()
	require.Equal(t, http.StatusOK, rsp.StatusCode)

	var rspData contract.LoginResponse
	err = json.NewDecoder(rsp.Body).Decode(&rspData)
	require.NoError(t, err)

	return rspData.JWT
}

func userGetById(t *testing.T, token string, userID int) *contract.User {
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/v1/user/"+fmt.Sprint(userID), nil)
	require.NoError(t, err)

	if token != "" {
		r.Header.Set("Authorization", "Bearer: "+token)
	}

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

func transactionAdd(t *testing.T, token string, userId int, value float64) {
	r, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user/"+fmt.Sprint(userId)+"/transaction", bytes.NewBufferString(`
{	
	"userID":`+fmt.Sprint(userId)+`,
	"value":`+fmt.Sprint(value)+`
}
`))
	r.Header.Set("Authorization", "Bearer: "+token)

	require.NoError(t, err)

	c := http.Client{
		Timeout: time.Second * 5,
	}
	rsp, err := c.Do(r)
	require.NoError(t, err)
	defer rsp.Body.Close()
	require.Equal(t, http.StatusNoContent, rsp.StatusCode)
}
