package service

import (
	"github.com/alexrondon89/coinscan-common/error"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

type currencySrv struct {
	logger    *logrus.Logger
	config    *config.Config
	coinGecko client.ClientIntf
}

func New(logger *logrus.Logger, config *config.Config, coinGecko client.ClientIntf) internal.ServiceIntf {
	return currencySrv{
		logger:    logger,
		config:    config,
		coinGecko: coinGecko,
	}
}

func (s currencySrv) GetPricesFromApis(c *fiber.Ctx) (client.ClientResp, error.Error) {
	prices, err := s.coinGecko.GetCoinPrice()
	if err != nil {
		return prices, err
	}
	return prices, nil
}
