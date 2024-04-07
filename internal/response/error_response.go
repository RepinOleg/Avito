package response

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func HandleErrorJson(w http.ResponseWriter, errorMsg string, statusCode int) {
	errorResponse := ErrorResponse{Err: errorMsg}
	errorJSON, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(errorJSON)
	if err != nil {
		log.Println(err)
	}
}
