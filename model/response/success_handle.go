package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`

	Result any `json:"result"`
}

func JSONSuccess(w http.ResponseWriter, status bool, statusCode int, result any) {
	res := SuccessResponse{
		Status:     status,
		StatusCode: statusCode,
		Result:     result,

		Message: "request success",
	}
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
