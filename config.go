package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	token    = "TOKEN"
	domain   = "DOMAIN"
	path     = "PATH"
	port     = "PORT"
	withCert = "WITH_CERT"
	certPath = "CERT_PATH"
	keyPath  = "KEY_PATH"
)

type webHookData struct {
	domain   string
	path     string
	port     string
	withCert bool
	certPath string
	keyPath  string
}

type config struct {
	botToken string
	webHook  webHookData
}

var conf config

func initConfig() {
	conf.webHook = webHookData{}
	var err error

	conf.botToken = os.Getenv(token)
	conf.webHook.withCert = false
	conf.webHook.domain = os.Getenv(domain)
	conf.webHook.path = os.Getenv(path)
	conf.webHook.certPath = os.Getenv(certPath)
	conf.webHook.keyPath = os.Getenv(keyPath)
	conf.webHook.port = os.Getenv(port)
	conf.webHook.withCert, err = strconv.ParseBool(os.Getenv(withCert))
	if err != nil {
		logrus.Warn(err)
	}
}
