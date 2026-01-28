package models

import (
	"time"

	"gorm.io/gorm"
)

// ExamSession represents an exam session for a user
type ExamSession struct {
	ID          uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      string         `gorm:"column:user_id;type:varchar(50);not null;index" json:"user_id"` // Hardcoded user ID from URL
	SessionCode string         `gorm:"column:session_code;type:varchar(100);uniqueIndex;not null" json:"session_code"`
	Status      string         `gorm:"column:status;type:varchar(20);default:'NOT_STARTED'" json:"status"` // NOT_STARTED, IN_PROGRESS, COMPLETED, EXPIRED
	StartedAt   *time.Time     `gorm:"column:started_at" json:"started_at"`
	CompletedAt *time.Time     `gorm:"column:completed_at" json:"completed_at"`
	ExpiresAt   time.Time      `gorm:"column:expires_at;not null" json:"expires_at"`
	Duration    int            `gorm:"column:duration;default:120" json:"duration"` // Duration in minutes (default 2 hours)
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	ExamQuestions []ExamQuestion `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"exam_questions,omitempty"`
	UserAnswers   []UserAnswer   `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"user_answers,omitempty"`
	ExamResults   []ExamResult   `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"exam_results,omitempty"`
}

// TableName specifies the table name for ExamSession model
func (ExamSession) TableName() string {
	return "exam_sessions"
}

// ExamQuestion represents the assigned questions for a specific exam session
// This table ensures each user gets different random questions per category
type ExamQuestion struct {
	ID            uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExamSessionID uint           `gorm:"column:exam_session_id;not null;index" json:"exam_session_id"`
	QuestionID    uint           `gorm:"column:question_id;not null;index" json:"question_id"`
	Category      string         `gorm:"column:category;type:varchar(50);not null;index" json:"category"` // MANAJERIAL, SOSIAL_KULTURAL, TEKNIS, WAWANCARA
	OrderNumber   int            `gorm:"column:order_number;not null" json:"order_number"`                // Question order in the exam (1-20)
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	ExamSession ExamSession `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"exam_session,omitempty"`
	Question    Question    `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"question,omitempty"`
}

// TableName specifies the table name for ExamQuestion model
func (ExamQuestion) TableName() string {
	return "exam_questions"
}

// UserAnswer represents a user's answer to a specific question in their exam
type UserAnswer struct {
	ID               uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExamSessionID    uint           `gorm:"column:exam_session_id;not null;index" json:"exam_session_id"`
	ExamQuestionID   uint           `gorm:"column:exam_question_id;not null;index" json:"exam_question_id"`
	QuestionID       uint           `gorm:"column:question_id;not null;index" json:"question_id"`
	QuestionOptionID uint           `gorm:"column:question_option_id;not null;index" json:"question_option_id"`
	Score            int            `gorm:"column:score;not null" json:"score"`
	AnsweredAt       time.Time      `gorm:"column:answered_at;not null" json:"answered_at"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	ExamSession    ExamSession    `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"exam_session,omitempty"`
	ExamQuestion   ExamQuestion   `gorm:"foreignKey:ExamQuestionID;constraint:OnDelete:CASCADE" json:"exam_question,omitempty"`
	Question       Question       `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"question,omitempty"`
	QuestionOption QuestionOption `gorm:"foreignKey:QuestionOptionID;constraint:OnDelete:CASCADE" json:"question_option,omitempty"`
}

// TableName specifies the table name for UserAnswer model
func (UserAnswer) TableName() string {
	return "user_answers"
}

// ExamResult represents the result of an exam session per category
type ExamResult struct {
	ID             uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExamSessionID  uint           `gorm:"column:exam_session_id;not null;index" json:"exam_session_id"`
	Category       string         `gorm:"column:category;type:varchar(50);not null" json:"category"` // MANAJERIAL, SOSIAL_KULTURAL, TEKNIS, WAWANCARA
	TotalQuestions int            `gorm:"column:total_questions;default:5" json:"total_questions"`   // Always 5 per category
	TotalAnswered  int            `gorm:"column:total_answered;not null" json:"total_answered"`
	TotalScore     int            `gorm:"column:total_score;not null" json:"total_score"`
	MaxScore       int            `gorm:"column:max_score;default:20" json:"max_score"` // 5 questions * 4 max score = 20
	Percentage     float64        `gorm:"column:percentage;not null" json:"percentage"`
	Grade          string         `gorm:"column:grade;type:varchar(5)" json:"grade"` // A, B, C, D, E
	IsPassed       bool           `gorm:"column:is_passed;default:false" json:"is_passed"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	ExamSession ExamSession `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"exam_session,omitempty"`
}

// TableName specifies the table name for ExamResult model
func (ExamResult) TableName() string {
	return "exam_results"
}

// ExamSummary represents the overall exam summary for a user
type ExamSummary struct {
	ID                uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExamSessionID     uint           `gorm:"column:exam_session_id;not null;uniqueIndex" json:"exam_session_id"`
	UserID            string         `gorm:"column:user_id;type:varchar(50);not null;index" json:"user_id"`
	TotalQuestions    int            `gorm:"column:total_questions;default:20" json:"total_questions"` // 4 categories * 5 questions = 20
	TotalAnswered     int            `gorm:"column:total_answered;not null" json:"total_answered"`
	TotalScore        int            `gorm:"column:total_score;not null" json:"total_score"`
	MaxScore          int            `gorm:"column:max_score;default:80" json:"max_score"` // 20 questions * 4 max score = 80
	OverallPercentage float64        `gorm:"column:overall_percentage;not null" json:"overall_percentage"`
	OverallGrade      string         `gorm:"column:overall_grade;type:varchar(5)" json:"overall_grade"`
	IsPassed          bool           `gorm:"column:is_passed;default:false" json:"is_passed"`
	CompletedAt       time.Time      `gorm:"column:completed_at;not null" json:"completed_at"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	ExamSession ExamSession `gorm:"foreignKey:ExamSessionID;constraint:OnDelete:CASCADE" json:"exam_session,omitempty"`
}

// TableName specifies the table name for ExamSummary model
func (ExamSummary) TableName() string {
	return "exam_summaries"
}
