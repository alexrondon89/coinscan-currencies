package server

import "github.com/gofiber/fiber/v2"

type CurrencyIntf interface {
	GetPrices(c *fiber.Ctx) error
}
