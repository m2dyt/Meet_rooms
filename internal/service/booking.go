package service

import (
    "errors"
    "time"

    "booking/internal/models"
    "booking/internal/repository"
)

type BookingService interface {
    CreateBooking(userID, slotID string, createConferenceLink bool) (*models.Booking, error)
    CancelBooking(bookingID, userID string) (*models.Booking, error)
    ListAllBookings(page, pageSize int) ([]models.Booking, int64, error)
    MyBookings(userID string) ([]models.Booking, error)
}

type bookingService struct {
    bookingRepo repository.BookingRepository
    slotRepo    repository.SlotRepository
    userRepo    repository.UserRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, slotRepo repository.SlotRepository, userRepo repository.UserRepository) BookingService {
    return &bookingService{
        bookingRepo: bookingRepo,
        slotRepo:    slotRepo,
        userRepo:    userRepo,
    }
}

func (s *bookingService) CreateBooking(userID, slotID string, createConferenceLink bool) (*models.Booking, error) {
    slot, err := s.slotRepo.FindByID(slotID)
    if err != nil {
        return nil, errors.New("slot not found")
    }
    if slot.Start.Before(time.Now().UTC()) {
        return nil, errors.New("cannot book slot in the past")
    }
    existing, _ := s.bookingRepo.FindBySlotIDAndActive(slotID)
    if existing != nil && existing.ID != "" {
        return nil, errors.New("slot is already booked")
    }
    _, err = s.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }
    booking := &models.Booking{
        SlotID: slotID,
        UserID: userID,
        Status: "active",
    }
    if createConferenceLink {
        link := "https://meet.example.com/" + slotID[:8]
        booking.ConferenceLink = &link
    }
    if err := s.bookingRepo.Create(booking); err != nil {
        return nil, err
    }
    return booking, nil
}

func (s *bookingService) CancelBooking(bookingID, userID string) (*models.Booking, error) {
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        return nil, errors.New("booking not found")
    }
    if booking.UserID != userID {
        return nil, errors.New("cannot cancel another user's booking")
    }
    if booking.Status == "cancelled" {
        return booking, nil
    }
    if err := s.bookingRepo.UpdateStatus(bookingID, "cancelled"); err != nil {
        return nil, err
    }
    booking.Status = "cancelled"
    return booking, nil
}

func (s *bookingService) ListAllBookings(page, pageSize int) ([]models.Booking, int64, error) {
    return s.bookingRepo.ListAll(page, pageSize)
}

func (s *bookingService) MyBookings(userID string) ([]models.Booking, error) {
    return s.bookingRepo.ListByUserFuture(userID)
}