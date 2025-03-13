package ali

import (
	"one-api/dto"
	"strings"
)

func fileMsg2Ai(request *dto.GeneralOpenAIRequest) {
	// 处理文本对话的请求
	var newMessages []dto.Message
	for index, message := range request.Messages {
		if index == 0 && !message.IsStringContent() {
			addSysytemMessage := false
			for _, content := range message.ParseContent() {
				if content.Type == "file_info" {
					addSysytemMessage = true
					break
				}
			}
			if addSysytemMessage {
				newMessage := dto.Message{
					Role: "system",
				}
				newMessage.SetStringContent("你是一个文件对话助手")
				newMessages = append(newMessages, newMessage)
			}
		}
		if !message.IsStringContent() {
			fileInput := false
			contents := message.ParseContent()
			for _, mediaMsg := range contents {
				// 如果是file格式，则处理message
				if mediaMsg.Type == "file_info" {
					fileInput = true
					break
				}
			}
			if fileInput {
				var fileIds []string
				var userMsg string
				for _, mediaMsg := range contents {
					if mediaMsg.Type == "text" {
						userMsg = mediaMsg.Text
					}
					if mediaMsg.Type == "file_info" {
						fileIds = append(fileIds, "fileid://"+mediaMsg.FileInfo.(dto.MessageFileInfo).FileId)
					}
				}
				sysMessage := dto.Message{
					Role: "system",
				}
				sysMessage.SetStringContent(strings.Join(fileIds, ","))
				userMessage := dto.Message{
					Role: "user",
				}
				userMessage.SetStringContent(userMsg)
				newMessages = append(newMessages, sysMessage, userMessage)
			} else {
				newMessages = append(newMessages, message)
			}
		} else {
			newMessages = append(newMessages, message)
		}
	}
	request.Messages = newMessages
}
