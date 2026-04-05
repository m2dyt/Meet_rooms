package service

import (
    "booking/internal/models"
    "booking/internal/repository"
)

type RoomService interface {
    CreateRoom(name string, description *string, capacity *int) (*models.Room, error)
    ListRooms() ([]models.Room, error)
}

type roomService struct {
    roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) RoomService {
    return &roomService{roomRepo: roomRepo}
}

func (s *roomService) CreateRoom(name string, description *string, capacity *int) (*models.Room, error) {
    room := &models.Room{
        Name:        name,
        Description: description,
        Capacity:    capacity,
    }
    err := s.roomRepo.Create(room)
    return room, err
}

func (s *roomService) ListRooms() ([]models.Room, error) {
    return s.roomRepo.List()
}