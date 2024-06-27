package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port                 string
	Secret               string
	LimitCountPerRequest int64
}

func ServerConfig() string {
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")

	appServer := fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT"))
	log.Print("Server Running at :", appServer)
	return appServer
}
