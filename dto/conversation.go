package dto

type CreateConversationRequest struct {
	Title string `json:"title"`
	Model string `json:"model"`
	// ModelConfig model.ModelConfig `json:"model_config"`
}
