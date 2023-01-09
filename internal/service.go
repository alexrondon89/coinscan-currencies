package internal

import (
	"context"
)

type ServiceIntf interface {
	GetPricesFromApis(c context.Context) ([]ServiceResp, error)
}

type ServiceResp struct {
	CoinGecko     []Info `json:"coinGecko,omitempty"`
	CoinMarketCap []Info `json:"coinMarketCap,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
}

type Info struct {
	Name     string  `json:"name,omitempty"`
	Symbol   string  `json:"symbol,omitempty"`
	UsdPrice float64 `json:"usdPrice,omitempty"`
	Error    string  `json:"error,omitempty"`
}
