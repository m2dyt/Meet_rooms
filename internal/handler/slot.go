package handler

import (
    "net/http"

    "github.com/gorilla/mux"
    "booking/internal/service"
)

type SlotHandler struct {
    slotService service.SlotService
}

func NewSlotHandler(slotService service.SlotService) *SlotHandler {
    return &SlotHandler{slotService: slotService}
}

// ListSlots – список доступных слотов (admin и user)
func (h *SlotHandler) ListSlots(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    roomID := vars["roomId"]
    if roomID == "" {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "roomId is required")
        return
    }
    date := r.URL.Query().Get("date")
    if date == "" {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "date parameter is required")
        return
    }
    slots, err := h.slotService.GetAvailableSlots(roomID, date)
    if err != nil {
        if err.Error() == "room not found" {
            writeError(w, http.StatusNotFound, "ROOM_NOT_FOUND", err.Error())
        } else {
            writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        }
        return
    }
    writeJSON(w, http.StatusOK, map[string]interface{}{"slots": slots})
}