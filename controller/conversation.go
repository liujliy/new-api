package controller

import (
	"net/http"
	"one-api/model"

	"github.com/gin-gonic/gin"
)

// GetConversations 获取用户的所有会话
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

// CreateConversation 创建新的会话
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
