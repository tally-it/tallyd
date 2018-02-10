package config

import (
	"context"
	"os"

	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/log"

	"github.com/BurntSushi/toml"
	"github.com/go-validator/validator"
)

var pkgLogger = log.New("config")

type Config struct {
	Mysql     *Mysql
	LDAP      *LDAP
	JWT       *JWT
	Bootstrap *Bootstrap
}

type Bootstrap struct {
	User       string `validate:"nonzero"`
	Password   string `validate:"nonzero"`
	KeepActive bool
}

type LDAP struct {
	Active             bool
	Protocol           string `validate:"nonzero"`
	Port               int    `validate:"nonzero"`
	Host               string `validate:"nonzero"`
	UseSSL             bool
	Bind               string `validate:"nonzero"`
	CAFilePath         string
	SkipInsecureVerify bool
}

type JWT struct {
	Secret    string `validate:"min=8"`
	ValidTime int64  `validate:"nonzero"`
}

type Mysql struct {
	Host     string `validate:"nonzero"`
	Port     string `validate:"nonzero"`
	Database string `validate:"nonzero"`
	User     string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

func ReadFile(filepath string) (*Config, error) {
	logger := pkgLogger.ForFunc(context.Background(), "ReadFile")
	logger.Debug("enter config")
	conf := &Config{}
	if _, err := toml.DecodeFile(filepath, conf); err != nil {
		logger.WithError(err).Error("failed to load config file")
		return nil, errors.InternalServerError("failed to load config file", err)
	}

	if conf.Bootstrap == nil {
		logger.Error("missing bootstrap config")
		return nil, errors.InternalServerError("failed to read config", nil)
	}

	err := validator.Validate(conf.Bootstrap)
	if err != nil {
		logger.WithError(err).Error("failed to validate Bootstrap config")
		return nil, errors.InternalServerError("failed to read config", err)
	}

	err = validator.Validate(conf.JWT)
	if err != nil {
		logger.WithError(err).Error("failed to validate JWT config")
		return nil, errors.InternalServerError("failed to read config", err)
	}

	err = validator.Validate(conf.Mysql)
	if err != nil {
		logger.WithError(err).Error("failed to validate mysql config")
		return nil, errors.InternalServerError("failed to read config", err)
	}

	if conf.LDAP != nil && conf.LDAP.Active {
		err = validator.Validate(conf.LDAP)
		if err != nil {
			logger.WithError(err).Error("failed to validate ldap config")
			return nil, errors.InternalServerError("failed to read config", err)
		}
	}

	return conf, nil
}

func ReadConfig() (config Config) {
	var conf Config
	if _, err := toml.DecodeFile(os.Getenv("GOPATH")+"/src/github.com/marove2000/hack-and-pay/misc/config/config.toml", &conf); err != nil {
		println(err)
	}
	return conf
}
