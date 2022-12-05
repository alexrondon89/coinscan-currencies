package platform

import (
	"github.com/alexrondon89/coinscan-common/error"
	"net/http"
)

var HttpRespErr = newClientErr("error doing http request", "0001CLIENTCURRENCIES", http.StatusInternalServerError)
var UnmarshalErr = newClientErr("error unmarshalling body response", "0002CLIENTCURRENCIES", http.StatusInternalServerError)

func newClientErr(message string, internalCode string, http int) error.ErrorType {
	return error.ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     http,
	}
}
