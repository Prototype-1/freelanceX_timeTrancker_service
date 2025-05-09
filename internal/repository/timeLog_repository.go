package repository

import (
"gorm.io/gorm"
	"github.com/google/uuid"
	"github.com/Prototype-1/freelanceX_timeTrancker_service/internal/model"
	"time"
)

type TimeLogRepository interface {
	CreateTimeLog(log *model.TimeLog) (*model.TimeLog, error)
	GetTimeLogsByUser(userID uuid.UUID, projectID uuid.UUID, dateFrom, dateTo *time.Time) ([]model.TimeLog, error)
	GetTimeLogsByProject(projectID uuid.UUID, dateFrom, dateTo *time.Time) ([]model.TimeLog, error)
	GetTimeLogByID(logID uuid.UUID) (*model.TimeLog, error)
	UpdateTimeLog(log *model.TimeLog) (*model.TimeLog, error)
	DeleteTimeLog(logID uuid.UUID) error
}

type timeLogRepository struct {
	db *gorm.DB
}

func NewTimeLogRepository(db *gorm.DB) TimeLogRepository {
	return &timeLogRepository{db: db}
}

func (r *timeLogRepository) CreateTimeLog(log *model.TimeLog) (*model.TimeLog, error) {
	if err := r.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (r *timeLogRepository) GetTimeLogsByUser(userID uuid.UUID, projectID uuid.UUID, dateFrom, dateTo *time.Time) ([]model.TimeLog, error) {
	var logs []model.TimeLog
	query := r.db.Where("user_id = ?", userID)

	if projectID != uuid.Nil {
		query = query.Where("project_id = ?", projectID)
	}

	if dateFrom != nil {
		query = query.Where("start_time >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where("end_time <= ?", *dateTo)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *timeLogRepository) GetTimeLogsByProject(projectID uuid.UUID, dateFrom, dateTo *time.Time) ([]model.TimeLog, error) {
	var logs []model.TimeLog
	query := r.db.Where("project_id = ?", projectID)

	if dateFrom != nil {
		query = query.Where("start_time >= ?", *dateFrom)
	}

	if dateTo != nil {
		query = query.Where("end_time <= ?", *dateTo)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *timeLogRepository) GetTimeLogByID(logID uuid.UUID) (*model.TimeLog, error) {
	var log model.TimeLog
	if err := r.db.First(&log, "id = ?", logID).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *timeLogRepository) UpdateTimeLog(log *model.TimeLog) (*model.TimeLog, error) {
	if err := r.db.Save(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (r *timeLogRepository) DeleteTimeLog(logID uuid.UUID) error {
	if err := r.db.Delete(&model.TimeLog{}, logID).Error; err != nil {
		return err
	}
	return nil
}
