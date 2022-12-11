package config

type Config struct {
	Service   string
	Version   string
	Port      string
	Coingecko service
	Redis     redis
}

type service struct {
	Url url
}

type url struct {
	BaseUrl   string
	Endpoints map[string]string
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
