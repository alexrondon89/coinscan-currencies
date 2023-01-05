package coinmarketcap

import (
	"context"
	"encoding/json"
	"errors"
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
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(platform.HttpRespErr, err))
	}

	resp, err := req.AddHeader(cmc.config.CoinMarketCap.Header).Exec()
	if err != nil {
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(platform.HttpRespErr, err))
	}

	respObject := CoinMarketCapResp{}
	err = json.Unmarshal(resp.Body, &respObject)
	if err != nil {
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(platform.HttpRespErr, err))
	}

	return buildClientResponse(coins, respObject, nil)
}

func buildClientResponse(coins string, buildClientResponse CoinMarketCapResp, err error.Error) ([]client.ClientResp, error.Error) {
	coinList := strings.Split(coins, ",")
	var cliRespElements []client.ClientResp
	for _, coin := range coinList[:len(coinList)-1] {
		found := false
		for name, value := range buildClientResponse.Data {
			if strings.ToLower(coin) == strings.ToLower(name) {
				element := client.ClientResp{
					Name:     strings.ToLower(value.Name),
					Symbol:   strings.ToLower(value.Symbol),
					UsdPrice: value.Quote.Usd.Price,
				}
				cliRespElements = append(cliRespElements, element)
				found = true
				break
			}
		}
		if !found {
			element := client.ClientResp{
				Symbol: strings.ToLower(coin),
				Error:  error.New(platform.HttpRespErr, errors.New("invalid coin")).Error(),
			}
			cliRespElements = append(cliRespElements, element)
		}
	}

	return cliRespElements, nil
}
