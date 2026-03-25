package models

import "time"

type Account struct {
	ID 			uint 	`gorm:"primaryKey"`
	Email		string 	`gorm:"unique; not null"`
	Username	string 	`gorm:"not null"`
	Password 	string 	`gorm:"not null"`
	ColorScheme	string 	`gorm:"not null"`
	CreatedAt	time.Time 
}