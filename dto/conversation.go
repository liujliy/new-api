package dto

type DeleteConversationRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
}
