package dto

type SubmitMessageTTSReq struct {
	MessageId uint    `json:"message_id"`
	Text      string  `json:"text" binding:"required"`
	Volume    float64 `json:"volume"`
	Speed     float64 `json:"speed"`
}

type RefreshTTSTaskReq struct {
	ID int `json:"id" binding:"required"`
}

type VolcTTSSubmitRequest struct {
	AppId     string  `json:"appid"`
	ReqId     string  `json:"reqid"`
	Text      string  `json:"text"`
	Format    string  `json:"format"`
	VoiceType string  `json:"voice_type"`
	Volume    float64 `json:"volume"`
	Speed     float64 `json:"speed"`
}

type VolcTTSSubmitResponse struct {
	TaskId     string `json:"task_id"`
	TaskStatus int64  `json:"task_status"`
	TextLength int64  `json:"text_length"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

type VolcTTSTaskResult struct {
	TaskId     string `json:"task_id"`
	TaskStatus int    `json:"task_status"`
	TextLength int    `json:"task_length"`
	AudioUrl   string `json:"audio_url"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}
