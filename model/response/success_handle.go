package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message string `json:"message"`

	Result any `json:"result"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func JSONSuccess(w http.ResponseWriter, status int, data any) {
	res := SuccessResponse{
		Result: data,

		Message: "request success",
	}
	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonRes)
}
