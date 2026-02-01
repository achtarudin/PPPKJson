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

// DashboardResponse represents the dashboard data response
type DashboardResponse struct {
	UserID       string                `json:"user_id" example:"1234"`
	HasExam      bool                  `json:"has_exam" example:"true"`
	ExamStatus   string                `json:"exam_status" example:"IN_PROGRESS" enums:"NO_EXAM,NOT_STARTED,IN_PROGRESS,COMPLETED,EXPIRED"`
	ExamSession  *ExamSessionResponse  `json:"exam_session,omitempty"`
	ExamResults  *ExamResultsResponse  `json:"exam_results,omitempty"`
	ProgressInfo *ProgressInfoResponse `json:"progress_info,omitempty"`
}

// ProgressInfoResponse represents exam progress information
type ProgressInfoResponse struct {
	TotalQuestions    int `json:"total_questions" example:"20"`
	AnsweredQuestions int `json:"answered_questions" example:"8"`
	RemainingTime     int `json:"remaining_time_minutes" example:"95"`
}

// DashboardData represents dashboard information for a user (internal use)
type DashboardData struct {
	UserID       string        `json:"user_id"`
	HasExam      bool          `json:"has_exam"`
	ExamStatus   string        `json:"exam_status"` // NO_EXAM, NOT_STARTED, IN_PROGRESS, COMPLETED, EXPIRED
	ExamSession  interface{}   `json:"exam_session,omitempty"`
	ExamSummary  interface{}   `json:"exam_summary,omitempty"`
	ExamResults  interface{}   `json:"exam_results,omitempty"`
	ProgressInfo *ProgressInfo `json:"progress_info,omitempty"`
}

// ProgressInfo represents exam progress information (internal use)
type ProgressInfo struct {
	TotalQuestions    int `json:"total_questions"`
	AnsweredQuestions int `json:"answered_questions"`
	RemainingTime     int `json:"remaining_time_minutes"`
}

// UserListDashboardResponse represents the response for list of users dashboard
type UserListDashboardResponse struct {
	TotalUsers int                    `json:"total_users" example:"25"`
	Users      []UserDashboardSummary `json:"users"`
}

// UserDashboardSummary represents summary information for each user
type UserDashboardSummary struct {
	UserID      string   `json:"user_id" example:"1234"`
	ExamStatus  string   `json:"exam_status" example:"COMPLETED" enums:"NO_EXAM,NOT_STARTED,IN_PROGRESS,COMPLETED,EXPIRED"`
	SessionCode string   `json:"session_code,omitempty" example:"EXAM_1234_1643356800"`
	StartedAt   *string  `json:"started_at,omitempty" example:"2026-01-28T10:00:00Z"`
	CompletedAt *string  `json:"completed_at,omitempty" example:"2026-01-28T11:30:00Z"`
	TotalScore  *int     `json:"total_score,omitempty" example:"12"`
	MaxScore    *int     `json:"max_score,omitempty" example:"16"`
	Percentage  *float64 `json:"percentage,omitempty" example:"75.0"`
	Grade       *string  `json:"grade,omitempty" example:"B"`
	IsPassed    *bool    `json:"is_passed,omitempty" example:"true"`
}

// DetailedAnswer represents a detailed user answer with question and score information
type DetailedAnswer struct {
	ExamQuestionID   uint      `json:"exam_question_id" example:"1"`
	QuestionID       uint      `json:"question_id" example:"15"`
	QuestionText     string    `json:"question_text" example:"Atasan Anda melakukan rekayasa laporan..."`
	Category         string    `json:"category" example:"MANAJERIAL"`
	SelectedOptionID uint      `json:"selected_option_id" example:"59"`
	SelectedOption   string    `json:"selected_option" example:"Dalam hati tidak menyetujui hal tersebut"`
	Score            int       `json:"score" example:"3"`
	MaxScore         int       `json:"max_score" example:"4"`
	IsCorrect        bool      `json:"is_correct" example:"false"`
	CorrectOptionID  uint      `json:"correct_option_id" example:"60"`
	CorrectOption    string    `json:"correct_option" example:"Menolak dengan tegas dan melaporkan kepada atasan"`
	CorrectScore     int       `json:"correct_score" example:"4"`
	AnsweredAt       time.Time `json:"answered_at" example:"2026-01-28T11:15:00Z"`
}
