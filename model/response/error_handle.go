package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

type MyErrorHandler struct{}

func JSONErrorResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response data to JSON format", http.StatusInternalServerError)
	}
}

func (h *MyErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse := ErrorResponse{
		Message: err.Error(),
		Code:    http.StatusInternalServerError,
	}

	JSONErrorResponse(w, http.StatusInternalServerError, errorResponse)
}
