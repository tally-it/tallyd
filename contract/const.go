package contract

type AuthType string

const (
	AuthTypeLDAP   AuthType = "ldap"
	AuthTypePasswd          = "passwd"
)
