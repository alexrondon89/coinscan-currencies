package client

import (
	"context"
	"github.com/alexrondon89/coinscan-common/error"
)

type ClientIntf interface {
	GetCoinPrice(c context.Context, coin string) (ClientResp, error.Error)
}

type ClientResp struct {
	Name     string  `json:"name,omitempty"`
	Symbol   string  `json:"symbol,omitempty"`
	UsdPrice float64 `json:"usdPrice,omitempty"`
}
