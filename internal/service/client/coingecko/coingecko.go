package coingecko

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-common/error"
	httpCli "github.com/alexrondon89/coinscan-common/http/client"
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
	path := strings.Replace(cg.config.CoinClients.CoinGecko.Url.Endpoints["coininfo"], ":coinid", coin, 1)
	req, err := httpCli.New("GET", cg.config.CoinClients.CoinGecko.Url.BaseUrl, path, nil)
	if err != nil {
		errType := platform.HttpRespErr
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(errType.Message, errType.HttpCode, err))
	}

	resp, err := req.AddHeader(cg.config.CoinClients.CoinGecko.Header).Exec()
	if err != nil {
		errType := platform.HttpRespErr
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(errType.Message, errType.HttpCode, err))
	}

	respObject := CoinGeckoResp{}
	err = json.Unmarshal(resp.Body, &respObject)
	if err != nil {
		errType := platform.HttpRespErr
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(errType.Message, errType.HttpCode, err))
	}

	if respObject.Error != "" {
		errType := platform.HttpRespErr
		return buildClientResponse(CoinGeckoResp{Name: coin}, error.New(errType.Message, errType.HttpCode, err))
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
