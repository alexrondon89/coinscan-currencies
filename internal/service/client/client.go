package client

import (
	"context"
	"github.com/alexrondon89/coinscan-common/error"
)

type ClientIntf interface {
	GetCoinPrice(c context.Context) (ClientResp, error.Error)
}

type ClientResp struct {
	Name string `json:"name,omitempty"`
	Info Coin   `json:"Info,omitempty"`
}

type Coin struct {
	Symbol    string  `json:"symbol,omitempty"`
	UsdPrice  float64 `json:"usdPrice,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
}
