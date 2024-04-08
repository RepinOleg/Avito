package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Error struct {
	Message string `json:"error"`
}

func (e *Error) Error() string {
	return e.Message
}

type AccessError struct {
	Message string `json:"error"`
}

func (e *AccessError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string `json:"error"`
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func HandleError(w http.ResponseWriter, err error) {
	var accessErr *AccessError
	if errors.As(err, &accessErr) {
		http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
		return
	}

	var notFoundErr *NotFoundError
	if errors.As(err, &notFoundErr) {
		http.Error(w, "Баннер не найден", http.StatusNotFound)
		return
	}

	HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
}

func HandleErrorJson(w http.ResponseWriter, errorMsg string, statusCode int) {
	errorResponse := Error{Message: errorMsg}

	errorJSON, _ := json.Marshal(errorResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err := w.Write(errorJSON)
	if err != nil {
		log.Println(err)
	}
}
