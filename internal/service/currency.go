package service

import (
	"context"
	"encoding/json"
	"github.com/alexrondon89/coinscan-common/error"
	"github.com/alexrondon89/coinscan-common/redis"
	"github.com/alexrondon89/coinscan-currencies/internal/platform"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/alexrondon89/coinscan-currencies/cmd/config"
	"github.com/alexrondon89/coinscan-currencies/internal"
	"github.com/alexrondon89/coinscan-currencies/internal/service/client"
)

type currencySrv struct {
	logger    *logrus.Logger
	config    *config.Config
	coinGecko client.ClientIntf
	redis     redis.RedisIntf
}

func New(logger *logrus.Logger, config *config.Config, coinGecko client.ClientIntf, redis redis.RedisIntf) internal.ServiceIntf {
	cuSrv := currencySrv{
		logger:    logger,
		config:    config,
		coinGecko: coinGecko,
		redis:     redis,
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
				price, err := s.coinGecko.GetCoinPrice(context.Background())
				if err == nil {
					if value, err := json.Marshal(price); err == nil {
						err := s.redis.SetItem(context.Background(), price.Name, value, time.Duration(s.config.Redis.Cache.ExpirationTime)*time.Second)
						if err != nil {
							s.logger.Error("error saving with SetItem method in redis: " + err.Error())
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

func (s currencySrv) GetPricesFromApis(c context.Context) (client.ClientResp, error.Error) {
	prices, err := s.redis.GetItem(c, "BitcOin")
	if err != nil && err.Error() != "value not found" {
		return client.ClientResp{}, error.New(platform.GetItemsRedisErr, err)
	}

	if len(prices) == 0 {
		s.logger.Info("item not found in redis service... calling client coingecko")
		resp, err := s.coinGecko.GetCoinPrice(c)
		if err != nil {
			return resp, err
		}
		return resp, nil
	}

	s.logger.Info("item found from redis client")
	clientResp := client.ClientResp{}
	err = json.Unmarshal([]byte(prices), &clientResp)
	if err != nil {
		return client.ClientResp{}, error.New(platform.UnmarshalErr, err)
	}

	return clientResp, nil
}
