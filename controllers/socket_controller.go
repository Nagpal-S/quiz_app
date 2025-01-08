package controllers

import (
	"fmt"
	"net/http"
	"quizapp/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type SocketController struct {
	DB *gorm.DB
}

// Upgrader to handle WebSocket upgrade
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// HandleWebSocket sends JSON data to a WebSocket client
func (sc *SocketController) HandleWebSocket(c *gin.Context) {
	var categories []models.QuizCategory

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Predefined intervals
	predefinedIntervals := []int{1, 16, 31, 46}

	// Helper function to get the nearest previous interval
	getPreviousInterval := func(current time.Time) time.Time {
		currentMinute := current.Truncate(time.Minute)
		currentSecond := current.Second()

		// Loop in reverse to find the most recent valid interval
		for i := len(predefinedIntervals) - 1; i >= 0; i-- {
			if currentSecond >= predefinedIntervals[i] {
				return currentMinute.Add(time.Duration(predefinedIntervals[i]) * time.Second)
			}
		}

		// If no valid interval in the current minute, move to the last interval of the previous minute
		previousMinute := currentMinute.Add(-time.Minute)
		return previousMinute.Add(time.Duration(predefinedIntervals[len(predefinedIntervals)-1]) * time.Second)
	}

	// Helper function to fetch and send quiz data for a specific interval
	sendQuizData := func(targetTime time.Time) {
		if err := sc.DB.Where("quiz_time <= ? AND quiz_end_time >= ?", targetTime, targetTime).Find(&categories).Error; err != nil {
			conn.WriteJSON(gin.H{"status": "0", "message": "DB error while fetching categories data", "details": err.Error()})
			return
		}

		var data gin.H
		if len(categories) == 0 {
			data = gin.H{"status": "0", "message": "No Active quiz found.", "details": []gin.H{}}
		} else {
			data = gin.H{"status": "1", "message": "Quiz list", "details": gin.H{}}
			for _, value := range categories {
				var questions models.QuizQuestion
				if err := sc.DB.Where("`category_id` = ? AND `from` <= ? AND `to` >= ?", value.ID, targetTime, targetTime).First(&questions).Error; err != nil {
					continue
				}

				questionsData := gin.H{
					"question_id":         questions.ID,
					"total_question":      value.NumOfQuestions,
					"question_number":     questions.QuestionNumber,
					"question":            questions.Question,
					"a":                   questions.OptionA,
					"b":                   questions.OptionB,
					"c":                   questions.OptionC,
					"d":                   questions.OptionD,
					"correct_answer":      questions.CorrectAnswer,
					"question_level":      questions.Level,
					"question_start_time": questions.From,
					"question_end_time":   questions.To,
				}

				data["details"].(gin.H)[fmt.Sprintf("%d", value.ID)] = questionsData
			}
		}

		// Send data to the client
		if err := conn.WriteJSON(data); err != nil {
			conn.WriteJSON(gin.H{"status": "0", "message": "Failed to send data", "details": err.Error()})
		}
	}

	// Fetch and send data for the previous interval immediately
	previousInterval := getPreviousInterval(time.Now())
	sendQuizData(previousInterval)

	// Sleep until the next predefined interval
	currentTime := time.Now()
	var nextStartTime time.Time
	for _, interval := range predefinedIntervals {
		potentialStart := currentTime.Truncate(time.Minute).Add(time.Duration(interval) * time.Second)
		if currentTime.Before(potentialStart) {
			nextStartTime = potentialStart
			break
		}
	}
	if nextStartTime.IsZero() {
		nextStartTime = currentTime.Truncate(time.Minute).Add(time.Minute).Add(time.Duration(predefinedIntervals[0]) * time.Second)
	}
	time.Sleep(time.Until(nextStartTime))

	// Continuously send structured JSON data to the client at each interval
	for {
		sendQuizData(time.Now())

		// Wait for the next interval (15 seconds gap in your current logic)
		time.Sleep(15 * time.Second)
	}
}
