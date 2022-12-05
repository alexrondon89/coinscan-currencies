package config

type Config struct {
	Service   string
	Version   string
	Port      string
	Coingecko service
	Cache     cache
}

type service struct {
	Url url
}

type url struct {
	BaseUrl   string
	Endpoints map[string]string
}

type cache struct {
	TimeToUpdate        uint16
	NumberOfElements    uint16
	MaxNumberOfElements uint16
}
