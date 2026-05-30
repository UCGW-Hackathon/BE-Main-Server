package dto

type ChatSendRequest struct {
	Content     *string `json:"content,omitempty"`
	MessageType string  `json:"message_type" binding:"required"`
}
