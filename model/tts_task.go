package model

import "time"

type TtsTask struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Platform    string    `gorm:"type:varchar(50);not null" json:"platform"`
	TaskId      string    `gorm:"type:varchar(255);not null" json:"task_id"`
	Text        string    `gorm:"type:text;not null" json:"text"`
	Status      string    `gorm:"type:varchar(50);not null" json:"status"`
	Result      string    `gorm:"type:varchar(500);" json:"result"`
	Description string    `gorm:"type:varchar(500);" json:"description"`
	MessageId   *uint     `gorm:"type:int;" json:"message_id"`
	CreatedAt   time.Time `gorm:"not null;" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;" json:"updated_at"`
}

func (ttsTask *TtsTask) Insert() (int64, error) {
	err := DB.Create(ttsTask).Error
	return ttsTask.ID, err
}

func (ttsTask *TtsTask) Update() error {
	newTtsTask := *ttsTask
	DB.First(&ttsTask, ttsTask.ID)
	return DB.Model(ttsTask).Updates(newTtsTask).Error
}

func GetTtsMessageById(id int) *TtsTask {
	ttsTask := &TtsTask{}
	DB.First(ttsTask, id)
	return ttsTask
}

func GetProcessingTtsTask() []*TtsTask {
	var ttsTasks []*TtsTask
	var err error
	err = DB.Where("status = ?", "processing").Find(&ttsTasks).Error
	if err != nil {
		return nil
	}
	return ttsTasks
}
