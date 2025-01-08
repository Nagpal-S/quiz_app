package main

import (
	"log"
	"quizapp/models"
	"quizapp/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "quizapp/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Quiz App API's
//	@version		1.0
//	@description	This is list API's to be used in Quiz App.
//	@termsOfService	http://swagger.io/terms/

//	@host		localhost:8080
//	@BasePath	/

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/quizdb?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	models.MigrateUser(db)
	models.MigrateQuizCategory(db)
	models.MigrateQuizQuestion(db)
	models.MigrateUserWallet(db)
	models.MigrateUserTransactions(db)
	models.MigrateUserJoinContest(db)
	models.MigrateUserJoinContestHistory(db)
	models.MigrateContestRules(db)
	models.MigrateContestPrize(db)
	models.MigrateUserContestResult(db)
	models.MigrateContestPointsChart(db)
	models.MigrateUserContestLeaderboard(db)
	models.MigrateBanners(db)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Use(func(c *gin.Context) {
		// Setting the Content-Type header to application/json
		c.Header("Content-Type", "application/json")
		// Allowing all origins (CORS)
		c.Header("Access-Control-Allow-Origin", "*")
		// Allowing the methods GET, POST, OPTIONS, PUT, DELETE
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		// Allowing specific headers
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Continue to the next middleware/handler
		c.Next()
	})

	routes.RegisterUser(router, db)
	routes.QuizRoutes(router, db)
	routes.RegisterRoutes(router)
	routes.SettingsRoutes(router, db)
	routes.SetupRoutes(router, db)

	router.Run(":8080")

}
