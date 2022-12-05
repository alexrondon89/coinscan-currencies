package main

import (
	"github.com/alexrondon89/coinscan-common/logger"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/cmd/server"
	"github.com/alexrondon89/coinscan-currencies/cmd/server/handler"
	"github.com/alexrondon89/coinscan-currencies/internal/service"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client/coingecko"
	"log"
)

func main() {
	configSrv, err := config.Load()
	if err != nil {
		log.Fatal("coinScan currencies service could not start due to error in configApp initialization: ", err.Error())
	}

	logger := logger.NewLogger()

	coingeckCli := coingecko.New(logger, configSrv)
	srv := service.New(logger, configSrv, coingeckCli)
	hndlr := handler.NewCurrencyHandler(logger, configSrv, srv)
	server.New(logger, configSrv, hndlr).
		AddPricesScanRoutes().
		Start()
}
