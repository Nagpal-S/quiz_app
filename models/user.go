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
	Image    string    `json:"image" gorm:"type:text; default:NULL"`
	Register string    `json:"register" gorm:"type:varchar(1);default:'0'"` // Enum "0" or "1"
	Gender   string    `json:"gender" gorm:"type:varchar(6);default:Male;check:gender IN ('Male', 'Female', 'Others')"`
	Created  time.Time `json:"created" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// MigrateUser migrates the User table
func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
	// if err := db.AutoMigrate(&User{}); err != nil {
	// 	log.Println("Migration failed:", err)
	// }
}

type UserTransactions struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	UserId          uint      `gorm:"not null;index" json:"user_id"` // Foreign key
	Title           string    `json:"title" gorm:"type:varchar(44)"`
	TransactionType string    `gorm:"type:enum('CREDIT','DEBIT');not null; check:transaction_type IN ('DEBIT', 'CREDIT')" json:"transaction_type"`
	Amount          float64   `json:"amount" gorm:"type:float"`
	Created         time.Time `json:"created" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

func MigrateUserTransactions(db *gorm.DB) {
	db.AutoMigrate(&UserTransactions{})
	// if err := db.AutoMigrate(&User{}); err != nil {
	// 	log.Println("Migration failed:", err)
	// }
}
