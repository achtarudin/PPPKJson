package handlers

import (
	"cutbray/pppk-json/internal/repositories/exam_service"
	"cutbray/pppk-json/internal/repositories/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GinExamHandler struct {
	examService *exam_service.ExamService
}

func NewGinExamHandler(db *gorm.DB) *GinExamHandler {
	return &GinExamHandler{
		examService: exam_service.NewExamService(db),
	}
}

// Response structures for Gin handlers
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Success message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:"Error message"`
}

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

type QuestionResponse struct {
	ExamQuestionID uint                     `json:"exam_question_id" example:"1"`
	QuestionID     uint                     `json:"question_id" example:"15"`
	Category       string                   `json:"category" example:"MANAJERIAL" enums:"MANAJERIAL,SOSIAL KULTURAL,TEKNIS,WAWANCARA"`
	OrderNumber    int                      `json:"order_number" example:"1"`
	QuestionText   string                   `json:"question_text" example:"Atasan Anda melakukan rekayasa laporan..."`
	Options        []QuestionOptionResponse `json:"options"`
}

type QuestionOptionResponse struct {
	ID         uint   `json:"id" example:"59"`
	OptionText string `json:"option_text" example:"Dalam hati tidak menyetujui hal tersebut"`
}

type CategoryStatsResponse struct {
	Category       string `json:"category" example:"MANAJERIAL"`
	TotalQuestions int    `json:"total_questions" example:"5"`
	AnsweredCount  int    `json:"answered_count" example:"0"`
}

type SubmitAnswerRequest struct {
	ExamQuestionID   uint `json:"exam_question_id" binding:"required" example:"1"`
	QuestionOptionID uint `json:"question_option_id" binding:"required" example:"59"`
}

type ExamResultsResponse struct {
	Summary           ExamSummaryResponse  `json:"summary"`
	ResultsByCategory []ExamResultResponse `json:"results_by_category"`
}

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

// RegisterRoutes registers all exam-related routes
func (h *GinExamHandler) RegisterRoutes(router *gin.Engine) {
	examGroup := router.Group("/api/v1/exam")
	{
		examGroup.GET("/:userID", h.GetOrCreateExam)
		examGroup.POST("/:userID/start", h.StartExam)
		examGroup.POST("/:userID/answer", h.SubmitAnswer)
		examGroup.POST("/:userID/complete", h.CompleteExam)
		examGroup.GET("/:userID/results", h.GetExamResults)
	}
}

// GetOrCreateExam creates or gets existing exam session
// @Summary Create or get exam session
// @Description Creates a new exam session with 20 random questions (5 per category) or returns existing active session
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} APIResponse{data=ExamSessionResponse} "Exam session ready"
// @Failure 400 {object} APIResponse "Invalid user ID"
// @Failure 500 {object} APIResponse "Failed to create exam session"
// @Router /exam/{userID} [get]
func (h *GinExamHandler) GetOrCreateExam(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Try to get existing active exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		// If no active session, create new one
		examSession, err = h.examService.CreateExamSession(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "Failed to create exam session: " + err.Error(),
			})
			return
		}
	}

	// Convert to response format
	response := h.convertToExamSessionResponse(examSession)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Exam session ready",
		Data:    response,
	})
}

// StartExam starts the exam session
// @Summary Start exam
// @Description Starts the exam timer and changes status to IN_PROGRESS
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} APIResponse{data=map[string]interface{}} "Exam started successfully"
// @Failure 400 {object} APIResponse "Invalid user ID"
// @Failure 404 {object} APIResponse "Exam session not found"
// @Failure 500 {object} APIResponse "Failed to start exam"
// @Router /exam/{userID}/start [post]
func (h *GinExamHandler) StartExam(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Exam session not found",
		})
		return
	}

	// Start the exam
	if err := h.examService.StartExam(c.Request.Context(), examSession.ID); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to start exam: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Exam started successfully",
		Data: gin.H{
			"session_id": examSession.ID,
			"status":     "IN_PROGRESS",
			"started_at": time.Now(),
		},
	})
}

// SubmitAnswer submits an answer for a question
// @Summary Submit answer
// @Description Submits user's answer for a specific question
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Param request body SubmitAnswerRequest true "Answer submission"
// @Success 200 {object} APIResponse "Answer submitted successfully"
// @Failure 400 {object} APIResponse "Invalid request"
// @Failure 404 {object} APIResponse "Exam session not found"
// @Failure 500 {object} APIResponse "Failed to submit answer"
// @Router /exam/{userID}/answer [post]
func (h *GinExamHandler) SubmitAnswer(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var request SubmitAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Get exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Exam session not found",
		})
		return
	}

	// Submit answer
	err = h.examService.SubmitAnswer(c.Request.Context(), examSession.ID, request.ExamQuestionID, request.QuestionOptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to submit answer: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Answer submitted successfully",
	})
}

// CompleteExam completes the exam and calculates results
// @Summary Complete exam
// @Description Completes the exam session and calculates final results
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} APIResponse{data=map[string]interface{}} "Exam completed successfully"
// @Failure 400 {object} APIResponse "Invalid user ID"
// @Failure 404 {object} APIResponse "Exam session not found"
// @Failure 500 {object} APIResponse "Failed to complete exam"
// @Router /exam/{userID}/complete [post]
func (h *GinExamHandler) CompleteExam(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Exam session not found",
		})
		return
	}

	// Complete exam
	if err := h.examService.CompleteExam(c.Request.Context(), examSession.ID); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to complete exam: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Exam completed successfully",
		Data: gin.H{
			"session_id":   examSession.ID,
			"status":       "COMPLETED",
			"completed_at": time.Now(),
		},
	})
}

// GetExamResults gets exam results for a user
// @Summary Get exam results
// @Description Retrieves detailed exam results including summary and category breakdown
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} APIResponse{data=ExamResultsResponse} "Exam results retrieved"
// @Failure 400 {object} APIResponse "Invalid user ID"
// @Failure 404 {object} APIResponse "Exam results not found"
// @Router /exam/{userID}/results [get]
func (h *GinExamHandler) GetExamResults(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	summary, results, err := h.examService.GetExamResults(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "Exam results not found",
		})
		return
	}

	// Convert to response format
	summaryResponse := ExamSummaryResponse{
		ID:                summary.ID,
		ExamSessionID:     summary.ExamSessionID,
		UserID:            summary.UserID,
		TotalQuestions:    summary.TotalQuestions,
		TotalAnswered:     summary.TotalAnswered,
		TotalScore:        summary.TotalScore,
		MaxScore:          summary.MaxScore,
		OverallPercentage: summary.OverallPercentage,
		OverallGrade:      summary.OverallGrade,
		IsPassed:          summary.IsPassed,
		CompletedAt:       summary.CompletedAt,
	}

	resultsResponse := make([]ExamResultResponse, len(results))
	for i, result := range results {
		resultsResponse[i] = ExamResultResponse{
			ID:             result.ID,
			ExamSessionID:  result.ExamSessionID,
			Category:       result.Category,
			TotalQuestions: result.TotalQuestions,
			TotalAnswered:  result.TotalAnswered,
			TotalScore:     result.TotalScore,
			MaxScore:       result.MaxScore,
			Percentage:     result.Percentage,
			Grade:          result.Grade,
			IsPassed:       result.IsPassed,
		}
	}

	response := ExamResultsResponse{
		Summary:           summaryResponse,
		ResultsByCategory: resultsResponse,
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Exam results retrieved",
		Data:    response,
	})
}

// Helper function to convert model to response
func (h *GinExamHandler) convertToExamSessionResponse(examSession *models.ExamSession) ExamSessionResponse {
	questions := make([]QuestionResponse, len(examSession.ExamQuestions))
	categoryStats := make(map[string]*CategoryStatsResponse)

	for i, eq := range examSession.ExamQuestions {
		options := make([]QuestionOptionResponse, len(eq.Question.Options))
		for j, opt := range eq.Question.Options {
			options[j] = QuestionOptionResponse{
				ID:         opt.ID,
				OptionText: opt.OptionText,
				// Note: We don't expose the score to prevent cheating
			}
		}

		questions[i] = QuestionResponse{
			ExamQuestionID: eq.ID,
			QuestionID:     eq.QuestionID,
			Category:       eq.Category,
			OrderNumber:    eq.OrderNumber,
			QuestionText:   eq.Question.QuestionText,
			Options:        options,
		}

		// Build category stats
		if categoryStats[eq.Category] == nil {
			categoryStats[eq.Category] = &CategoryStatsResponse{
				Category:       eq.Category,
				TotalQuestions: 0,
				AnsweredCount:  0, // TODO: Calculate from user_answers
			}
		}
		categoryStats[eq.Category].TotalQuestions++
	}

	// Convert map to slice
	categoryStatsSlice := make([]CategoryStatsResponse, 0, len(categoryStats))
	for _, stats := range categoryStats {
		categoryStatsSlice = append(categoryStatsSlice, *stats)
	}

	return ExamSessionResponse{
		SessionID:     examSession.ID,
		UserID:        examSession.UserID,
		SessionCode:   examSession.SessionCode,
		Status:        examSession.Status,
		ExpiresAt:     examSession.ExpiresAt,
		Duration:      examSession.Duration,
		Questions:     questions,
		CategoryStats: categoryStatsSlice,
	}
}
