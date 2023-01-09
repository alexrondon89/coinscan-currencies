package server

import (
	"github.com/gofiber/fiber/v2"
)

type CurrencyIntf interface {
	GetPrices(c *fiber.Ctx) error
	ErrorHandler(ctx *fiber.Ctx, err error) error
}
