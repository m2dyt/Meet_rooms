package handler

import (
    "encoding/json"
    "net/http"

    "booking/internal/service"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

// DummyLogin godoc
// @Summary      Тестовый JWT
// @Description  Возвращает JWT для указанной роли (admin/user). Фиксированные UUID.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body object true "Роль" SchemaExample({"role":"admin"})
// @Success      200  {object}  map[string]interface{} "token"
// @Failure      400  {object}  map[string]interface{} "error"
// @Router       /dummyLogin [post]
func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Role string `json:"role"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    token, err := h.authService.DummyLogin(req.Role)
    if err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

// Register godoc
// @Summary      Регистрация
// @Description  Регистрирует нового пользователя (email, пароль, роль)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body object true "Данные регистрации" SchemaExample({"email":"user@example.com","password":"123","role":"user"})
// @Success      201  {object}  map[string]interface{} "user"
// @Failure      400  {object}  map[string]interface{} "error"
// @Router       /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    user, err := h.authService.Register(req.Email, req.Password, req.Role)
    if err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    writeJSON(w, http.StatusCreated, map[string]interface{}{"user": user})
}

// Login godoc
// @Summary      Вход
// @Description  Авторизует пользователя по email и паролю, возвращает JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body object true "Учётные данные" SchemaExample({"email":"user@example.com","password":"123"})
// @Success      200  {object}  map[string]interface{} "token"
// @Failure      401  {object}  map[string]interface{} "error"
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
        return
    }
    token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"token": token})
}