package config

type Config struct {
	Service   string
	Version   string
	Port      string
	Coingecko service
	Cache     cache
}

type service struct {
	BaseUrl string
}

type cache struct {
	TimeToUpdate        uint16
	NumberOfElements    uint16
	MaxNumberOfElements uint16
}
