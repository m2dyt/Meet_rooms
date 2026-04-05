package service

import (
    "errors"
    "log"
    "time"

    "booking/internal/models"
    "booking/internal/repository"
)

type SlotService interface {
    GetAvailableSlots(roomID string, dateStr string) ([]models.Slot, error)
}

type slotService struct {
    slotRepo     repository.SlotRepository
    scheduleRepo repository.ScheduleRepository
}

func NewSlotService(slotRepo repository.SlotRepository, scheduleRepo repository.ScheduleRepository) SlotService {
    return &slotService{slotRepo: slotRepo, scheduleRepo: scheduleRepo}
}

func (s *slotService) GetAvailableSlots(roomID string, dateStr string) ([]models.Slot, error) {
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return nil, errors.New("invalid date format, use YYYY-MM-DD")
    }
    // Ensure slots exist for this date (generate on the fly)
    if err := s.generateSlotsForDate(roomID, date); err != nil {
        log.Printf("generateSlotsForDate error: %v", err)
        // If room has no schedule, just return empty list
        return []models.Slot{}, nil
    }
    return s.slotRepo.FindAvailableByRoomAndDate(roomID, date)
}

func (s *slotService) generateSlotsForDate(roomID string, date time.Time) error {
    schedule, err := s.scheduleRepo.FindByRoomID(roomID)
    if err != nil {
        log.Printf("No schedule for room %s: %v", roomID, err)
        return errors.New("no schedule")
    }
    // Check if day of week matches
    weekday := int(date.Weekday())
    ourWeekday := weekday
    if ourWeekday == 0 {
        ourWeekday = 7
    }
    log.Printf("Date: %s, weekday: %d, ourWeekday: %d, schedule days: %v", date.Format("2006-01-02"), weekday, ourWeekday, schedule.DaysOfWeek)
    match := false
    for _, d := range schedule.DaysOfWeek {
        if int(d) == ourWeekday {
            match = true
            break
        }
    }
    if !match {
        log.Printf("Day of week not in schedule")
        return errors.New("no schedule for this day")
    }

    // Обрезаем время до HH:MM (удаляем секунды, если они есть)
    startTimeStr := schedule.StartTime
    if len(startTimeStr) > 5 {
        startTimeStr = startTimeStr[:5]
    }
    endTimeStr := schedule.EndTime
    if len(endTimeStr) > 5 {
        endTimeStr = endTimeStr[:5]
    }

    startTime, err := time.Parse("15:04", startTimeStr)
    if err != nil {
        log.Printf("Error parsing startTime: %v", err)
        return err
    }
    endTime, err := time.Parse("15:04", endTimeStr)
    if err != nil {
        log.Printf("Error parsing endTime: %v", err)
        return err
    }

    slotStart := time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
    slotEnd := time.Date(date.Year(), date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

    log.Printf("Generating slots from %s to %s", slotStart, slotEnd)

    var slots []models.Slot
    for t := slotStart; t.Before(slotEnd); t = t.Add(30 * time.Minute) {
        end := t.Add(30 * time.Minute)
        if end.After(slotEnd) {
            break
        }
        exists, _ := s.slotRepo.CheckExists(roomID, t, end)
        if !exists {
            slots = append(slots, models.Slot{
                RoomID: roomID,
                Start:  t,
                End:    end,
            })
        }
    }
    if len(slots) > 0 {
        log.Printf("Creating %d slots", len(slots))
        return s.slotRepo.CreateInBatch(slots)
    }
    log.Printf("No new slots to create")
    return nil
}