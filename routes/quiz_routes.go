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

		quizGroup.GET("get-categories/:user_id", QuizController.GetQuizCategories)
		quizGroup.GET("get-question-by-category/:category_id", QuizController.GetQuizByCategory)
		quizGroup.POST("user-join-contest", QuizController.UserJoinContest)
		quizGroup.GET("get-contest-joined-by-user/:user_id", QuizController.GetContestJoinedByUser)
		quizGroup.GET("get-rules-list-by-category/:category_id", QuizController.GetRulesByCategory)
		quizGroup.GET("get-contest-prize-list-by-category/:category_id", QuizController.GetContestPrizes)
		quizGroup.POST("user-question-answer/", QuizController.UserContestResponse)
		quizGroup.GET("get-user-contest-result/:user_id/:category_id", QuizController.GetUserContestReport)
		quizGroup.GET("get-contest-leaderboard/:category_id", QuizController.GetUserContestLeaderboard)
		quizGroup.GET("get-user-contest-history/:user_id", QuizController.GetUserPlayedContest)
		quizGroup.POST("create-leaderboard", QuizController.CreateLeaderboard)

	}

}
