package dto

type PaymentCreateRequest struct {
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	PaymentProofURL *string `json:"payment_proof_url,omitempty"`
	PromoCode       *string `json:"promo_code,omitempty"`
}

type SandboxPaymentCallbackRequest struct {
	PaymentID string `json:"payment_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

type MidtransWebhookRequest struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type PaymentResponse struct {
	PaymentID       string  `json:"payment_id"`
	OrderID         string  `json:"order_id"`
	InvoiceID       string  `json:"invoice_id"`
	Amount          int     `json:"amount"`
	PaymentStatus   string  `json:"payment_status"`
	Token           *string `json:"token,omitempty"`
	RedirectURL     *string `json:"redirect_url,omitempty"`
	CreatedAt       string  `json:"created_at"`
}
