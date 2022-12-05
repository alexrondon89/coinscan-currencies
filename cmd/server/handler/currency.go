package handler

import (
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/cmd/server"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type currency struct {
	logger      *logrus.Logger
	config      *config.Config
	currencySrv internal.ServiceIntf
}

func NewCurrencyHandler(logger *logrus.Logger, config *config.Config, currencySrv internal.ServiceIntf) server.CurrencyIntf {
	return currency{
		logger:      logger,
		config:      config,
		currencySrv: currencySrv,
	}
}

func (cu currency) GetPrices(c *fiber.Ctx) error {
	prices, err := cu.currencySrv.GetPricesFromApis(c)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err.Type())
	}

	return c.Status(200).JSON(prices)
}
