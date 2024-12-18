package models

import (
	"time"

	"gorm.io/gorm"
)

type UserWallet struct {
	ID      uint64    `gorm:"primaryKey;autoIncrement"`
	UserId  uint64    `json:"user_id" gorm:"type:int(64)"`
	Amount  float64   `json:"amount" gorm:"type:float"`
	Created time.Time `json:"crated" gorm:"type:datetime"`
}

func MigrateUserWallet(db *gorm.DB) {
	db.AutoMigrate(&UserWallet{})
}
