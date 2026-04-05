package database

import (
    "fmt"
    "log"
    "time"

    "booking/internal/config"
    "booking/internal/models"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(time.Hour)
    return db, nil
}

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Room{},
        &models.Schedule{},
        &models.Slot{},
        &models.Booking{},
    )
}

func EnsureDummyUsers(db *gorm.DB) error {
    adminID := "11111111-1111-1111-1111-111111111111"
    userID := "22222222-2222-2222-2222-222222222222"

    var admin models.User
    if err := db.Where("id = ?", adminID).First(&admin).Error; err != nil {
        admin = models.User{
            ID:        adminID,
            Email:     "admin@example.com",
            Password:  "$2a$10$dummyhashdummyhashdummyhashdummyhashdummyhash",
            Role:      "admin",
            CreatedAt: time.Now(),
        }
        if err := db.Create(&admin).Error; err != nil {
            return err
        }
        log.Println("Created dummy admin user")
    }

    var user models.User
    if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
        user = models.User{
            ID:        userID,
            Email:     "user@example.com",
            Password:  "$2a$10$dummyhashdummyhashdummyhashdummyhashdummyhash",
            Role:      "user",
            CreatedAt: time.Now(),
        }
        if err := db.Create(&user).Error; err != nil {
            return err
        }
        log.Println("Created dummy regular user")
    }
    return nil
}

func Seed(db *gorm.DB) error {
    log.Println("Seeding test data...")
    if err := EnsureDummyUsers(db); err != nil {
        return err
    }

    // Create a test room if none exists
    var count int64
    db.Model(&models.Room{}).Count(&count)
    if count == 0 {
        desc := "Test meeting room"
        cap := 10
        room := models.Room{
            Name:        "Conference Room A",
            Description: &desc,
            Capacity:    &cap,
            CreatedAt:   time.Now(),
        }
        if err := db.Create(&room).Error; err != nil {
            return err
        }
        log.Printf("Created test room: %s", room.Name)
    }
    return nil
}