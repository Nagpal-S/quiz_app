package routes

import (
	"quizapp/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes sets up all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {

	SocketController := controllers.SocketController{DB: db}

	socketGroup := router.Group("/ws")
	{
		socketGroup.GET("/", SocketController.HandleWebSocket)

	}

}
