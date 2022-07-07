package util

import (
	"io/ioutil"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	LOG_CONFIG string = "log_config.yaml"
)

//Creates a ZAP logger from the file log_config.yaml
func Init_logger() *zap.Logger {
	yfile, err := ioutil.ReadFile(LOG_CONFIG)
	if err != nil {
		panic(err)
	}

	var cfg zap.Config
	if err := yaml.Unmarshal(yfile, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
