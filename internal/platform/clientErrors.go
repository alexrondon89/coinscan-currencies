package platform

import (
	"github.com/alexrondon89/coinscan-common/error"
	"net/http"
)

var HttpCliErr = newClientErr("error creating http instance", "0001CLIENTCURRENCIES", http.StatusInternalServerError)
var HttpRespErr = newClientErr("error doing http request", "0002CLIENTCURRENCIES", http.StatusInternalServerError)
var UnmarshalErr = newClientErr("error unmarshalling body response", "0003CLIENTCURRENCIES", http.StatusInternalServerError)

func newClientErr(message string, internalCode string, http int) error.ErrorType {
	return error.ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     http,
	}
}
