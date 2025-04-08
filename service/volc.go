package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"one-api/common"
	"one-api/dto"
	"one-api/model"
	"os"
)

func GetTTSResult(taskId string) (*dto.VolcTTSTaskResult, error) {
	// Create a new request with the desired URL and parameters
	req, err := http.NewRequest("GET", "https://openspeech.bytedance.com/api/v1/tts_async/query?appid="+os.Getenv("VLOC_TTS_APP_ID")+"&task_id="+taskId, nil)
	if err != nil {
		return nil, errors.New("查询语音合成请求初始化失败")
	}
	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer;"+os.Getenv("VLOC_TTS_ACCESS_TOKEN"))
	req.Header.Set("Resource-Id", "volc.tts_async.default")

	client := GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("查询语音合成请求失败")
	}
	responseBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("语音合成请求失败，状态码：" + resp.Status)
	}

	var response dto.VolcTTSTaskResult

	err = common.DecodeJson(responseBody, &response)
	if err != nil {
		log.Println("解析语音合成结果失败", err.Error())
		return nil, errors.New("解析语音合成结果失败")
	}
	if response.Code != 0 {
		return nil, errors.New("语音合成失败，错误：" + response.Message)
	}
	return &response, nil
}

func SubmitTTS(submitReq dto.SubmitMessageTTSReq) (int64, error) {
	// Set the request body
	submitBody := dto.VolcTTSSubmitRequest{
		AppId:     os.Getenv("VLOC_TTS_APP_ID"),
		ReqId:     common.GetUUID(),
		Text:      submitReq.Text,
		Format:    "mp3",
		VoiceType: "BV001_streaming",
	}
	jsonData, err := json.Marshal(submitBody)
	if err != nil {
		return 0, errors.New("语音合成请求参数错误")
	}
	requestBody := bytes.NewBuffer(jsonData)
	// Create a new request with the desired URL and parameters
	req, err := http.NewRequest("POST", "https://openspeech.bytedance.com/api/v1/tts_async/submit", requestBody)
	if err != nil {
		return 0, errors.New("语音合成请求初始化失败")
	}
	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer;"+os.Getenv("VLOC_TTS_ACCESS_TOKEN"))
	req.Header.Set("Resource-Id", "volc.tts_async.default")

	client := GetHttpClient()
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("语音合成请求失败")
	}
	responseBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("语音合成请求失败，状态码：" + resp.Status)
	}

	var response dto.VolcTTSSubmitResponse
	err = common.DecodeJson(responseBody, &response)
	// save tts task
	ttsTask := &model.TtsTask{
		Platform: "volc",
		TaskId:   response.TaskId,
		Text:     submitReq.Text,
		Status:   "processing",
	}
	if submitReq.MessageId != 0 {
		ttsTask.MessageId = &submitReq.MessageId
	}
	return ttsTask.Insert()
}
