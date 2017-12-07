package v1

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	DBServer   string
	DBDatabase string
	DBPort     string
	DBUser     string
	DBPassword string
	JWTSecret   string
	JWTValidTime int64
	NewUserOnlyByAdmin bool
	LDAPActive bool
	LDAPProtocol string
	LDAPPort int
	LDAPHost string
	LDAPUseSSL bool
	LDAPBind string
}

func ReadConfig() (config Config) {
	var conf Config
	if _, err := toml.DecodeFile(os.Getenv("GOPATH") + "/src/github.com/marove2000/hack-and-pay/misc/config/config.toml", &conf); err != nil {
		println(err)
	}
	return conf
}
