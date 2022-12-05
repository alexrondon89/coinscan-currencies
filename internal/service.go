package internal

import (
	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
	"github.com/gofiber/fiber/v2"
)

type ServiceIntf interface {
	GetPricesFromApis(c *fiber.Ctx) (client.ClientResp, error.Error)
}
