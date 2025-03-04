package dto

type CreateMessageRequest struct {
	ConversationID string
	ExchangeID     string
	Role           string
	Content        string
	ContentType    string
}

type ClearMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	ExchangeID     string `json:"exchange_id"`
}
