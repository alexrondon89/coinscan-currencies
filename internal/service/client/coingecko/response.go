package coingecko

type CoinGeckoResp struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	MarketData Prices `json:"market_data"`
	Timestamp  string `json:"timestamp"`
}

type Prices struct {
	CurrentPrice Currency `json:"current_price"`
}

type Currency struct {
	Usd float64 `json:"usd"`
}
