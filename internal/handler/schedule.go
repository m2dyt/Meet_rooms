package handler

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "booking/internal/middleware"
    "booking/internal/service"
)

type ScheduleHandler struct {
    scheduleService service.ScheduleService
}

func NewScheduleHandler(scheduleService service.ScheduleService) *ScheduleHandler {
    return &ScheduleHandler{scheduleService: scheduleService}
}

// CreateSchedule – создание расписания (только admin)
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
    role, ok := r.Context().Value(middleware.RoleKey).(string)
    if !ok || role != "admin" {
        writeError(w, http.StatusForbidden, "FORBIDDEN", "only admin can create schedules")
        return
    }
    vars := mux.Vars(r)
    roomID := vars["roomId"]
    if roomID == "" {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "roomId is required")
        return
    }
    var req struct {
        DaysOfWeek []int  `json:"daysOfWeek"`
        StartTime  string `json:"startTime"`
        EndTime    string `json:"endTime"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    schedule, err := h.scheduleService.CreateSchedule(roomID, req.DaysOfWeek, req.StartTime, req.EndTime)
    if err != nil {
        switch err.Error() {
        case "room not found":
            writeError(w, http.StatusNotFound, "ROOM_NOT_FOUND", err.Error())
        case "schedule already exists for this room":
            writeError(w, http.StatusConflict, "SCHEDULE_EXISTS", err.Error())
        default:
            writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        }
        return
    }
    writeJSON(w, http.StatusCreated, map[string]interface{}{"schedule": schedule})
}