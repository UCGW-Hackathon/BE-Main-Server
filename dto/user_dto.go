package dto

type UpdateUserProfileRequest struct {
	FullName  *string  `json:"full_name,omitempty"`
	Phone     *string  `json:"phone,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type UpdateUserLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Address   string  `json:"address" binding:"required"`
}
