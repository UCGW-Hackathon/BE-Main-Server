package dto

type OrderLocationRequest struct {
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	Address       string  `json:"address" binding:"required"`
	AddressDetail *string `json:"address_detail,omitempty"`
}

type OrderCreateRequest struct {
	WorkerID           string               `json:"worker_id" binding:"required"`
	ServiceID          string               `json:"service_id" binding:"required"`
	Title              string               `json:"title" binding:"required"`
	Description        string               `json:"description" binding:"required"`
	Location           OrderLocationRequest `json:"location" binding:"required"`
	PreferredDate      *string              `json:"preferred_date,omitempty"`
	PreferredTimeStart *string              `json:"preferred_time_start,omitempty"`
	PreferredTimeEnd   *string              `json:"preferred_time_end,omitempty"`
	Urgency            *string              `json:"urgency,omitempty"`
	Photos             []string             `json:"photos,omitempty"`
	Notes              *string              `json:"notes,omitempty"`
}

type OrderCancelRequest struct {
	Reason         string  `json:"reason" binding:"required"`
	ReasonCategory *string `json:"reason_category,omitempty"`
}
