package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type webHookData struct {
	Domain   string
	Path     string
	Port     string
	WithCert bool
	CertPath string
	KeyPath  string
}

type CacheData struct {
	Network string // Ex: 'tcp'
	Address string
}

type config struct {
	BotToken string
	BotName  string
	WebHook  webHookData
	Cache    CacheData
}

func initConfig() (config, error) {
	viper.SetConfigName("config")            // name of config file
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/tuenviobot/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.tuenviobot") // call multiple times to add many search paths
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		logrus.Errorf("Fatal error config file: %s \n", err)
	}
	var c config
	if err := viper.Unmarshal(&c); err != nil {
		logrus.Errorf("unable to decode config into struct, %v", err)
		return config{}, err
	}
	return c, nil
}
