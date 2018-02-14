package ldap

import (
	"context"
	"crypto/tls"
	"strconv"

	"github.com/tally-it/tallyd/config"
	"github.com/tally-it/tallyd/errors"
	"github.com/tally-it/tallyd/log"

	"gopkg.in/ldap.v2"
	"crypto/x509"
	"io/ioutil"
)

var pkgLogger = log.New("sql")

type LDAP struct {
	conn     *ldap.Conn
	isActive bool
}

func New(conf *config.LDAP) (*LDAP, error) {
	logger := pkgLogger.ForFunc(context.Background(), "New")
	logger.Debug("enter LDAP")

	if conf == nil || !conf.Active {
		return &LDAP{}, nil
	}

	// read certificate
	cert, err := ioutil.ReadFile(conf.CAFilePath)
	if err != nil {
		logger.WithError(err).WithField("path", conf.CAFilePath).Error("failed to read LDAP CA file")
		return nil, errors.InternalServerError("failed to connect to LDAP", err)
	}

	// insert certificate to pool
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	tlsConfig := &tls.Config{InsecureSkipVerify: conf.SkipInsecureVerify, RootCAs:certPool, ServerName: conf.Host}

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

func (l *LDAP) Login(ctx context.Context, name, pass string) error {
	logger := pkgLogger.ForFunc(ctx, "Login")
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
