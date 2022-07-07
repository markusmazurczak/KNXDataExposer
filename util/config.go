package util

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	APP_CONFIG      string = "app_config"
	APP_CONFIG_TYPE string = "yaml"
	ENV_PREFIX      string = "exposer"
)

/*Initializes the application configuration.

Loads the yaml configuration from app_config.yaml.
If no "path" parameter is given then the current directory of this file is used.

Every config option can be overriden by env vars.
This means if you configure an enviroment variable with the same name as an option from the config file
the value from the environment variable will take precedence.

If you want to configure an parameter via ENV you have to prefix the var name with "EXPOSER".
Lets assume that there is a variable named "port" in app_config.yaml. If you want to override the value
using an ENV variable that variable have to be named "EXPOSER_PORT"
*/
func Init_config(logger *zap.Logger, path ...string) (config interface{}, err error) {
	var p string
	logger.Sugar().Infof("%v", path)
	if len(path) == 0 || strings.TrimSpace(path[0]) == "" {
		p = "."
	} else {
		p = path[0]
	}
	return load_config(logger, p)
}

func load_config(logger *zap.Logger, path string) (config interface{}, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(APP_CONFIG)
	viper.SetConfigType(APP_CONFIG_TYPE)

	viper.SetEnvPrefix(ENV_PREFIX)
	logger.Debug("Config initialization", zap.String("path", path), zap.String("file_name", APP_CONFIG),
		zap.String("file_type", APP_CONFIG_TYPE), zap.String("env_prefix", ENV_PREFIX))

	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
