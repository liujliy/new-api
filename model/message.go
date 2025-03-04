package model

import (
	"one-api/dto"
	"time"
)

type Message struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID string    `gorm:"index;type:varchar(255);not null" json:"conversation_id"`
	ExchangeID     string    `gorm:"type:varchar(255);not null" json:"exchange_id"`
	Role           string    `gorm:"type:varchar(20);not null" json:"role"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	ContentType    string    `gorm:"type:varchar(20);not null" json:"content_type"`
	CreatedAt      time.Time `gorm:"not null;index" json:"created_at"`
}

// 根据会话ID获取消息
func GetMessagesByConversationID(conversationID string) ([]*Message, error) {
	var messages []*Message
	// 构建基础查询，添加 conversationID 条件并按时间戳排序
	err := DB.Where("conversation_id = ?", conversationID).
		Order("created_at asc").
		Find(&messages).Error
	return messages, err
}

// 创建新的消息
func (message *Message) Insert() error {
	err := DB.Create(message).Error
	return err
}

// 清理会话
func ClearMessages(param dto.ClearMessageRequest) error {
	query := DB.Where("conversation_id = ?", param.ConversationID)
	if param.ExchangeID != "" {
		query = query.Where("exchange_id = ?", param.ExchangeID)
	}
	err := query.Delete(&Message{}).Error
	return err
}
