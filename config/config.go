package config

import (
	"osvauld/infra/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		// Ignore the error if the config file is not found
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Errorf("error to read config, %v", err)
			return err
		}
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return err
	}

	return nil
}

func GetJWTSecret() string {
	jwtSecret := viper.GetString("SECRET")
	return jwtSecret
}
