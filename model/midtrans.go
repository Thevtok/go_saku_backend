package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/ReygaFitra/inc-final-project.git/model/response"
	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func generateSerialNumber() int {
	// Generate random 4-digit number
	rand.Seed(time.Now().UnixNano())
	serialNumber := rand.Intn(10000)

	return serialNumber
}
func CreateMidtransTransactionFromDeposit(depo *Deposit, user *User) (string, error) {
	serialNumber := generateSerialNumber()
	orderID := fmt.Sprintf("DEPOSIT-%04d", serialNumber)

	customerName := depo.AccountHolderName
	totalAmount := int64(depo.Amount)

	midtrans.ServerKey = utils.DotEnv("SERVER_KEY")
	midtrans.ClientKey = utils.DotEnv("CLIENT_KEY")

	payment := midtrans.TransactionDetails{
		OrderID:  orderID,
		GrossAmt: totalAmount,
	}

	customer := &midtrans.CustomerDetails{
		FName: customerName,
		Email: user.Email,
		Phone: user.Phone_Number,
		BillAddr: &midtrans.CustomerAddress{
			Address: user.Address,
		},
	}

	request := &snap.Request{
		TransactionDetails: payment,
		CustomerDetail:     customer,
		EnabledPayments:    snap.AllSnapPaymentType,
	}

	// Buat client Snap Midtrans
	snapAPIClient := snap.Client{
		ServerKey:  midtrans.ServerKey,
		Env:        midtrans.Sandbox,
		HttpClient: midtrans.GetHttpClient(midtrans.Environment),
	}

	// Buat transaksi dengan Snap Midtrans
	resp, err := snapAPIClient.CreateTransaction(request)
	if err != nil {
		return "", fmt.Errorf("failed to create Midtrans transaction: %v", err)
	}

	if resp.Token != "" {
		response := &response.MidtransResponse{
			Status:     true,
			StatusCode: 201,
			Message:    "request success",
			Result: response.MidtransBody{
				Token:       resp.Token,
				RedirectURL: resp.RedirectURL,
			},
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return "", fmt.Errorf("failed to marshal Midtrans response: %v", err)
		}

		return string(responseJSON), nil
	}

	return "", fmt.Errorf("no payment token found")
}
