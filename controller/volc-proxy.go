package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"one-api/common"
	"one-api/dto"
	"one-api/model"
	"one-api/service"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volcengine/volc-sdk-golang/base"
)

func ProxyAIGC(c *gin.Context) {
	action := c.Query("Action")
	version := c.Query("Version")

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://rtc.volcengineapi.com?Action="+action+"&Version="+version, strings.NewReader(string(bodyBytes)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	credential := base.Credentials{
		Region:          "cn-north-1",
		Service:         "rtc",
		AccessKeyID:     os.Getenv("VOLC_ACCESSKEY"),
		SecretAccessKey: os.Getenv("VOLC_SECRETKEY"),
	}
	signedReq := credential.Sign(req)

	resp, err := client.Do(signedReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"message": resp.Status,
			"success": false,
		})
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    result,
	})
}

func GetRTCToken(c *gin.Context) {
	roomID := c.Query("room_id")
	userID := c.Query("user_id")

	token := common.NewRTCToken(
		os.Getenv("VOLC_APP_ID"),
		os.Getenv("VOLC_APP_KEY"),
		roomID,
		userID,
	)
	token.ExpireTime(time.Now().Add(time.Hour * 24))
	token.AddPrivilege(common.PrivSubscribeStream, time.Time{})
	token.AddPrivilege(common.PrivPublishStream, time.Now().Add(time.Minute))

	tokenStr, err := token.Serialize()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    tokenStr,
	})
}

func GetTtsTaskById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ttsTask := model.GetTtsMessageById(id)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    ttsTask,
	})
	return
}

func SubmitTtsTask(c *gin.Context) {
	var submitReq dto.SubmitMessageTTSReq
	if err := c.ShouldBindJSON(&submitReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	taskId, err := service.SubmitTTS(submitReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    taskId,
	})
	return
}

func RefreshTtsTaskResult(c *gin.Context) {
	var refreshReq dto.RefreshTTSTaskReq
	if err := c.ShouldBindJSON(&refreshReq); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	ttsTask := model.GetTtsMessageById(refreshReq.ID)
	if ttsTask == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "任务不存在",
		})
		return
	}
	ttsTaskResult, err := service.GetTTSResult(ttsTask.TaskId)
	if err != nil {
		// 更新任务状态
		ttsTask.Status = "fail"
		ttsTask.Description = err.Error()
		ttsTask.Update()
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if ttsTaskResult.TaskStatus == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "任务正在生成中",
		})
		return
	}
	if ttsTaskResult.TaskStatus == 1 {
		ttsTask.Status = "success"
		ttsTask.Result = ttsTaskResult.AudioUrl
		ttsTask.Description = "任务生成成功"
		ttsTask.Update()
		if ttsTask.MessageId != nil {
			// 更新消息的语音结果，减少重复生成
			message := &model.Message{ID: *ttsTask.MessageId}
			message.AudioUrl = &ttsTaskResult.AudioUrl
			message.Update()
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "成功",
			"data":    ttsTaskResult.AudioUrl,
		})
		return
	}
	if ttsTaskResult.TaskStatus == 2 {
		// 更新任务状态
		ttsTask.Status = "fail"
		ttsTask.Description = "任务生成失败"
		ttsTask.Update()
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "任务生成失败",
		})
		return
	}
}

func UpdateTtsTaskBulk() {
	for {
		time.Sleep(time.Duration(3) * time.Second)
		common.SysLog("TTS任务进度轮询开始")
		tasks := model.GetProcessingTtsTask()
		if len(tasks) == 0 {
			continue
		}
		common.SysLog("检测到生成中的TTS任务")
		// 查询任务的生成结果
		for _, task := range tasks {
			taskId := task.TaskId
			ttsTask, err := service.GetTTSResult(taskId)
			if err != nil {
				task.Status = "fail"
				task.Description = err.Error()
				task.Update()
				continue
			}
			if ttsTask.TaskStatus == 0 {
				// 任务正在生成中，继续等待
				continue
			}
			if ttsTask.TaskStatus == 1 {
				// 任务生成成功，更新任务状态为成功
				task.Status = "success"
				task.Result = ttsTask.AudioUrl
				task.Description = "任务生成成功"
				task.Update()
				if task.MessageId != nil {
					// 更新消息的语音结果，减少重复生成
					message := &model.Message{ID: *task.MessageId}
					message.AudioUrl = &ttsTask.AudioUrl
					message.Update()
				}
				continue
			}
			if ttsTask.TaskStatus == 2 {
				// 任务生成失败，更新任务状态为失败
				task.Status = "fail"
				task.Description = "任务生成失败"
				task.Update()
				continue
			}
		}
	}
}
