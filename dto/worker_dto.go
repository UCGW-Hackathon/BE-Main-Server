package dto

type WorkerProfileUpdateRequest struct {
	FullName       *string `json:"full_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Address        *string `json:"address,omitempty"`
	Specialization *string `json:"specialization,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	BasePrice      *int    `json:"base_price,omitempty"`
	PriceUnit      *string `json:"price_unit,omitempty"`
	Services       []any   `json:"services,omitempty"`
}

type WorkerAvailabilityRequest struct {
	IsAvailable bool `json:"is_available"`
}

type WorkerOrderAcceptRequest struct {
	EstimatedArrivalMinutes *int    `json:"estimated_arrival_minutes,omitempty"`
	Note                    *string `json:"note,omitempty"`
}

type WorkerOrderRejectRequest struct {
	Reason         *string `json:"reason,omitempty"`
	ReasonCategory string  `json:"reason_category" binding:"required"`
}

type WorkerOrderStatusRequest struct {
	Status    string   `json:"status" binding:"required"`
	Note      *string  `json:"note,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type WorkerGenerateInvoiceRequest struct {
	BaseServiceFee int     `json:"base_service_fee" binding:"required"`
	WorkerNotes    *string `json:"worker_notes,omitempty"`
}

type WorkerLocationUpdateRequest struct {
	Latitude       float64  `json:"latitude" binding:"required"`
	Longitude      float64  `json:"longitude" binding:"required"`
	Heading        *float64 `json:"heading,omitempty"`
	SpeedKmh       *float64 `json:"speed_kmh,omitempty"`
	AccuracyMeters *float64 `json:"accuracy_meters,omitempty"`
	OrderID        *string  `json:"order_id,omitempty"`
}
