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
	err := DB.Create(file).Error
	return err
}

func (file *File) Update() error {
	err := DB.Save(file).Error
	return err
}
