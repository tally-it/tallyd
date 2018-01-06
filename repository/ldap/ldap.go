package ldap

import (
	"crypto/tls"
	"strconv"

	"github.com/marove2000/hack-and-pay/config"
	"github.com/marove2000/hack-and-pay/errors"

	"github.com/sirupsen/logrus"
	"gopkg.in/ldap.v2"
)

const pkg = "ldap."

type LDAP struct {
	conn     *ldap.Conn
	isActive bool
}

func New(conf *config.LDAP) (*LDAP, error) {
	logger := logrus.WithField("func", pkg+"New")
	logger.Debug("enter LDAP")

	// TODO Check Certificate
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// Connect to LDAP
	ldap, err := ldap.DialTLS(conf.Protocol, conf.Host+":"+strconv.Itoa(conf.Port), tlsConfig)
	if err != nil {
		logger.WithError(err).Error("failed to connect to LDAP")
		return nil, errors.InternalServerError("failed to connect to LDAP", err)
	}

	return &LDAP{
		conn:     ldap,
		isActive: conf.Active,
	}, nil
}

func (l *LDAP) Login(name, pass string) error {
	logger := logrus.WithField("func", pkg+"LDAP.Login")
	logger.Debug("enter LDAP")

	if !l.isActive {
		logger.Error("ldap is not active")
		return errors.InternalServerError("ldap not active", nil)
	}

	err := l.conn.Bind("cn="+name+",ou=people,dc=binary-kitchen,dc=de", pass)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
			logger.WithError(err).Error("failed to connect to LDAP")
			return errors.InternalServerError("failed to connect to LDAP", err)
		}
		logger.WithError(err).Warn("failed to login to LDAP")
		return errors.Unauthorized()
	}

	return nil
}

func (l *LDAP) IsActive() bool {
	return l.isActive
}
