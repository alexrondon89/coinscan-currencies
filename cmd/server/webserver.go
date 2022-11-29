package server

import (
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WebServer struct {
	logger   *logrus.Logger
	config   *config.Config
	instance *fiber.App
	handler  handler.Handler
}

func New(logger *logrus.Logger, config *config.Config, handler handler.Handler) WebServer {
	return WebServer{
		logger:  logger,
		handler: handler,
	}
}

func (ws *WebServer) AddPricesScanRoutes() *WebServer {
	ws.instance.Get("/prices", ws.handler.GetCoinPrices)
	return ws
}

func (ws *WebServer) Start(config *config.Config) {
	defer ws.instance.Shutdown()
	err := ws.instance.Listen(config.Port)
	if err != nil {
		ws.logger.Fatal("coinScan-currencies service could not start: ", err.Error())
	}
}
