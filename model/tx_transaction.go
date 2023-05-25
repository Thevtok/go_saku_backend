package model

type Transaction struct {
	TxID                    int    `json:"tx_id"`
	TransactionType         string `json:"transaction_type"`
	TransactionDate         string `json:"transaction_date"`
	DepositBankName         string `json:"deposit_bank_name"`
	DepositBankNumber       string `json:"deposit_bank_number"`
	DepositAccountBankName  string `json:"deposit_account_bank_name"`
	DepositAmount           int    `json:"deposit_amount"`
	WithdrawBankName        string `json:"withdraw_bank_name"`
	WithdrawBankNumber      string `json:"withdraw_bank_number"`
	WithdrawAccountBankName string `json:"withdraw_account_bank_name"`
	WithdrawAmount          int    `json:"withdraw_amount"`
	TransferSenderName      string `json:"transfer_sender_name"`
	TransferSenderPhone     string `json:"transfer_sender_phone"`
	TransferRecipientName   string `json:"transfer_recipient_name"`
	TransferRecipientPhone  string `json:"transfer_recipient_phone"`
	TransferAmount          int    `json:"transfer_amount"`

	RedeemPEID   string `json:"redeem_pe_id"`
	RedeemAmount int    `json:"redeem_amount"`
	RedeemReward string `json:"redeem_reward"`
}

type Deposit struct {
	DepositID         int    `json:"deposit_id"`
	UserID            string `json:"user_id"`
	Amount            int    `json:"amount"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
	TransactionID     int    `json:"transaction_id"`
	TxID              int    `json:"tx_id"`
	TransactionType   string `json:"transaction_type"`
	TransactionDate   string `json:"transaction_date"`
}

type Withdraw struct {
	WithdrawID        int    `json:"withdraw_id"`
	UserID            string `json:"user_id"`
	Amount            int    `json:"amount"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`
	TransactionID     int    `json:"transaction_id"`
	TxID              int    `json:"tx_id"`
	TransactionType   string `json:"transaction_type"`
	TransactionDate   string `json:"transaction_date"`
}
type Transfer struct {
	TransferID           int    `json:"transfer_id"`
	SenderID             string `json:"sender_id"`
	RecipientID          string `json:"recipient_id"`
	Amount               int    `json:"amount"`
	SenderPhoneNumber    string `json:"sender_phone_number"`
	RecipientPhoneNumber string `json:"recipient_phone_number"`
	SenderName           string `json:"sender_name"`
	RecipientName        string `json:"recipient_name"`
	TransactionID        int    `json:"transaction_id"`
	TxID                 int    `json:"tx_id"`
	TransactionType      string `json:"transaction_type"`
	TransactionDate      string `json:"transaction_date"`
}

type Redeem struct {
	TransactionID int    `json:"transaction_id"`
	UserID        string `json:"user_id"`
	PEID          int    `json:"pe_id"`
	Amount        int    `json:"amount"`
}
