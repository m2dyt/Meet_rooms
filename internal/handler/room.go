package handler

import (
    "encoding/json"
    "net/http"

    "booking/internal/middleware"
    "booking/internal/service"
)

type RoomHandler struct {
    roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
    return &RoomHandler{roomService: roomService}
}

// ListRooms – список переговорок (доступно всем авторизованным)
func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
    rooms, err := h.roomService.ListRooms()
    if err != nil {
        writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]interface{}{"rooms": rooms})
}

// CreateRoom – создание переговорки (только admin)
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
    // Получаем роль из контекста
    role, ok := r.Context().Value(middleware.RoleKey).(string)
    if !ok || role != "admin" {
        writeError(w, http.StatusForbidden, "FORBIDDEN", "only admin can create rooms")
        return
    }

    var req struct {
        Name        string  `json:"name"`
        Description *string `json:"description"`
        Capacity    *int    `json:"capacity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    if req.Name == "" {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "name is required")
        return
    }
    room, err := h.roomService.CreateRoom(req.Name, req.Description, req.Capacity)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }
    writeJSON(w, http.StatusCreated, map[string]interface{}{"room": room})
}