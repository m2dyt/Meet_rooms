package repository

import (
    "time"

    "booking/internal/models"
    "gorm.io/gorm"
)

type BookingRepository interface {
    Create(booking *models.Booking) error
    FindBySlotIDAndActive(slotID string) (*models.Booking, error)
    FindByID(id string) (*models.Booking, error)
    UpdateStatus(id, status string) error
    ListAll(page, pageSize int) ([]models.Booking, int64, error)
    ListByUserFuture(userID string) ([]models.Booking, error)
}

type bookingRepository struct {
    db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
    return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *models.Booking) error {
    return r.db.Create(booking).Error
}

func (r *bookingRepository) FindBySlotIDAndActive(slotID string) (*models.Booking, error) {
    var booking models.Booking
    err := r.db.Where("slot_id = ? AND status = 'active'", slotID).First(&booking).Error
    return &booking, err
}

func (r *bookingRepository) FindByID(id string) (*models.Booking, error) {
    var booking models.Booking
    err := r.db.Where("id = ?", id).First(&booking).Error
    return &booking, err
}

func (r *bookingRepository) UpdateStatus(id, status string) error {
    return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookingRepository) ListAll(page, pageSize int) ([]models.Booking, int64, error) {
    var bookings []models.Booking
    var total int64
    offset := (page - 1) * pageSize
    r.db.Model(&models.Booking{}).Count(&total)
    err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&bookings).Error
    return bookings, total, err
}

func (r *bookingRepository) ListByUserFuture(userID string) ([]models.Booking, error) {
    var bookings []models.Booking
    now := time.Now().UTC()
    err := r.db.Joins("JOIN slots ON slots.id = bookings.slot_id").
        Where("bookings.user_id = ? AND bookings.status = 'active' AND slots.start >= ?", userID, now).
        Order("slots.start").
        Find(&bookings).Error
    return bookings, err
}