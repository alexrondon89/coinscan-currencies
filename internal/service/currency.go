package service

import (
	"context"
	"encoding/json"
	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-common/redis"
	"github.com/alexrondon89/coinscan-currencies/internal/platform"
	"github.com/sirupsen/logrus"
	"sort"
	"time"

	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

type currencySrv struct {
	logger        *logrus.Logger
	config        *config.Config
	coingecko     client.ClientIntf
	coinmarketcap client.ClientIntf
	redis         redis.RedisIntf
}

func New(logger *logrus.Logger, config *config.Config, coingecko, coinmarketcap client.ClientIntf, redis redis.RedisIntf) internal.ServiceIntf {
	cuSrv := currencySrv{
		logger:        logger,
		config:        config,
		coingecko:     coingecko,
		coinmarketcap: coinmarketcap,
		redis:         redis,
	}
	cuSrv.updateCacheLastPrices()
	return cuSrv
}

func (s currencySrv) updateCacheLastPrices() {
	go func() {
		ticker := time.NewTicker(time.Duration(s.config.Redis.Cache.TimeToUpdate) * time.Second)
		for {
			select {
			case <-ticker.C:
				serviceResp, err := s.getPrices()
				if err == nil {
					if value, err := json.Marshal(serviceResp); err == nil {
						expiration := time.Duration(s.config.Redis.Cache.ExpirationTime) * time.Second
						err := s.redis.SetItem(context.Background(), serviceResp.Timestamp, value, expiration)
						if err != nil {
							s.logger.Error("error saving with SetItem method in redis: " + err.Error())
							continue
						}
						s.logger.Info("prices updated in redis cache... ")
						continue
					}
					s.logger.Error("problem marshalling info for redis: " + err.Error())
				}
				s.logger.Error("problem getting coin price: " + err.Error())
			}
		}
	}()
}

func (s currencySrv) GetPricesFromApis(c context.Context) ([]internal.ServiceResp, error.Error) {
	prices, _, err := s.redis.Scan(c, 0, "*", int64(s.config.Redis.Cache.ItemsToRecover)-1)
	if err != nil && err.Error() != "value not found" {
		return []internal.ServiceResp{}, error.New(platform.GetItemsRedisErr, err)
	}

	if len(prices) == 0 {
		s.logger.Info("item not found in redis service... calling client coingecko")
		resp, err := s.getPrices()
		if err != nil {
			return nil, err
		}
		return []internal.ServiceResp{resp}, nil
	}

	sort.Strings(prices)
	response := []internal.ServiceResp{}
	for _, key := range prices {
		item, err := s.redis.GetItem(c, key)
		if err != nil {
			return nil, error.New(platform.GetItemsRedisErr, err)
		}

		clientResp := internal.ServiceResp{}
		err = json.Unmarshal([]byte(item), &clientResp)
		if err != nil {
			return nil, error.New(platform.UnmarshalErr, err)
		}
		response = append(response, clientResp)
	}

	return response, nil
}

func (s currencySrv) getPrices() (internal.ServiceResp, error.Error) {
	key := time.Now().UTC().String()
	priceBtc, err := s.coingecko.GetCoinPrice(context.Background(), "bitcoin")
	priceEth, err := s.coingecko.GetCoinPrice(context.Background(), "ethereum")
	if err != nil {
		return internal.ServiceResp{}, error.New(platform.RequestClientErr, err)
	}

	redisObj := internal.ServiceResp{
		Coingecko: internal.Coins{
			Bitcoin: internal.Info{
				Symbol:   priceBtc.Symbol,
				UsdPrice: priceBtc.UsdPrice,
			},
			Ethereum: internal.Info{
				Symbol:   priceEth.Symbol,
				UsdPrice: priceEth.UsdPrice,
			},
		},
		Timestamp: key,
	}

	return redisObj, nil
}
