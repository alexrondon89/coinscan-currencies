package internal

import (
	"context"
	"github.com/alexrondon89/coinscan-common/error"
)

type ServiceIntf interface {
	GetPricesFromApis(c context.Context) ([]ServiceResp, error.Error)
}

type ServiceResp struct {
	Coingecko     Coins
	CoinMarketCao Coins
	Timestamp     string `json:"timestamp,omitempty"`
}

type Coins struct {
	Bitcoin  Info
	Ethereum Info
}
type Info struct {
	Symbol   string  `json:"symbol,omitempty"`
	UsdPrice float64 `json:"usdPrice,omitempty"`
}
