package main

import (
	"osvauld/config"
	"osvauld/infra/database"
	"osvauld/infra/logger"
	"osvauld/routers"
	"time"

	"github.com/spf13/viper"
)

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, _ := config.DbConfiguration()

	if err := database.DbConnection(masterDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	//later separate migration

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
