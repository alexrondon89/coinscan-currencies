package config

type Config struct {
	Service     string
	Version     string
	Port        string
	CoinClients coinClients
	Redis       redis
	Coins       map[string]string
	Http        http
}

type cache struct {
	TimeToUpdate   uint16
	ExpirationTime uint16
	ItemsToRecover uint16
}

type redis struct {
	Host     string
	Port     string
	Password string
	Db       int
	Cache    cache
}

type http struct {
	Server server
	Client client
}

type server struct {
	ReadTimeout  uint8
	WriteTimeout uint8
	IdleTimeout  uint8
}

type client struct {
	Timeout uint8
}

type coinClients struct {
	CoinGecko     service
	CoinMarketCap service
}

type service struct {
	Url    url
	Header map[string][]string
}

type url struct {
	BaseUrl   string
	Endpoints map[string]string
}
