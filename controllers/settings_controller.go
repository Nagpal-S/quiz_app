package controllers

import (
	"quizapp/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingController struct {
	DB *gorm.DB
}

// =============================================== GetBannersList start ========================================================

// GetBannersList This API will provide list of banners
//
//	@Summary		This API will provide list of banners
//	@Description	This API will provide list of banners
//	@Schemes
//	@Tags		Settings
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	GetBannersListResponse
//	@Router		/settings/get-banners-list [get]
func (db *SettingController) GetBannersList(c *gin.Context) {

	var Banners []models.Banners

	if err := db.DB.Find(&Banners).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while getting banner data. " + err.Error(),
		})

	}

	if len(Banners) == 0 {

		c.JSON(200, gin.H{

			"status":  "0",
			"message": "Banners list found empty.",
		})
		return

	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Banners list found.",
		"details": Banners,
	})

}

type GetBannersListResponse struct {
	Details []BannerDetails `json:"details"`                               // Array of banner details
	Message string          `json:"message" example:"Banners list found."` // Response message
	Status  string          `json:"status" example:"1"`                    // Status of the request
}

type BannerDetails struct {
	ID      int    `json:"ID" example:"1"`                             // Unique ID of the banner
	Banner  string `json:"banner" example:"image-url"`                 // URL of the banner image
	Created string `json:"crated" example:"2024-12-31T17:00:00+05:30"` // Timestamp of banner creation
}
