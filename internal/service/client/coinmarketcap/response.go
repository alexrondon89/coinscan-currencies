package coinmarketcap

type CoinMarketCapResp struct {
	Status Status
	Data   map[string]Coin
}

type Status struct {
	Timestamp    string
	ErrorCode    int8
	ErrorMessage string
	Elapsed      int8
	CreditCount  int8
	Notice       string
}

type Coin struct {
	Name   string
	Symbol string
	Slug   string
	Quote  Quote
}

type Quote struct {
	Usd Usd
}

type Usd struct {
	Price float64
}
