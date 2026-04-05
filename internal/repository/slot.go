package repository

import (
    "time"

    "booking/internal/models"
    "gorm.io/gorm"
)

type SlotRepository interface {
    Create(slot *models.Slot) error
    CreateInBatch(slots []models.Slot) error
    FindAvailableByRoomAndDate(roomID string, date time.Time) ([]models.Slot, error)
    FindByID(id string) (*models.Slot, error)
    CheckExists(roomID string, start, end time.Time) (bool, error)
}

type slotRepository struct {
    db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) SlotRepository {
    return &slotRepository{db: db}
}

func (r *slotRepository) Create(slot *models.Slot) error {
    return r.db.Create(slot).Error
}

func (r *slotRepository) CreateInBatch(slots []models.Slot) error {
    return r.db.CreateInBatches(slots, 100).Error
}

func (r *slotRepository) FindAvailableByRoomAndDate(roomID string, date time.Time) ([]models.Slot, error) {
    startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
    endOfDay := startOfDay.Add(24 * time.Hour)
    var slots []models.Slot
    err := r.db.Table("slots").
        Select("slots.*").
        Joins("LEFT JOIN bookings ON bookings.slot_id = slots.id AND bookings.status = 'active'").
        Where("slots.room_id = ? AND slots.start >= ? AND slots.start < ? AND bookings.id IS NULL", roomID, startOfDay, endOfDay).
        Order("slots.start").
        Find(&slots).Error
    return slots, err
}

func (r *slotRepository) FindByID(id string) (*models.Slot, error) {
    var slot models.Slot
    err := r.db.Where("id = ?", id).First(&slot).Error
    return &slot, err
}

func (r *slotRepository) CheckExists(roomID string, start, end time.Time) (bool, error) {
    var count int64
    err := r.db.Model(&models.Slot{}).Where("room_id = ? AND start = ? AND end = ?", roomID, start, end).Count(&count).Error
    return count > 0, err
}