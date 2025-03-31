package model

import "time"

type File struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"index;not null" json:"user_id"`
	Username    string    `gorm:"type:varchar(255)" json:"username"`
	ChannelId   int       `gorm:"index;not null" json:"channel_id"`
	ChannelName string    `gorm:"type:varchar(255)" json:"channel_name"`
	FileID      string    `gorm:"type:varchar(255)" json:"file_id"`
	FileName    string    `gorm:"type:varchar(255)" json:"file_name"`
	Status      string    `gorm:"type:varchar(255)" json:"status"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}

func (file *File) Insert() error {
	var existFile File
	// 检查是否存在相同的记录
	err := DB.Where("user_id = ? AND channel_id = ? AND file_id = ?", file.UserID, file.ChannelId, file.FileID).First(&existFile).Error
	if err == nil {
		// 如果存在，则更新记录
		return file.Update()
	}
	err = DB.Create(file).Error
	return err
}

func (file *File) Update() error {
	err := DB.Save(file).Error
	return err
}
