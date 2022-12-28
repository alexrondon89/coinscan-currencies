package main

import (
	"github.com/alexrondon89/coinscan-common/logger"
	"github.com/alexrondon89/coinscan-common/redis"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/cmd/server"
	"github.com/alexrondon89/coinscan-currencies/cmd/server/handler"
	"github.com/alexrondon89/coinscan-currencies/internal/service"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client/coingecko"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client/coinmarketcap"

	"log"
)

func main() {
	configSrv, err := config.Load()
	if err != nil {
		log.Fatal("coinScan currencies service could not start due to error in configApp initialization: ", err.Error())
	}

	logger := logger.NewLogger()
	redisSrv, err := redis.New(configSrv.Redis.Host, configSrv.Redis.Port, configSrv.Redis.Password, configSrv.Redis.Db)
	if err != nil {
		log.Fatal("redis client could not be created: ", err.Error())
	}

	coingeckCli := coingecko.New(logger, configSrv)
	coinMarketCli := coinmarketcap.New(logger, configSrv)
	srv := service.New(logger, configSrv, coingeckCli, coinMarketCli, redisSrv)
	hndlr := handler.NewCurrencyHandler(logger, configSrv, srv)

	server.New(logger, configSrv, hndlr).
		AddPricesScanRoutes().
		Start()
}
