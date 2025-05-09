package model

import (
	"time"
	"github.com/google/uuid"
)

const (
	MANUAL = "manual"
	AUTO   = "auto"
)

type TimeLog struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"log_id"`
	UserID    uuid.UUID  `json:"user_id"`
	ProjectID uuid.UUID  `json:"project_id"`
	TaskName  string     `json:"task_name"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	Duration  int        `json:"duration_minutes"` 
	Source    string     `json:"source"`       
	CreatedAt time.Time  `json:"created_at"`
}

func (TimeLog) TableName() string {
	return "time_logs"
}
