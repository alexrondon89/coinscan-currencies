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
	handler  CurrencyIntf
}

func New(logger *logrus.Logger, config *config.Config, handler CurrencyIntf) *WebServer {
	instance := fiber.New()
	return &WebServer{
		logger:   logger,
		config:   config,
		handler:  handler,
		instance: instance,
	}
}

func (ws *WebServer) AddPricesScanRoutes() *WebServer {
	ws.instance.Get("/prices", ws.handler.GetPrices)
	return ws
}

func (ws *WebServer) Start() {
	defer ws.instance.Shutdown()
	err := ws.instance.Listen(ws.config.Port)
	if err != nil {
		ws.logger.Fatal("coinScan-currencies service could not start: ", err.Error())
	}
}
