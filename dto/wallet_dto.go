package dto

type WalletWithdrawRequest struct {
	Amount            int     `json:"amount" binding:"required"`
	BankName          string  `json:"bank_name" binding:"required"`
	AccountNumber     string  `json:"account_number" binding:"required"`
	AccountHolderName string  `json:"account_holder_name" binding:"required"`
	Notes             *string `json:"notes,omitempty"`
}
