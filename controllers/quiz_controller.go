package controllers

import (
	"errors"
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
//	@Param			user_id	path	string	true	"User ID"
//	@Success	200	{object}	categoryInfo
//	@Router		/quizes/get-categories/{user_id} [get]
func (qc *QuizController) GetQuizCategories(c *gin.Context) {
	userID := c.Param("user_id")

	var quizCategories []models.QuizCategory

	// Fetch active quiz categories where quiz_time is greater than or equal to the current time
	if err := qc.DB.Where("active = 1 AND quiz_time >= ?", time.Now()).Find(&quizCategories).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "0",
			"message": "Database error occurred.",
		})
		return
	}

	// If no categories found, return a 404 response
	if len(quizCategories) == 0 {
		c.JSON(404, gin.H{
			"status":  "0",
			"message": "Category list data found empty.",
		})
		return
	}

	// Prepare the response with user participation status
	var response []gin.H
	for _, quiz := range quizCategories {
		var userJoinInfo models.UserJoinContest

		// Check if the user has already joined the quiz for the current category
		err := qc.DB.Where("user_id = ? AND category_id = ?", userID, quiz.ID).First(&userJoinInfo).Error
		hasJoined := err == nil // User has joined if no error is returned

		response = append(response, gin.H{
			"id":                       quiz.ID,
			"active":                   quiz.Active,
			"title":                    quiz.Title,
			"total_price":              quiz.TotalPrice,
			"icon":                     quiz.Icon,
			"num_of_users_can_join":    quiz.NumOfUsersCanJoin,
			"num_of_users_have_joined": quiz.NumOfUsersHaveJoined,
			"quiz_time":                quiz.QuizTime,
			"join_amount":              quiz.JoinAmount,
			"created":                  quiz.Created,
			"user_has_joined":          hasJoined, // Include user participation status
		})
	}

	// Return the response
	c.JSON(200, gin.H{
		"status":  "1",
		"message": "Category list data found.",
		"details": response,
	})
}

// categoryInfo Struct to represent the API response for Swagger documentation.
type categoryInfo struct {
	Details []struct {
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
		UserHasJoined        bool   `json:"user_has_joined" example:"true"`
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
//	@Success	200	{object}	QuizResponse
//	@Router		/quizes/get-question-by-category/{category_id} [get]
func (qc *QuizController) GetQuizByCategory(c *gin.Context) {

	categoryId := c.Param("category_id")

	var quiz []models.QuizQuestion
	var category models.QuizCategory

	if err := qc.DB.Where("id = ? AND active = 1 AND quiz_time <= ?", categoryId, time.Now()).First(&category).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(404, gin.H{

				"status":  "0",
				"message": "Category not found or category not exist.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching category info.",
			})

		}

		return

	}

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

	var response = gin.H{
		"category_name":      category.Title,
		"category_id":        category.ID,
		"quiz_start_time":    category.QuizTime,
		"quiz_end_time":      category.QuizEndTime,
		"questions_duration": category.EachQuestionTimeDuration,
		"questions":          quiz,
	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Quiz question list found.",
		"details": response,
	})

}

type QuizResponse struct {
	Details struct {
		CategoryID        int        `json:"category_id" example:"1"`
		CategoryName      string     `json:"category_name" example:"GK"`
		Questions         []Question `json:"questions"`
		QuestionsDuration int        `json:"questions_duration" example:"60"`
		QuizStartTime     string     `json:"quiz_start_time" example:"2024-12-24T14:00:00+05:30"`
		QuizEndTime       string     `json:"quiz_end_time" example:"2024-12-24T18:57:33+05:30"`
	} `json:"details"`
	Message string `json:"message" example:"Quiz question list found."`
	Status  string `json:"status" example:"1"`
}

type Question struct {
	ID            int    `json:"id" example:"1"`
	CategoryID    int    `json:"category_id" example:"1"`
	Level         string `json:"level" example:"easy"`
	Question      string `json:"question" example:"Where is Delhi?"`
	OptionA       string `json:"option_a" example:"Unites States of America"`
	OptionB       string `json:"option_b" example:"England"`
	OptionC       string `json:"option_c" example:"India"`
	OptionD       string `json:"option_d" example:"Sri Lanka"`
	CorrectAnswer string `json:"correct_answer" example:"c"`
	CreatedAt     string `json:"created_at" example:"2024-12-18T14:16:53+05:30"`
}

// UserJoinContest This API will make user to join a contest
//
//	@Summary		This API will make user to join a contest
//	@Description	This API will make user to join a contest
//	@Schemes
//	@Tags		Quizes
//	@Accept		application/x-www-form-urlencoded
//	@Produce	json
//	@Param		user_id	formData		int	true	"user id"
//	@Param		category_id	formData		string	true	"quiz category id"
//	@Success	200	{object}	JoinContestResponse
//	@Router		/quizes/user-join-contest [post]
func (qc *QuizController) UserJoinContest(c *gin.Context) {

	userId := c.PostForm("user_id")
	categoryId := c.PostForm("category_id")

	var user models.User
	var category models.QuizCategory
	var wallet models.UserWallet
	var userJoinInfo models.UserJoinContest
	var UserTransactions models.UserTransactions

	// user info
	if err := qc.DB.Where("id = ? AND register = '1'", userId).First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(404, gin.H{

				"status":  "0",
				"message": "Invalid user_id or user not found. " + err.Error(),
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching user info. " + err.Error(),
			})

		}

		return
	}

	// category info
	if err := qc.DB.Where("id = ? AND active = 1", categoryId).First(&category).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(404, gin.H{

				"status":  "0",
				"message": "Invalid category_id or category not found. " + err.Error(),
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching category details. " + err.Error(),
			})

		}

		return

	}

	// user wallet info
	if err := qc.DB.Where("user_id = ?", userId).First(&wallet).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(404, gin.H{

				"status":  "0",
				"message": "User wallet not found. " + err.Error(),
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while getting wallet details. " + err.Error(),
			})

		}

		return

	}

	// check if user has already joined or not
	if err := qc.DB.Where("user_id = ? AND  category_id = ?", userId, categoryId).First(&userJoinInfo).Error; err != nil {

	} else {

		c.JSON(422, gin.H{

			"status":  "0",
			"message": "Can not join contest again.",
		})
		return

	}

	if category.JoinAmount > int(wallet.Amount) {

		c.JSON(422, gin.H{

			"status":  "0",
			"message": "Your wallet balance is too low to join the contest. Please recharge your wallet.",
		})
		return

	}

	if category.NumOfUsersCanJoin == category.NumOfUsersHaveJoined {

		c.JSON(422, gin.H{

			"status":  "0",
			"message": "Unable to join contest, slots not available.",
		})
		return

	}

	category.NumOfUsersHaveJoined += 1
	wallet.Amount -= float64(category.JoinAmount)

	// update wallet
	if err := qc.DB.Model(&wallet).Update("amount", wallet.Amount).Error; err != nil {
		c.JSON(500, gin.H{
			"status":  "0",
			"message": "DB error while updating wallet",
		})
		return
	}

	userJoinInfo.CategoryID = category.ID
	userJoinInfo.UserID = user.ID

	// save user info contest
	if err := qc.DB.Create(&userJoinInfo).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while saving join info. " + err.Error(),
		})
		return

	}

	if err := qc.DB.Model(&category).Update("NumOfUsersHaveJoined", category.NumOfUsersHaveJoined).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while updating contest details. " + err.Error(),
		})
		return

	}

	UserTransactions.Amount = float64(category.JoinAmount)
	UserTransactions.UserId = user.ID
	UserTransactions.TransactionType = "DEBIT"
	UserTransactions.Title = category.Title + " contest joined."

	if err := qc.DB.Create(&UserTransactions).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while adding transaction info. " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Contest joined successfully.",
	})

}

type JoinContestResponse struct {
	Status  string `json:"status" example:"1"`
	Message string `json:"message" example:"Contest joined successfully."`
}

// GetContestJoinedByUser This API will list contest joined by user
//
//	@Summary		This API will list contest joined by user
//	@Description	This API will list contest joined by user
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		user_id	path		string	true	"user id"
//	@Success	200	{object}	JoinedContestResponse
//	@Router		/quizes/get-contest-joined-by-user/{user_id} [get]
func (qc *QuizController) GetContestJoinedByUser(c *gin.Context) {

	userId := c.Param("user_id")

	var user models.User
	var userJoinInfo []models.UserJoinContest

	if err := qc.DB.Where("id = ? AND register = 1", userId).First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(404, gin.H{

				"status":  "0",
				"message": "Invalid user_id or user not found. " + err.Error(),
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching user info. " + err.Error(),
			})

		}
		return

	}

	if err := qc.DB.Where("user_id = ?", userId).Find(&userJoinInfo).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while fetching user contest join info. " + err.Error(),
		})
		return

	}

	if len(userJoinInfo) == 0 {

		c.JSON(404, gin.H{

			"status":  "0",
			"message": "Contest list found empty.",
		})
		return

	}

	var response []gin.H
	for _, contests := range userJoinInfo {

		var category models.QuizCategory

		err := qc.DB.Where("id = ?", contests.CategoryID).First(&category).Error
		if err != nil {
			continue
		}

		response = append(response, gin.H{

			"contest_name":   category.Title,
			"contest_date":   category.QuizTime,
			"contest_amount": category.TotalPrice,
			"contest_id":     category.ID,
		})

	}

	c.JSON(200, gin.H{
		"status":  "1",
		"message": "Contest list found.",
		"details": response,
	})

}

type JoinedContestResponse struct {
	Status  string        `json:"status" example:"1"`
	Message string        `json:"message" example:"Contest list found."`
	Details []ContestInfo `json:"details"`
}

type ContestInfo struct {
	ContestName   string `json:"contest_name" example:"GK"`
	ContestDate   string `json:"contest_date" example:"2024-12-21T18:00:00+05:30"`
	ContestAmount int    `json:"contest_amount" example:"10000"`
	ContestID     int    `json:"contest_id" example:"1"`
}

// GetRulesByCategory This API will provide contest rules by category/contest id
//
//	@Summary		This API will provide contest rules by category/contest id
//	@Description	This API will provide contest rules by category/contest id
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		category_id	path		string	true	"category_id id"
//	@Success	200	{object}	RulesResponse
//	@Router		/quizes/get-rules-list-by-category/{category_id} [get]
func (qc *QuizController) GetRulesByCategory(c *gin.Context) {

	categoryId := c.Param("category_id")

	var rules []models.ContestRules

	if err := qc.DB.Where("category_id = ?", categoryId).Find(&rules).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while getting rules list.",
		})
		return

	}

	if len(rules) == 0 {

		c.JSON(404, gin.H{

			"status":  "0",
			"message": "Rules list found empty.",
		})
		return
	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Rules list found",
		"details": rules,
	})

}

type RulesResponse struct {
	Details []Rule `json:"details"`
	Message string `json:"message" example:"Rules list found"`
	Status  string `json:"status" example:"1"`
}

type Rule struct {
	ID         int    `json:"id" example:"1"`
	CategoryID int    `json:"category_id" example:"1"`
	Rule       string `json:"rule" example:"Complete question on time"`
	CreatedAt  string `json:"created_at" example:"2024-12-25T16:11:37+05:30"`
}

// GetContestPrizes This API will provide contest prizes list by category/contest id
//
//	@Summary		This API will provide contest prizes list by category/contest id
//	@Description	This API will provide contest prizes list by category/contest id
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		category_id	path		string	true	"category_id id"
//	@Success	200	{object}	PrizesResponse
//	@Router		/quizes/get-contest-prize-list-by-category/{category_id} [get]
func (qc *QuizController) GetContestPrizes(c *gin.Context) {

	categoryId := c.Param("category_id")

	var prizes []models.ContestPrize

	if err := qc.DB.Where("category_id = ?", categoryId).Find(&prizes).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while fetching prizes data.",
		})
		return

	}

	if len(prizes) == 0 {

		c.JSON(404, gin.H{

			"status":  "0",
			"message": "Prizes list found empty for this contest.",
		})
		return

	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Prizes list found.",
		"details": prizes,
	})

}

type PrizesResponse struct {
	Details []Prize `json:"details"`
	Message string  `json:"message" example:"Prizes list found."`
	Status  string  `json:"status" example:"1"`
}

type Prize struct {
	ID         int     `json:"id" example:"1"`
	CategoryID int     `json:"category_id" example:"1"`
	RankFrom   int     `json:"rank_from" example:"1"`
	RankTo     int     `json:"rank_to" example:"1"`
	Winning    float64 `json:"winning" example:"25000"`
	CreatedAt  string  `json:"created_at" example:"2024-12-26T19:08:03+05:30"`
}
