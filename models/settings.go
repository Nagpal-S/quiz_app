package models

import (
	"time"

	"gorm.io/gorm"
)

type Banners struct {
	ID      uint64    `gorm:"primaryKey;autoIncrement"`
	Banner  string    `json:"banner" gorm:"type:varchar(555)`
	Created time.Time `json:"crated" gorm:"type:datetime"`
}

func MigrateBanners(db *gorm.DB) {
	db.AutoMigrate(&Banners{})
}
