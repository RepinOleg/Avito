package response

import "net/http"

type ErrorResponse struct {
	Err string `json:"error"`
}

func (e ErrorResponse) Error() string {
	return e.Err
}

// HandleError Обработчик ошибок
func HandleError(w http.ResponseWriter, errorMsg string, statusCode int) {
	http.Error(w, errorMsg, statusCode)
}
