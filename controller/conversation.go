package controller

import (
	"net/http"
	"one-api/common"
	"one-api/dto"
	"one-api/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取用户的所有会话
func GetConversations(c *gin.Context) {
	userId := c.GetInt("id")
	conversations, err := model.GetConversationsByUserID(userId)
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
		"data":    conversations,
	})
}

func ListConversations(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = common.ItemsPerPage
	}
	username := c.Query("username")
	title := c.Query("title")
	conversationType := c.Query("type")
	startTime, _ := strconv.ParseInt(c.Query("start_time"), 10, 64)
	endTime, _ := strconv.ParseInt(c.Query("end_time"), 10, 64)
	conversations, total, err := model.ListConversations(username, title, conversationType, startTime, endTime, (page-1)*page, pageSize)
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
		"data": map[string]any{
			"items":     conversations,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// 创建新的会话
func CreateConversation(c *gin.Context) {
	var conversation model.Conversation
	if err := c.ShouldBindJSON(&conversation); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	conversation.UserID = c.GetInt("id")
	conversation.Username = c.GetString("username")
	conversationId, err := conversation.Insert()
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
		"data":    conversationId,
	})
}

func UpdateConversation(c *gin.Context) {
	var conversation model.Conversation
	if err := c.ShouldBindJSON(&conversation); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	err := conversation.Update()
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
	})
}

// 删除会话
func DeleteConversation(c *gin.Context) {
	var deleteConversationRequest dto.DeleteConversationRequest
	if err := c.ShouldBindJSON(&deleteConversationRequest); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	err := model.DeleteConversation(c.GetInt("id"), deleteConversationRequest.ConversationID)
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
	})
}
