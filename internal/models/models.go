package models

import (
    "time"
    "github.com/lib/pq"
)

type User struct {
    ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"`
    Role      string    `gorm:"type:varchar(10);not null" json:"role"`
    CreatedAt time.Time `json:"createdAt"`
}

type Room struct {
    ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Name        string     `gorm:"not null" json:"name"`
    Description *string    `json:"description"`
    Capacity    *int       `json:"capacity"`
    CreatedAt   time.Time  `json:"createdAt"`
}

type Schedule struct {
    ID         string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    RoomID     string        `gorm:"type:uuid;not null;uniqueIndex" json:"roomId"`
    DaysOfWeek pq.Int64Array `gorm:"type:integer[]" json:"daysOfWeek"`
    StartTime  string        `gorm:"type:time;not null" json:"startTime"`
    EndTime    string        `gorm:"type:time;not null" json:"endTime"`
}

type Slot struct {
    ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    RoomID    string    `gorm:"type:uuid;not null;index" json:"roomId"`
    Start     time.Time `gorm:"not null;uniqueIndex:idx_slot_room_start" json:"start"`
    End       time.Time `gorm:"not null" json:"end"`
}

type Booking struct {
    ID             string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    SlotID         string     `gorm:"type:uuid;not null;uniqueIndex:idx_active_slot,where:status='active'" json:"slotId"`
    UserID         string     `gorm:"type:uuid;not null;index" json:"userId"`
    Status         string     `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
    ConferenceLink *string    `json:"conferenceLink"`
    CreatedAt      time.Time  `json:"createdAt"`
}