package dto

type CreateConversationRequest struct {
	Title string `json:"title"`
	Model string `json:"model"`
	// ModelConfig model.ModelConfig `json:"model_config"`
}

type DeleteConversationRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
}
