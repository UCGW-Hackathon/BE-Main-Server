package dto

type RatingCreateRequest struct {
	Rating  int      `json:"rating" binding:"required"`
	Comment *string  `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

type CustomerRatingCreateRequest struct {
	Rating  int      `json:"rating" binding:"required"`
	Comment *string  `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}
