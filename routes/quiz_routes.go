package routes

import (
	"quizapp/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func QuizRoutes(router *gin.Engine, db *gorm.DB) {

	QuizController := controllers.QuizController{DB: db}

	quizGroup := router.Group("/quizes")
	{

		quizGroup.GET("get-categories", QuizController.GetQuizCategories)
		quizGroup.GET("get-question-by-category/:category_id", QuizController.GetQuizByCategory)

	}

}
