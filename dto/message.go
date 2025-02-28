package dto

type CreateMessageRequest struct {
	ConversationID string
	ExchangeID     string
	Role           string
	Content        string
	ContentType    string
}
