package internal

import (
	"context"
	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

type ServiceIntf interface {
	GetPricesFromApis(c context.Context) (client.ClientResp, error.Error)
}
