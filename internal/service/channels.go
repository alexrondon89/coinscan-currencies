package service

import (
	"context"
	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

type getCoinPrice = func(context.Context, string) ([]client.ClientResp, error.Error)

func requestPricesInCoinGecko(config *config.Config, getCoinPrice getCoinPrice) chan []internal.Info {
	ch := make(chan []internal.Info)

	for name, _ := range config.Coins {
		go func(coin string) {
			price, err := getCoinPrice(context.Background(), coin)
			if err != nil {
				ch <- []internal.Info{
					{
						Error: err.Error(),
					},
				}
				return
			}

			if len(price) > 0 {
				infoElement := []internal.Info{
					{
						Name:     price[0].Name,
						Symbol:   price[0].Symbol,
						UsdPrice: price[0].UsdPrice,
					},
				}
				ch <- infoElement
			}
		}(name)
	}

	return ch
}

func buildResponseForCoinGeckoChannel(config *config.Config, ch chan []internal.Info) []internal.Info {
	var response []internal.Info
	for i := 0; i < len(config.Coins); i++ {
		select {
		case value := <-ch:
			response = append(response, value...)
		}
	}
	close(ch)
	return response
}

func requestPricesInCoinMarketCap(config *config.Config, getCoinPrice getCoinPrice) chan []internal.Info {
	ch := make(chan []internal.Info)
	var queryParams string
	for _, symbol := range config.Coins {
		queryParams += symbol + ","
	}

	go func() {
		defer close(ch)
		price, err := getCoinPrice(context.Background(), queryParams)
		if err != nil {
			ch <- []internal.Info{
				{
					Error: err.Error(),
				},
			}
			return
		}

		var infoResponse []internal.Info
		for _, coin := range price {
			item := internal.Info{
				Name:     coin.Name,
				Symbol:   coin.Symbol,
				UsdPrice: coin.UsdPrice,
			}
			infoResponse = append(infoResponse, item)
		}

		ch <- infoResponse
	}()
	return ch
}

func buildResponseForCoinMarketCap(ch chan []internal.Info) []internal.Info {
	return <-ch
}
