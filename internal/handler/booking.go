package handler

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "booking/internal/middleware"
    "booking/internal/service"
)

type BookingHandler struct {
    bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
    return &BookingHandler{bookingService: bookingService}
}

// CreateBooking – создание брони (только user)
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
    role, ok := r.Context().Value(middleware.RoleKey).(string)
    if !ok || role != "user" {
        writeError(w, http.StatusForbidden, "FORBIDDEN", "only users can create bookings")
        return
    }
    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user")
        return
    }

    var req struct {
        SlotID               string `json:"slotId"`
        CreateConferenceLink bool   `json:"createConferenceLink"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    if req.SlotID == "" {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "slotId is required")
        return
    }
    booking, err := h.bookingService.CreateBooking(userID, req.SlotID, req.CreateConferenceLink)
    if err != nil {
        switch err.Error() {
        case "slot not found":
            writeError(w, http.StatusNotFound, "SLOT_NOT_FOUND", err.Error())
        case "slot is already booked":
            writeError(w, http.StatusConflict, "SLOT_ALREADY_BOOKED", err.Error())
        case "cannot book slot in the past":
            writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        default:
            writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        }
        return
    }
    writeJSON(w, http.StatusCreated, map[string]interface{}{"booking": booking})
}

// CancelBooking – отмена брони (только владелец)
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
    role, ok := r.Context().Value(middleware.RoleKey).(string)
    if !ok || role != "user" {
        writeError(w, http.StatusForbidden, "FORBIDDEN", "only users can cancel bookings")
        return
    }
    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user")
        return
    }

    vars := mux.Vars(r)
    bookingID := vars["bookingId"]
    if bookingID == "" {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "bookingId is required")
        return
    }
    booking, err := h.bookingService.CancelBooking(bookingID, userID)
    if err != nil {
        switch err.Error() {
        case "booking not found":
            writeError(w, http.StatusNotFound, "BOOKING_NOT_FOUND", err.Error())
        case "cannot cancel another user's booking":
            writeError(w, http.StatusForbidden, "FORBIDDEN", err.Error())
        default:
            writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        }
        return
    }
    writeJSON(w, http.StatusOK, map[string]interface{}{"booking": booking})
}

// ListAllBookings – список всех броней (только admin)
func (h *BookingHandler) ListAllBookings(w http.ResponseWriter, r *http.Request) {
    role, ok := r.Context().Value(middleware.RoleKey).(string)
    if !ok || role != "admin" {
        writeError(w, http.StatusForbidden, "FORBIDDEN", "only admin can list all bookings")
        return
    }
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
    if pageSize < 1 {
        pageSize = 20
    }
    if pageSize > 100 {
        pageSize = 100
    }
    bookings, total, err := h.bookingService.ListAllBookings(page, pageSize)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]interface{}{
        "bookings": bookings,
        "pagination": map[string]interface{}{
            "page":     page,
            "pageSize": pageSize,
            "total":    total,
        },
    })
}

// MyBookings – брони текущего пользователя (только user)
func (h *BookingHandler) MyBookings(w http.ResponseWriter, r *http.Request) {
    role, ok := r.Context().Value(middleware.RoleKey).(string)
    if !ok || role != "user" {
        writeError(w, http.StatusForbidden, "FORBIDDEN", "only users can view their bookings")
        return
    }
    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid user")
        return
    }
    bookings, err := h.bookingService.MyBookings(userID)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]interface{}{"bookings": bookings})
}