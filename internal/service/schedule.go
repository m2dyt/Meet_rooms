package service

import (
    "errors"

    "booking/internal/models"
    "booking/internal/repository"
    "github.com/lib/pq"
)

type ScheduleService interface {
    CreateSchedule(roomID string, daysOfWeek []int, startTime, endTime string) (*models.Schedule, error)
}

type scheduleService struct {
    scheduleRepo repository.ScheduleRepository
    roomRepo     repository.RoomRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, roomRepo repository.RoomRepository) ScheduleService {
    return &scheduleService{scheduleRepo: scheduleRepo, roomRepo: roomRepo}
}

func (s *scheduleService) CreateSchedule(roomID string, daysOfWeek []int, startTime, endTime string) (*models.Schedule, error) {
    // Проверка существования комнаты
    if _, err := s.roomRepo.FindByID(roomID); err != nil {
        return nil, errors.New("room not found")
    }
    // Проверка, нет ли уже расписания
    if existing, _ := s.scheduleRepo.FindByRoomID(roomID); existing != nil && existing.ID != "" {
        return nil, errors.New("schedule already exists for this room")
    }
    // Валидация дней недели
    for _, d := range daysOfWeek {
        if d < 1 || d > 7 {
            return nil, errors.New("daysOfWeek must be between 1 and 7")
        }
    }
    // Валидация времени
    if len(startTime) != 5 || len(endTime) != 5 || startTime[2] != ':' || endTime[2] != ':' {
        return nil, errors.New("invalid time format, use HH:MM")
    }
    // Преобразование []int в pq.Int64Array
    pqDays := make(pq.Int64Array, len(daysOfWeek))
    for i, v := range daysOfWeek {
        pqDays[i] = int64(v)
    }
    schedule := &models.Schedule{
        RoomID:     roomID,
        DaysOfWeek: pqDays,
        StartTime:  startTime,
        EndTime:    endTime,
    }
    err := s.scheduleRepo.Create(schedule)
    return schedule, err
}