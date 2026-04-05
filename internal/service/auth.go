package service

import (
    "errors"
    "time"

    "booking/internal/models"
    "booking/internal/repository"
    "booking/internal/utils"

    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    DummyLogin(role string) (string, error)
    Register(email, password, role string) (*models.User, error)
    Login(email, password string) (string, error)
}

type authService struct {
    userRepo  repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
    return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) DummyLogin(role string) (string, error) {
    var userID, email string
    if role == "admin" {
        userID = "11111111-1111-1111-1111-111111111111"
        email = "admin@example.com"
    } else if role == "user" {
        userID = "22222222-2222-2222-2222-222222222222"
        email = "user@example.com"
    } else {
        return "", errors.New("invalid role")
    }

    // Ensure user exists (should be created by EnsureDummyUsers)
    _, err := s.userRepo.FindByID(userID)
    if err != nil {
        // Create if not exists
        user := &models.User{
            ID:        userID,
            Email:     email,
            Password:  "$2a$10$dummyhashdummyhashdummyhashdummyhashdummyhash",
            Role:      role,
            CreatedAt: time.Now(),
        }
        if err := s.userRepo.Create(user); err != nil {
            return "", err
        }
    }

    token, err := utils.GenerateJWT(userID, role, s.jwtSecret)
    return token, err
}

func (s *authService) Register(email, password, role string) (*models.User, error) {
    // Check if user exists
    existing, _ := s.userRepo.FindByEmail(email)
    if existing != nil && existing.ID != "" {
        return nil, errors.New("email already registered")
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        Email:     email,
        Password:  string(hashed),
        Role:      role,
        CreatedAt: time.Now(),
    }
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    token, err := utils.GenerateJWT(user.ID, user.Role, s.jwtSecret)
    return token, err
}