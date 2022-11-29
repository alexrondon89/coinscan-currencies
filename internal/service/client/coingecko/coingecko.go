package coingecko

import (
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
	"github.com/sirupsen/logrus"
)

type coingecko struct {
	logger *logrus.Logger
	config *config.Config
}

func New(logger *logrus.Logger, config *config.Config) client.ClientIntf {
	return coingecko{
		logger: logger,
		config: config,
	}
}

func (cg *coingecko) GetPrices() (interface{}, error) {
	return nil, nil
}
