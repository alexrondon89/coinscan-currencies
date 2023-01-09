package platform

type ErrorModel struct {
	Message  string
	HttpCode int
}

func NewErrorModel(message string, httpCode int) ErrorModel {
	return ErrorModel{
		Message:  message,
		HttpCode: httpCode,
	}
}
