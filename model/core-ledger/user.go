package model

import "time"

type User struct {
	ID          uint64 `gorm:"primaryKey"`
	Email       string `gorm:"unique;not null"`
	Password    string
	FullName    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
