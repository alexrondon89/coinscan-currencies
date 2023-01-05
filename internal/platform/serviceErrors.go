package platform

import (
	"github.com/alexrondon89/coinscan-common/error"
	"net/http"
)

var GetItemsRedisErr = newServiceErr("error getting items from redis.", "0001SERVICECURRENCIES", http.StatusNotFound)
var ItemsRedisNotFoundErr = newServiceErr("item not found in redis.", "0002SERVICECURRENCIES", http.StatusNotFound)
var RequestClientErr = newServiceErr("error doing request. ", "0002SERVICECURRENCIES", http.StatusBadRequest)

func newServiceErr(message string, internalCode string, http int) error.ErrorType {
	return error.ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     http,
	}
}
