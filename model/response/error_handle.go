package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`

	Result any `json:"result"`
}

func JSONErrorResponse(w http.ResponseWriter, status bool, statusCode int, result any) {
	res := ErrorResponse{
		Status:     status,
		StatusCode: statusCode,
		Result:     result,

		Message: "request failed",
	}
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
