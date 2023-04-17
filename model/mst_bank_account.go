package model

type BankAccount struct {
	Account_ID         int64  `json:"Account_ID"`
	User_ID            int64  `json:"User_ID"`
	Bank_Name          string `json:"Bank_Name"`
	Account_Number     string `json:"Account_Number"`
	Account_HolderName string `json:"Account_HolderName"`
}
