package coingecko

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-common/http"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal/platform"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
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

func (cg coingecko) GetCoinPrice(c context.Context, coin string) ([]client.ClientResp, error.Error) {
	path := strings.Replace(cg.config.CoinGecko.Url.Endpoints["coininfo"], ":coinid", coin, 1)
	req, err := http.New("GET", cg.config.CoinGecko.Url.BaseUrl, path, nil)
	if err != nil {
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(platform.HttpRespErr, err))
	}

	resp, err := req.AddHeader(cg.config.CoinGecko.Header).Exec()
	if err != nil {
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(platform.HttpRespErr, err))
	}

	respObject := CoinGeckoResp{}
	err = json.Unmarshal(resp.Body, &respObject)
	if err != nil {
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(platform.HttpRespErr, err))
	}

	if respObject.Error != "" {
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(platform.HttpRespErr, errors.New(respObject.Error)))
	}

	return buildClientResponse(respObject, nil)
}

func buildClientResponse(respObject CoinGeckoResp, err error.Error) ([]client.ClientResp, error.Error) {

	return []client.ClientResp{
		{
			Name:     strings.ToLower(respObject.Name),
			Symbol:   strings.ToLower(respObject.Symbol),
			UsdPrice: respObject.MarketData.CurrentPrice.Usd,
			Error:    recoverErrorMsg(err),
		},
	}, nil
}

func recoverErrorMsg(err error.Error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
