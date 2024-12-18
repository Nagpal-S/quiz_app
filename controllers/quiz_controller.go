package controllers

import (
	"quizapp/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuizController struct {
	DB *gorm.DB
}

// GetQuizCategories This API will provide list of quiz categories
//
//	@Summary		This API will provide list of quiz categories
//	@Description	This API will provide list of quiz categories
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	categoryInfo
//	@Router		/quizes/get-categories [get]
func (qc *QuizController) GetQuizCategories(c *gin.Context) {

	var quizCategories []models.QuizCategory

	if err := qc.DB.Where("active = 1").Find(&quizCategories).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "Database error occured.",
		})
		return

	}

	if len(quizCategories) == 0 {

		c.JSON(404, gin.H{

			"status":  "0",
			"message": "Category list data found empty.",
		})
		return

	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Category list data found.",
		"details": quizCategories,
	})

}

type categoryInfo struct {
	Details struct {
		ID                   int    `json:"id" example:"1"`
		Active               string `json:"active" example:"1"`
		Title                string `json:"title" example:"GK"`
		TotalPrice           int    `json:"total_price" example:"100000"`
		Icon                 string `json:"icon" example:"https://quizbuck.s3.ap-south-1.amazonaws.com/uploads/1734090491_new.jpg"`
		NumOfUsersCanJoin    int    `json:"num_of_users_can_join" example:"20"`
		NumOfUsersHaveJoined int    `json:"num_of_users_have_joined" example:"0"`
		QuizTime             string `json:"quiz_time" example:"2024-12-17T18:00:00+05:30"`
		JoinAmount           int    `json:"join_amount" example:"100"`
		Created              string `json:"created" example:"2024-12-17T18:07:19+05:30"`
	} `json:"details"`
	Message string `json:"message" example:"Category list data found."`
	Status  string `json:"status" example:"1"`
}

// GetQuizByCategory This API will provide list of quiz questions
//
//	@Summary		This API will provide list of quiz questions
//	@Description	This API will provide list of quiz questions
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		category_id	path		string	true	"quiz category id"
//	@Success	200	{object}	quizInfo
//	@Router		/quizes/get-question-by-category/{category_id} [get]
func (qc *QuizController) GetQuizByCategory(c *gin.Context) {

	categoryId := c.Param("category_id")

	var quiz []models.QuizQuestion

	if err := qc.DB.Where("category_id = ?", categoryId).Find(&quiz).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "Server error while fetching the data.",
		})
		return

	}

	if len(quiz) == 0 {

		c.JSON(404, gin.H{

			"status":  "0",
			"message": "Quiz question list found empty.",
		})
		return

	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Quiz question list found.",
		"details": quiz,
	})

}

type quizInfo struct {
	Details struct {
		ID            int       `json:"id" example:"1"`
		CategoryID    int       `json:"category_id" example:"1"`
		Level         string    `json:"level" example:"easy"`
		Question      string    `json:"question" example:"Where is Delhi?"`
		OptionA       string    `json:"option_a" example:"Unites States of America"`
		OptionB       string    `json:"option_b" example:"England"`
		OptionC       string    `json:"option_c" example:"India"`
		OptionD       string    `json:"option_d" example:"Sri Lanka"`
		CorrectAnswer string    `json:"correct_answer" example:"c"`
		CreatedAt     time.Time `created_at" example:"2024-12-18T14:16:53+05:30"`
	} `json:"details"`
	Message string `json:"message" example:"New user created"`
	Status  string `json:"status" example:"1"`
}
