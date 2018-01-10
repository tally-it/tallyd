package config

import (
	"context"
	"os"

	"github.com/marove2000/hack-and-pay/errors"

	"github.com/BurntSushi/toml"
	"github.com/marove2000/hack-and-pay/log"
)

var pkgLogger = log.New("config")

type Config struct {
	DBServer           string
	DBDatabase         string
	DBPort             string
	DBUser             string
	DBPassword         string
	JWTSecret          string
	JWTValidTime       int64
	NewUserOnlyByAdmin bool
	LDAPActive         bool
	LDAPProtocol       string
	LDAPPort           int
	LDAPHost           string
	LDAPUseSSL         bool
	LDAPBind           string
}

type LDAP struct {
	Active   bool
	Protocol string
	Port     int
	Host     string
	UseSSL   bool
	Bind     string
}

type NewConfig struct {
	Mysql *Mysql
	LDAP  *LDAP
}

type Mysql struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func ReadFile(filepath string) (*NewConfig, error) {
	logger := pkgLogger.ForFunc(context.Background(), "ReadFile")
	logger.Debug("enter config")
	conf := &Config{}
	if _, err := toml.DecodeFile(filepath, conf); err != nil {
		logger.WithError(err).Error("failed to load config file")
		return nil, errors.InternalServerError("failed to load config file", err)
	}

	return &NewConfig{
		Mysql: &Mysql{
			User:     conf.DBUser,
			Database: conf.DBDatabase,
			Host:     conf.DBServer,
			Password: conf.DBPassword,
			Port:     conf.DBPort,
		},
		LDAP: &LDAP{
			Active:   conf.LDAPActive,
			Port:     conf.LDAPPort,
			Host:     conf.LDAPHost,
			Bind:     conf.LDAPBind,
			Protocol: conf.LDAPProtocol,
			UseSSL:   conf.LDAPUseSSL,
		},
	}, nil
}

func ReadConfig() (config Config) {
	var conf Config
	if _, err := toml.DecodeFile(os.Getenv("GOPATH")+"/src/github.com/marove2000/hack-and-pay/misc/config/config.toml", &conf); err != nil {
		println(err)
	}
	return conf
}
