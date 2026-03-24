package models

import "time"

type Account struct {
	ID 			uint 	`gorm:"primaryKey"`
	Username	string 	`gorm:"unique;not null"`
	Password 	string 	`gorm:"not null"`
	ColorScheme	string 	`gorm:"not null"`
	CreatedAt	time.Time 
}