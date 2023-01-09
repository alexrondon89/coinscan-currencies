package platform

import (
	"net/http"
)

var GetItemsRedisErr = NewErrorModel("error getting items from redis.", http.StatusInternalServerError)
var ItemsRedisNotFoundErr = NewErrorModel("item not found in redis.", http.StatusServiceUnavailable)
var RequestClientErr = NewErrorModel("error doing request.", http.StatusInternalServerError)
