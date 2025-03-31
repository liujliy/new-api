package dto

type FileResponse struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Status   string `json:"status"`
}
