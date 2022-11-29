package internal

import "github.com/gofiber/fiber/v2"

type ServiceIntf interface {
	GetPricesFromApis(c *fiber.Ctx) (interface{}, error)
}
