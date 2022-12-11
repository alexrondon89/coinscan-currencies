package service

import (
	"context"
	"encoding/json"
	"fmt"
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

func (s currencySrv) getPricesFromRedis(c context.Context, cursor uint64, match string, numberOfItems int64) ([]string, error.Error) {
	items, nextCursor, err := s.redis.Scan(c, cursor, match, numberOfItems)
	if err != nil {
		return nil, error.New(platform.GetItemsRedisErr, err)
	}
	itemsFound := len(items)
	if itemsFound == 0 {
		return items, nil
	}

	sort.Strings(items)
	s.logger.Info(fmt.Sprintf("%d elements found in redis", itemsFound))
	if itemsFound < int(numberOfItems) && nextCursor != cursor {
		difference := int(numberOfItems) - itemsFound
		s.logger.Info(fmt.Sprintf("new recursive call in redis is neccessary to find %d elements", difference))
		moreItems, err := s.getPricesFromRedis(c, nextCursor, match, int64(difference))
		if err != nil {
			return nil, error.New(platform.GetItemsRedisErr, err)
		}

		items = append(items, moreItems...)
	}

	return items, nil
}

func (s currencySrv) GetPricesFromApis(c context.Context) ([]internal.ServiceResp, error.Error) {
	items, err := s.getPricesFromRedis(c, 0, "*", int64(s.config.Redis.Cache.ItemsToRecover))
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		s.logger.Info("items not found in redis service... calling clients")
		resp, err := s.getPrices()
		if err != nil {
			return nil, err
		}
		return []internal.ServiceResp{resp}, nil
	}

	response := []internal.ServiceResp{}
	for _, key := range items {
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
