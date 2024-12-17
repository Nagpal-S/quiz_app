package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user model in the database
type User struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	Name     string    `json:"name" gorm:"type:varchar(22)"`
	Email    string    `json:"email" gorm:"type:varchar(55);"`
	Phone    string    `json:"phone" gorm:"type:varchar(22)unique"`
	Otp      string    `json:"otp" gorm:"type:int(11)"`
	Image    string    `json:"image" gorm:"type:text"`
	Register string    `json:"register" gorm:"type:varchar(1);default:'0'"`                      // Enum "0" or "1"
	Gender   string    `json:"gender" gorm:"type:varchar(6);check:gender IN ('Male', 'Female')"` // Enum type for Gender
	Created  time.Time `json:"created" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// MigrateUser migrates the User table
func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
	// if err := db.AutoMigrate(&User{}); err != nil {
	// 	log.Println("Migration failed:", err)
	// }
}
