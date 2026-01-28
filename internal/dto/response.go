package dto

import "time"

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Success message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error message"`
}

// ExamSessionResponse represents the exam session response
type ExamSessionResponse struct {
	SessionID     uint                    `json:"session_id" example:"1"`
	UserID        string                  `json:"user_id" example:"1234"`
	SessionCode   string                  `json:"session_code" example:"EXAM_1234_1643356800"`
	Status        string                  `json:"status" example:"NOT_STARTED" enums:"NOT_STARTED,IN_PROGRESS,COMPLETED,EXPIRED"`
	ExpiresAt     time.Time               `json:"expires_at" example:"2026-01-28T12:00:00Z"`
	Duration      int                     `json:"duration" example:"120"`
	Questions     []QuestionResponse      `json:"questions"`
	CategoryStats []CategoryStatsResponse `json:"category_stats"`
}

// QuestionResponse represents a question in the exam session
type QuestionResponse struct {
	ExamQuestionID uint                     `json:"exam_question_id" example:"1"`
	QuestionID     uint                     `json:"question_id" example:"15"`
	Category       string                   `json:"category" example:"MANAJERIAL" enums:"MANAJERIAL,SOSIAL KULTURAL,TEKNIS,WAWANCARA"`
	OrderNumber    int                      `json:"order_number" example:"1"`
	QuestionText   string                   `json:"question_text" example:"Atasan Anda melakukan rekayasa laporan..."`
	Options        []QuestionOptionResponse `json:"options"`
}

// QuestionOptionResponse represents a question option
type QuestionOptionResponse struct {
	ID         uint   `json:"id" example:"59"`
	OptionText string `json:"option_text" example:"Dalam hati tidak menyetujui hal tersebut"`
}

// CategoryStatsResponse represents statistics for a category
type CategoryStatsResponse struct {
	Category       string `json:"category" example:"MANAJERIAL"`
	TotalQuestions int    `json:"total_questions" example:"5"`
	AnsweredCount  int    `json:"answered_count" example:"0"`
}

// ExamResultsResponse represents the complete exam results
type ExamResultsResponse struct {
	Summary           ExamSummaryResponse  `json:"summary"`
	ResultsByCategory []ExamResultResponse `json:"results_by_category"`
}

// ExamSummaryResponse represents the exam summary
type ExamSummaryResponse struct {
	ID                uint      `json:"id" example:"1"`
	ExamSessionID     uint      `json:"exam_session_id" example:"1"`
	UserID            string    `json:"user_id" example:"1234"`
	TotalQuestions    int       `json:"total_questions" example:"20"`
	TotalAnswered     int       `json:"total_answered" example:"18"`
	TotalScore        int       `json:"total_score" example:"65"`
	MaxScore          int       `json:"max_score" example:"80"`
	OverallPercentage float64   `json:"overall_percentage" example:"81.25"`
	OverallGrade      string    `json:"overall_grade" example:"B"`
	IsPassed          bool      `json:"is_passed" example:"true"`
	CompletedAt       time.Time `json:"completed_at" example:"2026-01-28T11:30:00Z"`
}

// ExamResultResponse represents exam results by category
type ExamResultResponse struct {
	ID             uint    `json:"id" example:"1"`
	ExamSessionID  uint    `json:"exam_session_id" example:"1"`
	Category       string  `json:"category" example:"MANAJERIAL"`
	TotalQuestions int     `json:"total_questions" example:"5"`
	TotalAnswered  int     `json:"total_answered" example:"5"`
	TotalScore     int     `json:"total_score" example:"16"`
	MaxScore       int     `json:"max_score" example:"20"`
	Percentage     float64 `json:"percentage" example:"80.0"`
	Grade          string  `json:"grade" example:"B"`
	IsPassed       bool    `json:"is_passed" example:"true"`
}
