package gpreview

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

// GPReview : Global information of configuration.
var GPReview Config

// Config ...
type Config struct {
	MsTranslatorClientID     string `toml:"ms_tranlator_client_id"`
	MsTranslatorClientSecret string `toml:"ms_tranlator_client_secret"`

	SlackURL string `toml:"slack_url"`
}

// Load : load config
func Load() Config {
	tmlPath := "config.tml"
	if _, err := toml.DecodeFile(tmlPath, &GPReview); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Panicln("Config:Load:DecodeFile")
	}
	return GPReview
}
