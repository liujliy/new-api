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
	Title       string      `gorm:"type:varchar(255)" json:"title"`
	Model       string      `gorm:"type:varchar(50)" json:"model"`
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
	// 开始事务
	tx := DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 构建基础查询
	err := tx.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&conversations).Error

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return conversations, nil
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
