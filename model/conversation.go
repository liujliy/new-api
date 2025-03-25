package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID          string      `gorm:"primaryKey;type:varchar(255);not null" json:"id"`
	UserID      int         `gorm:"index;not null" json:"user_id"`
	Username    string      `gorm:"type:varchar(255)" json:"username"`
	Title       string      `gorm:"type:varchar(255)" json:"title"`
	Type        string      `gorm:"type:varchar(50)" json:"type"`
	ModelConfig ModelConfig `gorm:"type:json" json:"model_config"`
	CreatedAt   time.Time   `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"not null" json:"updated_at"`
}

type ModelConfig struct {
	Model        string  `json:"model"`
	SystemPrompt string  `json:"system_prompt"`
	Temperature  float32 `json:"temperature"`
}

func (m *ModelConfig) Scan(val interface{}) error {
	bytesValue, _ := val.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m ModelConfig) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// 根据用户ID获取会话列表
func GetConversationsByUserID(userID int) ([]*Conversation, error) {
	var conversations []*Conversation
	// 构建基础查询
	err := DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}

	return conversations, nil
}

// 获取会话列表
func ListConversations(username string, title string, startTime int64, endTime int64, startIndex int, limit int) (conversations []*Conversation, total int64, err error) {
	query := DB.Model(&Conversation{})
	if username != "" {
		query = query.Where("username = ?", username)
	}
	if title != "" {
		query = query.Where("title like %?%", title)
	}
	if startTime != 0 {
		query = query.Where("created_at >= ?", time.Unix(startTime, 0))
	}
	if endTime != 0 {
		query = query.Where("created_at <= ?", time.Unix(endTime, 0))
	}
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Order("created_at desc").Limit(limit).Offset(startIndex).Find(&conversations).Error
	if err != nil {
		return nil, 0, err
	}
	return conversations, total, nil
}

// 创建新的会话
func (conversation *Conversation) Insert() (string, error) {
	if conversation.Title == "" {
		conversation.Title = "新对话"
	}
	conversation.ID = uuid.New().String()

	err := DB.Create(conversation).Error

	return conversation.ID, err
}

// 更新会话
func (conversation *Conversation) Update() error {
	newConversation := *conversation
	DB.First(&conversation, conversation.ID)
	return DB.Model(conversation).Updates(newConversation).Error
}

// 更新会话标题
func UpdateConversationTitle(conversationID string, title string) {
	if title == "" {
		return
	}
	conversation := &Conversation{
		ID: conversationID,
	}
	err := DB.First(conversation, "id = ?", conversation.ID).Error
	if err != nil {
		return
	}
	// 已更新过无需再更新
	if conversation.Title != "新对话" {
		return
	}
	// 更新指定会话的标题
	DB.Model(conversation).Update("title", title)
}

// 删除用户的会话
func DeleteConversation(userID int, conversationID string) error {
	return DB.Where("user_id = ? and id = ?", userID, conversationID).Delete(&Conversation{}).Error
}
