package coinmarketcap

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-common/error"
	httpCli "github.com/alexrondon89/coinscan-common/http/client"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal/platform"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

const (
	timeout = http.StatusRequestTimeout
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

func (cmc coinmarketcap) GetCoinPrice(ctx context.Context, coins string) ([]client.ClientResp, error.Error) {
	path := strings.Replace(cmc.config.CoinClients.CoinMarketCap.Url.Endpoints["coininfo"], ":coinid", coins, 1)
	req, err := httpCli.New("GET", cmc.config.CoinClients.CoinMarketCap.Url.BaseUrl, path, nil)
	if err != nil {
		errType := platform.HttpRespErr
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(errType.Message, errType.HttpCode, err))
	}

	ctxWithTimeOut, cancelFunc := context.WithTimeout(ctx, time.Duration(cmc.config.Http.Client.Timeout)*time.Second)
	defer cancelFunc()

	resp, err := req.
		AddHeader(cmc.config.CoinClients.CoinMarketCap.Header).
		AddContext(ctxWithTimeOut).
		Exec()
	if err != nil {
		errType := platform.HttpRespErr
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(errType.Message, errType.HttpCode, err))
	}

	if statusCode := resp.Response.StatusCode; statusCode == timeout {
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(resp.Response.Status, statusCode, nil))
	}

	respObject := CoinMarketCapResp{}
	err = json.Unmarshal(resp.Body, &respObject)
	if err != nil {
		errType := platform.HttpRespErr
		return buildClientResponse(coins, CoinMarketCapResp{}, error.New(errType.Message, errType.HttpCode, err))
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
			errType := platform.InvalidCoinErr
			element := client.ClientResp{
				Symbol: strings.ToLower(coin),
				Error:  error.New(errType.Message, errType.HttpCode, nil).Error(),
			}
			cliRespElements = append(cliRespElements, element)
		}
	}

	return cliRespElements, nil
}
