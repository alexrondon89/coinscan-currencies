package coingecko

import (
	"context"
	"encoding/json"
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

func (cg coingecko) GetCoinPrice(c context.Context, coin string) (client.ClientResp, error.Error) {
	header := map[string][]string{"content-type": {"application/json"}}
	path := strings.Replace(cg.config.Coingecko.Url.Endpoints["coininfo"], ":coinid", coin, 1)
	resp, err := http.New("GET", cg.config.Coingecko.Url.BaseUrl, path).
		AddHeader(header).
		Exec()

	if err != nil {
		return client.ClientResp{}, error.New(platform.HttpRespErr, err)
	}

	respObject := CoinGeckoResp{}
	err = json.Unmarshal(resp.Body, &respObject)
	if err != nil {
		return client.ClientResp{}, error.New(platform.UnmarshalErr, err)
	}

	return buildClientResponse(respObject), nil
}

func buildClientResponse(respObject CoinGeckoResp) client.ClientResp {
	return client.ClientResp{
		Name:     respObject.Name,
		Symbol:   respObject.Symbol,
		UsdPrice: respObject.MarketData.CurrentPrice.Usd,
	}
}
