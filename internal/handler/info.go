package handler

import "net/http"

type InfoHandler struct{}

func NewInfoHandler() *InfoHandler {
    return &InfoHandler{}
}

// Info godoc
// @Summary      Health check
// @Description  Всегда возвращает 200 OK
// @Tags         System
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /_info [get]
func (h *InfoHandler) Info(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}