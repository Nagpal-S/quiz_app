package models

import (
	"time"

	"gorm.io/gorm"
)

type QuizCategory struct {
	ID                       uint64    `gorm:"primaryKey;autoIncrement"`
	Active                   string    `json:"active" gorm:"type:varchar(1);default:'0'"`
	Title                    string    `json:"title" gorm:"type:varchar(44)"`
	TotalPrice               int       `json:"total_price" gorm:"type:int(11)"`
	Icon                     string    `json:"icon" gorm:"type:text"`
	NumOfUsersCanJoin        int       `json:"num_of_users_can_join" gorm:"type:int(11)"`
	NumOfUsersHaveJoined     int       `json:"num_of_users_have_joined" gorm:"type:int(11)"`
	QuizTime                 time.Time `json:"quiz_time" gorm:"type:datetime"`
	NumOfQuestions           int       `json:"num_of_questions" gorm:"type:int(11)"`
	EachQuestionTimeDuration int       `json:"each_question_time_duration" gorm:"type:int(11)"`
	QuizEndTime              time.Time `json:"quiz_end_time" gorm:"type:datetime"`
	JoinAmount               int       `json:"join_amount" gorm:"type:int(11)"`
	Created                  time.Time `json:"crated" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type QuizQuestion struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID    uint64    `gorm:"not null;index" json:"category_id"` // Foreign key
	Level         string    `gorm:"type:enum('easy','medium','hard');not null" json:"level"`
	Question      string    `gorm:"type:text;not null" json:"question"`
	OptionA       string    `gorm:"type:varchar(44);not null" json:"option_a"`
	OptionB       string    `gorm:"type:varchar(44);not null" json:"option_b"`
	OptionC       string    `gorm:"type:varchar(44);not null" json:"option_c"`
	OptionD       string    `gorm:"type:varchar(44);not null" json:"option_d"`
	CorrectAnswer string    `gorm:"type:enum('a','b','c','d');not null" json:"correct_answer"`
	CreatedAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ContestRules struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID uint64    `gorm:"not null;index" json:"category_id"` // Foreign key
	Rule       string    `gorm:"not null;index" json:"rule"`        // Foreign key
	CreatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type UserJoinContest struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID uint64    `gorm:"not null;index" json:"category_id"` // Foreign key
	UserID     uint      `gorm:"not null;index" json:"user_id"`     // Foreign key
	JoinedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"joined_at"`
}

type UserJoinContestHistory struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	JoinID     uint64    `gorm:"not null;" json:"join_id"`          // Foreign key
	CategoryID uint64    `gorm:"not null;index" json:"category_id"` // Foreign key
	UserID     uint      `gorm:"not null;index" json:"user_id"`     // Foreign key
	JoinedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"joined_at"`
}

type ContestPrize struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID uint64    `gorm:"not null;index" json:"category_id"` // Foreign key
	RankFrom   uint      `gorm:"not null;index" json:"rank_from"`   // Foreign key
	RankTo     uint      `gorm:"not null;index" json:"rank_to"`     // Foreign key
	Winning    float64   `gorm:"not null;index" json:"winning"`     // Foreign key
	CreatedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type UserContestResults struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  uint64    `gorm:"not null;index" json:"category_id"` // Foreign key
	QuestionID  uint64    `gorm:"not null;index" json:"question_id"` // Foreign key
	UserID      uint      `gorm:"not null;index" json:"user_id"`     // Foreign key
	AnswerGiven string    `gorm:"not null;index" json:"answer_given"`
	AnswerType  string    `gorm:"type:enum('CORRECT', 'WRONG');not null;check:answer_type IN ('CORRECT', 'WRONG')" json:"answer_type"`
	TimeTaken   uint      `gorm:"not null;index" json:"time_taken"`
	Points      uint      `gorm:"not null;index" json:"points"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ContestPointsChart struct {
	ID                      uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TimeFrom                uint      `gorm:"not null;index" json:"time_from"`
	TimeTo                  uint      `gorm:"not null;index" json:"time_to"`
	BasePoint               uint      `gorm:"not null;index" json:"base_point"`
	CorrectAnswerPoint      uint      `gorm:"not null;index" json:"correct_answer_point"`
	TotalCorrectAnswerPoint uint      `gorm:"not null;index" json:"total_correct_answer_point"`
	WrongAnswerPoint        uint      `gorm:"not null;index" json:"wrong_answer_point"`
	TotalWrongAnswerPoint   uint      `gorm:"not null;index" json:"total_wrong_answer_point"`
	CreatedAt               time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type UserContestLeaderboard struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  uint64 `gorm:"not null;index" json:"category_id"` // Foreign key
	UserID      uint   `gorm:"not null;index" json:"user_id"`     // Foreign key
	Points      uint   `gorm:"not null;index" json:"points"`
	PrizeAmount uint   `gorm:"not null;index" json:"prize_amount"`
}

// MigrateUser migrates the User table
func MigrateQuizCategory(db *gorm.DB) {
	db.AutoMigrate(&QuizCategory{})
}

func MigrateQuizQuestion(db *gorm.DB) {
	db.AutoMigrate(&QuizQuestion{})
}

func MigrateUserJoinContest(db *gorm.DB) {
	db.AutoMigrate(&UserJoinContest{})
}

func MigrateUserJoinContestHistory(db *gorm.DB) {
	db.AutoMigrate(&UserJoinContestHistory{})
}

func MigrateContestRules(db *gorm.DB) {
	db.AutoMigrate(&ContestRules{})
}

func MigrateContestPrize(db *gorm.DB) {
	db.AutoMigrate(&ContestPrize{})
}

func MigrateUserContestResult(db *gorm.DB) {
	db.AutoMigrate(&UserContestResults{})
}

func MigrateContestPointsChart(db *gorm.DB) {
	db.AutoMigrate(&ContestPointsChart{})
}

func MigrateUserContestLeaderboard(db *gorm.DB) {
	db.AutoMigrate(&UserContestLeaderboard{})
}
