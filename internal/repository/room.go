package repository

import (
    "booking/internal/models"
    "gorm.io/gorm"
)

type RoomRepository interface {
    Create(room *models.Room) error
    List() ([]models.Room, error)
    FindByID(id string) (*models.Room, error)
}

type roomRepository struct {
    db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
    return &roomRepository{db: db}
}

func (r *roomRepository) Create(room *models.Room) error {
    return r.db.Create(room).Error
}

func (r *roomRepository) List() ([]models.Room, error) {
    var rooms []models.Room
    err := r.db.Find(&rooms).Error
    return rooms, err
}

func (r *roomRepository) FindByID(id string) (*models.Room, error) {
    var room models.Room
    err := r.db.Where("id = ?", id).First(&room).Error
    return &room, err
}