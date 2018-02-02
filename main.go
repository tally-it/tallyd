package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/marove2000/hack-and-pay/config"
	"github.com/marove2000/hack-and-pay/handler"
	"github.com/marove2000/hack-and-pay/log"
	"github.com/marove2000/hack-and-pay/repository/ldap"
	"github.com/marove2000/hack-and-pay/repository/sql"
	"github.com/marove2000/hack-and-pay/router"

	"github.com/sirupsen/logrus"
)

var (
	pkgLogger  = log.New("main")
	isDebug    = flag.Bool("debug", false, "enables debug mode")
	configPath = flag.String("config", os.Getenv("GOPATH")+"/src/github.com/marove2000/hack-and-pay/misc/config/config.toml", "path to config")
)

func main() {
	logger := pkgLogger.ForFunc(context.Background(), "main")
	logger.Info("startup")

	flag.Parse()
	if *isDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	conf, err := config.ReadFile(*configPath)
	if err != nil {
		logger.Fatal("bailing")
	}

	db, err := sql.New(conf.Mysql)
	if err != nil {
		logger.Fatal("bailing")
	}

	l, err := ldap.New(conf.LDAP)
	if err != nil {
		logger.Fatal("bailing")
	}

	r := router.NewRouter(handler.New(db, l))

	logger.Info("running...")

	logger.WithError(http.ListenAndServe(":8080", r)).Error("bailing")
}
