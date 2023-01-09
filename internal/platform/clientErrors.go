package platform

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

var HttpCliErr = NewErrorModel("error creating http instance", fiber.StatusInternalServerError)
var HttpRespErr = NewErrorModel("error doing http request", http.StatusInternalServerError)
var UnmarshalErr = NewErrorModel("error unmarshalling body response", http.StatusInternalServerError)
var InvalidCoinErr = NewErrorModel("Invalid Coin", http.StatusBadRequest)
