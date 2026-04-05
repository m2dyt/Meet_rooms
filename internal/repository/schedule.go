package repository

import (
    "booking/internal/models"
    "gorm.io/gorm"
)

type ScheduleRepository interface {
    Create(schedule *models.Schedule) error
    FindByRoomID(roomID string) (*models.Schedule, error)
}

type scheduleRepository struct {
    db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
    return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *models.Schedule) error {
    return r.db.Create(schedule).Error
}

func (r *scheduleRepository) FindByRoomID(roomID string) (*models.Schedule, error) {
    var schedule models.Schedule
    err := r.db.Where("room_id = ?", roomID).First(&schedule).Error
    return &schedule, err
}