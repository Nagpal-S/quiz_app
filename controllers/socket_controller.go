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
		// Allow connections from any origin
		return true
	},
}

// HandleWebSocket sends JSON data to a WebSocket client (mobile app)
func (sc *SocketController) HandleWebSocket(c *gin.Context) {

	var categories []models.QuizCategory

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Respond with error if WebSocket upgrade fails
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to upgrade to WebSocket",
		})
		return
	}
	defer conn.Close()

	var response []gin.H

	// Continuously send structured JSON data to the client
	for {

		if err := sc.DB.Where("quiz_time <= ? AND quiz_end_time >= ?", time.Now(), time.Now()).Find(&categories).Error; err != nil {

			errorData := gin.H{
				"status":  "0",
				"message": "DB error while fetching categories data",
				"details": err.Error(),
			}
			conn.WriteJSON(errorData)
			break

		}

		var data gin.H

		if len(categories) == 0 {

			data = gin.H{
				"status":  "0",
				"message": "No Active quiz found.",
				"details": response,
			}

		} else {

			data = gin.H{
				"status":  "1",
				"message": "Quiz list",
				"details": gin.H{},
			}

			for _, value := range categories {

				var questions models.QuizQuestion

				if err := sc.DB.Where("`category_id` = ? AND `from` <= ? AND `to` >= ?", value.ID, time.Now(), time.Now()).First(&questions).Error; err != nil {

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

		err := conn.WriteJSON(data)
		if err != nil {
			// Handle error during WebSocket communication
			errorData := gin.H{
				"status":  "0",
				"message": "Failed to send data",
				"details": err.Error(),
			}
			conn.WriteJSON(errorData)
			break
		}
		// Wait for 3 seconds before sending the next update
		time.Sleep(15 * time.Second)
	}
}
