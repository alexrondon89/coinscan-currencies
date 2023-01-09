package handler

import (
	errCommon "github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"

	"github.com/alexrondon89/coinscan-currencies/cmd/server"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

func (cu currency) GetPrices(ctx *fiber.Ctx) error {
	prices, err := cu.currencySrv.GetPricesFromApis(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(prices)
}

func (cu currency) ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// casting to our error type
	var e errCommon.ErrorType
	if errors.As(err, &e) {
		return ctx.Status(e.StatusCode()).JSON(e)
	}

	return ctx.Status(code).SendString("Internal Server Error")
}
