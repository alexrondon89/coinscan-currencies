package client

type ClientIntf interface {
	GetPrices() (interface{}, error)
}
