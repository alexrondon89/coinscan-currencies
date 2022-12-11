package service

type RedisObject struct {
	Timestamp Item `json:"timestamp"`
}

type Item struct {
	Symbol    string  `json:"symbol"`
	UsdPrice  float64 `json:"usdPrice"`
	Timestamp string  `json:"timestamp"`
}
