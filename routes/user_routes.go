package routes

import (
	"quizapp/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(router *gin.Engine, db *gorm.DB) {

	userController := controllers.UserController{DB: db}

	userGroup := router.Group("/users")
	{
		userGroup.GET("check-user/:phone", userController.GetUserByPhone)
		userGroup.GET("/verify-otp/:userId/:otp", userController.VerifyOTP)
		userGroup.GET("/get-user/:id", userController.CheckUser)
		userGroup.POST("/edit-user-profile/", userController.EditUserProfile)
		userGroup.POST("/initiate-user-transaction/", userController.InitiateUserTransaction)
		userGroup.GET("/get-user-wallet-details/:user_id", userController.GetUserWalletDetails)
	}
}
