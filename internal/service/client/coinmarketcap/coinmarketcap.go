package coinmarketcap

import (
	"context"
	"encoding/json"
	"github.com/alexrondon89/coinscan-currencies/internal/platform"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-common/http"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

type coinmarketcap struct {
	logger *logrus.Logger
	config *config.Config
}

func New(logger *logrus.Logger, config *config.Config) client.ClientIntf {
	return coinmarketcap{
		logger: logger,
		config: config,
	}
}

func (cmc coinmarketcap) GetCoinPrice(c context.Context, coins string) ([]client.ClientResp, error.Error) {
	path := strings.Replace(cmc.config.CoinMarketCap.Url.Endpoints["coininfo"], ":coinid", coins, 1)
	req, err := http.New("GET", cmc.config.CoinMarketCap.Url.BaseUrl, path, nil)
	if err != nil {
		return nil, error.New(platform.HttpCliErr, err)
	}

	resp, err := req.AddHeader(cmc.config.CoinMarketCap.Header).Exec()
	if err != nil {
		return nil, error.New(platform.HttpRespErr, err)
	}

	respObject := CoinMarketCapResp{}
	err = json.Unmarshal(resp.Body, &respObject)
	if err != nil {
		return nil, error.New(platform.UnmarshalErr, err)
	}

	return buildClientResponse(respObject)
}

func buildClientResponse(buildClientResponse CoinMarketCapResp) ([]client.ClientResp, error.Error) {
	var cliRespElements []client.ClientResp
	for _, coin := range buildClientResponse.Data {
		element := client.ClientResp{
			Name:     strings.ToLower(coin.Name),
			Symbol:   strings.ToLower(coin.Symbol),
			UsdPrice: coin.Quote.Usd.Price,
		}

		cliRespElements = append(cliRespElements, element)
	}

	return cliRespElements, nil
}
