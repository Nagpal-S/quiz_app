package routes

import (
	"quizapp/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SettingsRoutes(router *gin.Engine, db *gorm.DB) {

	SettingController := controllers.SettingController{DB: db}

	quizGroup := router.Group("/settings")
	{

		quizGroup.GET("get-banners-list", SettingController.GetBannersList)

	}

}
