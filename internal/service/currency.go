package service

import (
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type currencySrv struct {
	logger    *logrus.Logger
	config    *config.Config
	coinGecko client.Client
}

func New(logger *logrus.Logger, config *config.Config, coinGecko client.Client) internal.ServiceIntf {
	return currencySrv{
		logger:    logger,
		config:    config,
		coinGecko: coinGecko,
	}
}

func (s currencySrv) GetPricesFromApis(c *fiber.Ctx) (interface{}, error) {
	prices, err := s.coinGecko.GetPrices()
	return nil, nil
}
