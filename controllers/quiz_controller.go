package controllers

import (
	"errors"
	"quizapp/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuizController struct {
	DB *gorm.DB
}

// =============================================== GetQuizCategories start ========================================================

// GetQuizCategories This API will provide list of quiz categories
//
//	@Summary		This API will provide list of quiz categories
//	@Description	This API will provide list of quiz categories
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param			user_id	path	string	true	"User ID"
//	@Success	200	{object}	CategoryInfo
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

	// If no categories found, return a 204 response
	if len(quizCategories) == 0 {
		c.JSON(204, gin.H{
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

type CategoryInfo struct {
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

// =============================================== GetQuizByCategory start ========================================================

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

			c.JSON(204, gin.H{

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

		c.JSON(204, gin.H{

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

// =============================================== UserJoinContest start ========================================================

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
	var userJoinContestHistory models.UserJoinContestHistory
	var UserTransactions models.UserTransactions

	// user info
	if err := qc.DB.Where("id = ? AND register = '1'", userId).First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(200, gin.H{

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

			c.JSON(200, gin.H{

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

			c.JSON(200, gin.H{

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

	if err := qc.DB.Where("user_id = ? AND  category_id = ?", userId, categoryId).First(&userJoinInfo).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{

				"status":  "0",
				"message": "Error joining Quiz.",
			})
		} else {
			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error",
			})
		}

	}

	userJoinContestHistory.CategoryID = userJoinInfo.CategoryID
	userJoinContestHistory.JoinID = userJoinInfo.ID
	userJoinContestHistory.UserID = userJoinInfo.UserID
	userJoinContestHistory.JoinedAt = userJoinInfo.JoinedAt

	if err := qc.DB.Create(&userJoinContestHistory).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while creating join contest history.",
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

// =============================================== GetContestJoinedByUser start ========================================================

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

			c.JSON(204, gin.H{

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

		c.JSON(204, gin.H{

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

			"contest_name":              category.Title,
			"contest_date":              category.QuizTime,
			"contest_end_date":          category.QuizEndTime,
			"contest_question_duration": category.EachQuestionTimeDuration,
			"contest_amount":            category.TotalPrice,
			"contest_id":                category.ID,
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
	ContestName             string `json:"contest_name" example:"GK"`
	ContestDate             string `json:"contest_date" example:"2024-12-21T18:00:00+05:30"`
	ContestEndDate          string `json:"contest_end_date" example:"2024-12-21T18:00:00+05:40"`
	ContestQuestionDuration int    `json:"contest_question_duration" example:"15"`
	ContestAmount           int    `json:"contest_amount" example:"10000"`
	ContestID               int    `json:"contest_id" example:"1"`
}

// =============================================== GetRulesByCategory start ========================================================

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

		c.JSON(204, gin.H{

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

// =============================================== GetContestPrizes start ========================================================

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

		c.JSON(204, gin.H{

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

// =============================================== UserContestResponse start ========================================================

// UserContestResponse This API will record user response for questions
//
//	@Summary		This API will record user response for questions
//	@Description	This API will record user response for questions
//	@Schemes
//	@Tags		Quizes
//	@Accept		application/x-www-form-urlencoded
//	@Produce	json
//	@Param		user_id	formData		int	true	"user id"
//	@Param		category_id	formData		string	true	"quiz category id"
//	@Param		question_id	formData		string	true	"quiz question id"
//	@Param		answer_given	formData		string	true	"quiz answer_given that user has given like a, b, c, d"
//	@Param		answer_type	formData		string	true	"if user answer is corect or wrong pass CORRECT/WRONG"
//	@Param		time_taken	formData		string	true	"time taken by user to solve the answer in seconds"
//	@Success	200	{object}	QAResponse
//	@Router		/quizes/user-question-answer [post]
func (qc *QuizController) UserContestResponse(c *gin.Context) {

	userId := c.PostForm("user_id")
	questionId := c.PostForm("question_id")
	categoryId := c.PostForm("category_id")
	answerGiven := c.PostForm("answer_given")
	answerType := c.PostForm("answer_type")
	timeTaken := c.PostForm("time_taken")

	var userContestResult []models.UserContestResults
	var category models.QuizCategory
	var question models.QuizQuestion
	var userJoinContest models.UserJoinContest
	var contestPointsChart models.ContestPointsChart

	if err := qc.DB.Where("id = ?", categoryId).First(&category).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(204, gin.H{

				"status":  "0",
				"message": "Invalid category id or category not found.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching category details.",
			})

		}
		return

	}

	if category.QuizTime.After(time.Now()) {

		c.JSON(422, gin.H{

			"status":  "0",
			"message": "Contest not started yet.",
		})
		return

	}

	if err := qc.DB.Where("id = ? AND category_id = ?", questionId, categoryId).First(&question).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(204, gin.H{

				"status":  "0",
				"message": "Invalid question id or question found with this category.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while getting question details.",
			})

		}

		return

	}

	if err := qc.DB.Where("category_id = ? AND user_id = ?", categoryId, userId).First(&userJoinContest).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(204, gin.H{

				"status":  "0",
				"message": "User not found in this contest.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching user join data.",
			})

		}
		return

	}

	if err := qc.DB.Where("category_id = ? AND question_id = ? AND user_id = ?", categoryId, questionId, userId).Find(&userContestResult).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while fetching contest result.",
		})
		return

	}

	if len(userContestResult) > 0 {

		c.JSON(422, gin.H{

			"status":  "0",
			"message": "Contest respose alredy given",
		})
		return

	}

	timeTakenConv, err := strconv.ParseUint(timeTaken, 10, 32)
	if err != nil {

		c.JSON(400, gin.H{

			"status":  "0",
			"message": "Invalid time passed",
		})
		return

	}

	if err := qc.DB.Where("time_from <= ? AND time_to >= ?", timeTaken, timeTaken).First(&contestPointsChart).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(204, gin.H{

				"status":  "0",
				"message": "Time details not found.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching time details.",
			})

		}
		return

	}

	var points uint

	if answerType == "CORRECT" {

		points = contestPointsChart.TotalCorrectAnswerPoint

	} else {

		points = contestPointsChart.TotalWrongAnswerPoint

	}

	var userContestResultData models.UserContestResults

	userContestResultData.AnswerGiven = answerGiven
	userContestResultData.AnswerType = answerType
	userContestResultData.CategoryID = question.CategoryID
	userContestResultData.CreatedAt = time.Now()
	userContestResultData.Points = points
	userContestResultData.QuestionID = question.ID
	userContestResultData.TimeTaken = uint(timeTakenConv)
	userIdUint, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "Unexpected error",
		})

		return
	}
	userContestResultData.UserID = uint(userIdUint)

	if err := qc.DB.Create(&userContestResultData).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB erorr while addind data.",
		})
		return

	}

	c.JSON(201, gin.H{

		"status":  "1",
		"message": "User response recorded successfully.",
	})

}

type QAResponse struct {
	Message string `json:"message" example:"User response recorded successfully."`
	Status  string `json:"status" example:"1"`
}

// =============================================== GetUserContestReport start ========================================================

// GetUserContestReport This API will provide user contest report
//
//	@Summary		This API will provide user contest report
//	@Description	This API will provide user contest report
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		user_id	path		string	true	"user_id id"
//	@Param		category_id	path		string	true	"category_id id"
//	@Success	200	{object}	GetUserContestReportResponse
//	@Router		/quizes/get-user-contest-result/{user_id}/{category_id} [get]
func (qc *QuizController) GetUserContestReport(c *gin.Context) {

	userId := c.Param("user_id")
	categoryId := c.Param("category_id")

	var category models.QuizCategory
	var userJoinContest models.UserJoinContest
	var userJoinContestHistory models.UserJoinContestHistory
	var question []models.QuizQuestion
	// var leaderboard models.UserContestLeaderboard

	if err := qc.DB.Where("id = ?", categoryId).First(&category).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(200, gin.H{

				"status":  "0",
				"message": "Invalid category id or category not found.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching category details.",
			})

		}
		return

	}

	if err := qc.DB.Where("category_id = ? AND user_id = ?", categoryId, userId).First(&userJoinContestHistory).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(200, gin.H{

				"status":  "0",
				"message": "User not found in this contest.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching user join data.",
			})

		}
		return

	}

	if err := qc.DB.Where("category_id = ? AND user_id = ?", categoryId, userId).First(&userJoinContest).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(200, gin.H{

				"status":  "0",
				"message": "User not found in this contest.",
			})

		} else {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "DB error while fetching user join data.",
			})

		}
		return

	}

	if err := qc.DB.Where("category_id = ?", categoryId).Find(&question).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while getting question details.",
		})

		return

	}

	if len(question) == 0 {

		c.JSON(204, gin.H{

			"status":  "0",
			"message": "No questions found.",
		})
		return
	}

	var response []gin.H
	var points uint
	for _, ques := range question {

		var userAnswer string
		var userAnswerType string
		var ques_points uint

		var userContestResult models.UserContestResults
		if err := qc.DB.Where("category_id = ? AND question_id = ? AND user_id = ?", categoryId, ques.ID, userId).Take(&userContestResult).Error; err != nil {

			points += 0
			userAnswer = "N/A"
			userAnswerType = "NA"
			ques_points = 0

		} else {

			points += userContestResult.Points
			userAnswer = userContestResult.AnswerGiven
			userAnswerType = userContestResult.AnswerType
			ques_points = userContestResult.Points

		}

		response = append(response, gin.H{

			"question":         ques.Question,
			"answer_a":         ques.OptionA,
			"answer_b":         ques.OptionB,
			"answer_c":         ques.OptionC,
			"answer_d":         ques.OptionD,
			"correct_answer":   ques.CorrectAnswer,
			"user_answer":      userAnswer,
			"user_answer_type": userAnswerType,
			"points":           ques_points,
			"time_taken":       userContestResult.TimeTaken,
		})

	}

	// if err := qc.DB.Where("category_id = ? AND user_id = ?", categoryId, userId).First(&leaderboard).Error; err != nil {

	// 	if errors.Is(err, gorm.ErrRecordNotFound) {

	// 		// leaderboard.CategoryID = category.ID
	// 		// userIdUint, err := strconv.ParseUint(userId, 10, 32)
	// 		// if err != nil {

	// 		// 	c.JSON(500, gin.H{

	// 		// 		"status":  "0",
	// 		// 		"message": "Unexpected error",
	// 		// 	})

	// 		// 	return
	// 		// }
	// 		// leaderboard.UserID = uint(userIdUint)
	// 		// leaderboard.Points = points
	// 		// leaderboard.PrizeAmount = 0

	// 		// if createErr := qc.DB.Create(&leaderboard).Error; createErr != nil {
	// 		// 	c.JSON(500, gin.H{
	// 		// 		"status":  "0",
	// 		// 		"message": "DB error while creating leaderboard",
	// 		// 	})
	// 		// 	return
	// 		// }

	// 	} else {

	// 		c.JSON(500, gin.H{

	// 			"status":  "0",
	// 			"message": "DB error while creating leadeboard",
	// 		})
	// 		return

	// 	}

	// }

	// if err := qc.DB.Delete(&userJoinContest, userJoinContest.ID).Error; err != nil {

	// 	c.JSON(500, gin.H{

	// 		"status":  "0",
	// 		"message": "DB error while deletting join contest data.",
	// 	})
	// 	return

	// }

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "User result found",
		"details": gin.H{
			"total_points": points,
			"quiz_result":  response,
		},
	})

}

type GetUserContestReportResponse struct {
	Status  string  `json:"status" example:"1"`                  // @Description Status of the response
	Message string  `json:"message" example:"User result found"` // @Description Message in the response
	Details Details `json:"details"`                             // @Description Details of the quiz result
}

type Details struct {
	QuizResult  []QuizResult `json:"quiz_result"`  // @Description List of quiz results for the user
	TotalPoints int          `json:"total_points"` // @Description Total points scored by the user
}

type QuizResult struct {
	Question       string `json:"question" example:"Where is Delhi?"`          // @Description The question asked in the quiz
	AnswerA        string `json:"answer_a" example:"United States of America"` // @Description Option A of the question
	AnswerB        string `json:"answer_b" example:"England"`                  // @Description Option B of the question
	AnswerC        string `json:"answer_c" example:"India"`                    // @Description Option C of the question
	AnswerD        string `json:"answer_d" example:"Sri Lanka"`                // @Description Option D of the question
	CorrectAnswer  string `json:"correct_answer" example:"c"`                  // @Description The correct answer for the question
	UserAnswer     string `json:"user_answer" example:"c"`                     // @Description The answer given by the user
	UserAnswerType string `json:"user_answer_type" example:"CORRECT"`          // @Description Whether the answer is "CORRECT" or "WRONG"
	Points         int    `json:"points" example:"80"`                         // @Description Points scored for this particular question
	TimeTaken      int    `json:"time_taken" example:"1"`                      // @Description Points scored for this particular question
}

// =============================================== GetUserContestLeaderboard start ========================================================

// GetUserContestLeaderboard This API will provide user contest report
//
//	@Summary		This API will provide user contest report
//	@Description	This API will provide user contest report
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		category_id	path		string	true	"category_id id"
//	@Success	200	{object}	GetUserContestLeaderboardResponse
//	@Router		/quizes/get-contest-leaderboard/{category_id} [get]
func (qc *QuizController) GetUserContestLeaderboard(c *gin.Context) {

	categoryId := c.Param("category_id")
	var leaderboard []models.UserContestLeaderboard

	if err := qc.DB.Where("category_id = ?", categoryId).Order("points DESC").Find(&leaderboard).Error; err != nil {

		c.JSON(500, gin.H{
			"status":  "0",
			"message": "DB error while fetching leaderboard.",
		})
		return

	}

	if len(leaderboard) == 0 {

		c.JSON(204, gin.H{

			"status":  "0",
			"message": "Leaderboard data not found",
		})
		return

	}

	var response []gin.H
	for _, data := range leaderboard {

		var user models.User
		if err := qc.DB.Where("id = ?", data.UserID).First(&user).Error; err != nil {

			c.JSON(500, gin.H{

				"status":  "0",
				"message": "unexpected error",
			})
			return

		}

		response = append(response, gin.H{

			"user_name":    user.Name,
			"user_image":   user.Image,
			"points":       data.Points,
			"prize_amount": data.PrizeAmount,
		})

	}

	c.JSON(200, gin.H{

		"status":  "1",
		"message": "Leaderboard data found successfully.",
		"details": response,
	})

}

type GetUserContestLeaderboardResponse struct {
	Status  string             `json:"status" example:"1"`
	Message string             `json:"message" example:"Leaderboard data found successfully."`
	Details []LeaderboardEntry `json:"details"`
}

// LeaderboardEntry is the structure for each entry in the leaderboard
type LeaderboardEntry struct {
	Points      uint   `json:"points" example:"175"`
	PrizeAmount uint   `json:"prize_amount" example:"0"`
	UserImage   string `json:"user_image" example:"image"`
	UserName    string `json:"user_name" example:"snagpal"`
}

// =============================================== GetUserPlayedContest start ========================================================

// GetUserPlayedContest This API will list contest history joined by user
//
//	@Summary		This API will list contest history joined by user
//	@Description	This API will list contest history joined by user
//	@Schemes
//	@Tags		Quizes
//	@Accept		json
//	@Produce	json
//	@Param		user_id	path		string	true	"user id"
//	@Success	200	{object}	GetUserPlayedContestResponse
//	@Router		/quizes/get-user-contest-history/{user_id} [get]
func (qc *QuizController) GetUserPlayedContest(c *gin.Context) {

	userId := c.Param("user_id")

	var user models.User
	var userJoinInfoHistory []models.UserJoinContestHistory

	if err := qc.DB.Where("id = ? AND register = 1", userId).First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {

			c.JSON(204, gin.H{

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

	if err := qc.DB.Where("user_id = ?", userId).Find(&userJoinInfoHistory).Error; err != nil {

		c.JSON(500, gin.H{

			"status":  "0",
			"message": "DB error while fetching user contest join info. " + err.Error(),
		})
		return

	}

	if len(userJoinInfoHistory) == 0 {

		c.JSON(200, gin.H{

			"status":  "0",
			"message": "Contest list found empty.",
		})
		return

	}

	var response []gin.H
	for _, contests := range userJoinInfoHistory {

		var category models.QuizCategory
		var leaderboard models.UserContestLeaderboard

		err := qc.DB.Where("id = ?", contests.CategoryID).First(&category).Error
		if err != nil {
			continue
		}

		errr := qc.DB.Where("user_id = ? AND category_id = ?", user.ID, contests.CategoryID).First(&leaderboard).Error
		if errr != nil {
			continue
		}

		response = append(response, gin.H{

			"contest_name":              category.Title,
			"contest_date":              category.QuizTime,
			"contest_end_date":          category.QuizEndTime,
			"contest_question_duration": category.EachQuestionTimeDuration,
			"contest_amount":            category.TotalPrice,
			"contest_id":                category.ID,
			"points":                    leaderboard.Points,
			"prize_amount":              leaderboard.PrizeAmount,
		})

	}

	c.JSON(200, gin.H{
		"status":  "1",
		"message": "Contest list found.",
		"details": response,
	})

}

type GetUserPlayedContestResponse struct {
	Status  string               `json:"status" example:"1"`
	Message string               `json:"message" example:"Contest list found."`
	Details []ContestHistoryInfo `json:"details"`
}

type ContestHistoryInfo struct {
	ContestName             string `json:"contest_name" example:"GK"`
	ContestDate             string `json:"contest_date" example:"2024-12-21T18:00:00+05:30"`
	ContestEndDate          string `json:"contest_end_date" example:"2024-12-21T18:00:00+05:40"`
	ContestQuestionDuration int    `json:"contest_question_duration" example:"15"`
	ContestAmount           int    `json:"contest_amount" example:"10000"`
	ContestID               int    `json:"contest_id" example:"1"`
	Points                  int    `json:"points" example:"175"`
	PrizeAmount             int    `json:"prize_amount" example:"175"`
}

// =============================================== CreateLeaderboard cronjob start ========================================================

func (qc *QuizController) CreateLeaderboard(c *gin.Context) {

	var categories []models.QuizCategory

	if err := qc.DB.Where("quiz_end_time < ? AND leader_board_created = 0", time.Now()).Find(&categories).Error; err != nil {

		println("DB error while getting quiz list.")
		return

	}

	categoriesLength := len(categories)
	currentCategoryLength := 0

	if categoriesLength > 0 {

		for _, cat := range categories {

			currentCategoryLength++

			var userJoinContest []models.UserJoinContest

			if err := qc.DB.Where("category_id = ?", cat.ID).Find(&userJoinContest).Error; err != nil {

				println("DB error while getting quiz list.")

			}

			if len(userJoinContest) > 0 {

				for _, ujc := range userJoinContest {

					var leaderboard models.UserContestLeaderboard

					totalPoints := 0

					// Query the database to calculate the sum of points
					if err := qc.DB.Table("user_contest_results").Where("category_id = ? AND user_id = ?", ujc.CategoryID, ujc.UserID).Select("SUM(points)").Scan(&totalPoints).Error; err != nil {

						println("DB error while getting quiz list.")

					}

					leaderboard.CategoryID = ujc.CategoryID
					leaderboard.Points = uint(totalPoints)
					leaderboard.PrizeAmount = 0
					leaderboard.UserID = ujc.UserID

					if createErr := qc.DB.Create(&leaderboard).Error; createErr != nil {

						println("DB error while creating leaderboard.")

					}

					if err := qc.DB.Delete(&userJoinContest, ujc.ID).Error; err != nil {

						println("DB error while deleting user join contest.")

					}

				}

			}

			if currentCategoryLength == categoriesLength {

				cat.LeaderBoardCreated = "1"

				if err := qc.DB.Model(&cat).Select("leader_board_created").Save(&cat).Error; err != nil {

					println("DB error while updating category info.")

				}
			}

		}

	}

}

func (qc *QuizController) CloseEntries(c *gin.Context) {

	var quizCategories []models.QuizCategory

	if err := qc.DB.Where("stop_entries_time < ? AND stop_entries = '0'", time.Now()).Find(&quizCategories).Error; err != nil {

		println("DB error while fetching category info.")
		return
	}

	if len(quizCategories) > 0 {

		for _, categories := range quizCategories {

			categories.StopEntries = "1"

			if err := qc.DB.Model(&categories).Update("StopEntries", categories.StopEntries).Error; err != nil {

				println("DB error while fetching category info.")

			}

		}

	}

}
