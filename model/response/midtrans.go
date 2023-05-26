package response

type MidtransBody struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type MidtransResponse struct {
	Status     bool         `json:"status"`
	StatusCode int          `json:"statusCode"`
	Message    string       `json:"message"`
	Result     MidtransBody `json:"result"`
}

type PaymentNotification struct {
	VANumbers []struct {
		VANumber string `json:"va_number"`
		Bank     string `json:"bank"`
	} `json:"va_numbers"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	// ... tambahkan field lainnya sesuai kebutuhan
}
