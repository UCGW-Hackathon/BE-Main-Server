package dto

type PurchaseApproveRequest struct {
	Note *string `json:"note,omitempty"`
}

type PurchaseRejectRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type PurchaseClarifyRequest struct {
	Question string `json:"question" binding:"required"`
}

type PurchaseBulkApproveRequest struct {
	PurchaseIDs []string `json:"purchase_ids" binding:"required"`
	Note        *string  `json:"note,omitempty"`
}

type WorkerPurchaseCreateRequest struct {
	ItemName        string  `json:"item_name" binding:"required"`
	Category        string  `json:"category" binding:"required"`
	Quantity        float64 `json:"quantity" binding:"required"`
	Unit            string  `json:"unit" binding:"required"`
	UnitPrice       int     `json:"unit_price" binding:"required"`
	TotalPrice      int     `json:"total_price" binding:"required"`
	Reason          *string `json:"reason,omitempty"`
	ReceiptPhotoURL *string `json:"receipt_photo_url,omitempty"`
}

type WorkerPurchaseAIProcessRequest struct {
	RawInput     string         `json:"raw_input" binding:"required"`
	OrderContext map[string]any `json:"order_context,omitempty"`
}

type WorkerPurchaseUpdateRequest struct {
	ItemName   string  `json:"item_name" binding:"required"`
	Category   string  `json:"category" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required"`
	Unit       string  `json:"unit" binding:"required"`
	UnitPrice  int     `json:"unit_price" binding:"required"`
	TotalPrice int     `json:"total_price" binding:"required"`
	Reason     *string `json:"reason,omitempty"`
}

type WorkerPurchaseBulkSubmitRequest struct {
	PurchaseIDs []string `json:"purchase_ids" binding:"required"`
}

type WorkerPurchaseClarifyResponseRequest struct {
	Response        string  `json:"response" binding:"required"`
	UpdatedItemName *string `json:"updated_item_name,omitempty"`
	UpdatedReason   *string `json:"updated_reason,omitempty"`
}
