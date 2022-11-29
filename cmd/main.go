package main

import (
	"github.com/alexrondon89/coinscan-common/logger"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"log"
)

func main() {
	configSrv, err := config.Load()
	if err != nil {
		log.Fatal("coinScan currencies service could not start due to error in configApp initialization: ", err.Error())
	}

	logger := logger.NewLogger()

	server :=
}
