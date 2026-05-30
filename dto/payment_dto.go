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
