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

type TransferNotificationResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`

	TransferAmount    float64 `json:"transferAmount"`
	RecipientName     string  `json:"recipientName"`
	NotificationTitle string  `json:"notificationTitle"`
	NotificationBody  string  `json:"notificationBody"`
}

func JSONTransferNotification(w http.ResponseWriter, status bool, statusCode int, transferAmount float64, recipientName string) {
	res := TransferNotificationResponse{
		Status:            status,
		StatusCode:        statusCode,
		TransferAmount:    transferAmount,
		RecipientName:     recipientName,
		NotificationTitle: "Transfer Berhasil",
		NotificationBody:  "Anda telah berhasil melakukan transfer sebesar %f kepada %s.",
		Message:           "Transfer notification sent successfully",
	}

	jsonRes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonRes)
}
