package routes

import (
	"quizapp/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all application routes.
func RegisterRoutes(router *gin.Engine) {
	fileController := new(controllers.FileController)

	// Define the upload route
	api := router.Group("/upload")
	{
		api.POST("/file", fileController.UploadToS3)
	}
}
