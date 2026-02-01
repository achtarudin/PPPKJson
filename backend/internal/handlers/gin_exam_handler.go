package handlers

import (
	"cutbray/pppk-json/internal/dto"
	"cutbray/pppk-json/internal/repositories/exam_service"
	"net/http"
	"strings"
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

// RegisterRoutes registers all exam-related routes
func (h *GinExamHandler) RegisterRoutes(router *gin.Engine) {
	examGroup := router.Group("/api/v1/exam")
	{
		examGroup.GET("/:userID", h.GetOrCreateExam)
		examGroup.POST("/:userID/start", h.StartExam)
		examGroup.POST("/:userID/answer", h.SubmitAnswer)
		examGroup.POST("/:userID/complete", h.CompleteExam)
		examGroup.GET("/:userID/results", h.GetExamResults)
		examGroup.GET("/:userID/dashboard", h.GetDashboard)
		examGroup.GET("/:userID/answers", h.GetUserAnswers)
		examGroup.GET("/:userID/detailed-answers", h.GetDetailedUserAnswers)
	}

	// Admin/Dashboard routes
	dashboardGroup := router.Group("/api/v1/dashboard")
	{
		dashboardGroup.GET("/users", h.GetAllUsersDashboard)
	}
}

// GetOrCreateExam creates or gets existing exam session
// @Summary Create or get exam session
// @Description Creates a new exam session with 20 random questions (5 per category) or returns existing active session
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} dto.APIResponse{data=dto.ExamSessionResponse} "Exam session ready"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 500 {object} dto.APIResponse "Failed to create exam session"
// @Router /exam/{userID} [get]
func (h *GinExamHandler) GetOrCreateExam(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
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
			c.JSON(http.StatusInternalServerError, dto.APIResponse{
				Success: false,
				Error:   "Failed to create exam session: " + err.Error(),
			})
			return
		}
	}

	// Convert to response format
	response := dto.ToExamSessionResponse(examSession)

	c.JSON(http.StatusOK, dto.APIResponse{
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
// @Success 200 {object} dto.APIResponse{data=map[string]interface{}} "Exam started successfully"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 404 {object} dto.APIResponse "Exam session not found"
// @Failure 500 {object} dto.APIResponse "Failed to start exam"
// @Router /exam/{userID}/start [post]
func (h *GinExamHandler) StartExam(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Error:   "Exam session not found",
		})
		return
	}

	// Start the exam
	if err := h.examService.StartExam(c.Request.Context(), examSession.ID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to start exam: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
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
// @Param request body dto.SubmitAnswerRequest true "Answer submission"
// @Success 200 {object} dto.APIResponse "Answer submitted successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request"
// @Failure 404 {object} dto.APIResponse "Exam session not found"
// @Failure 500 {object} dto.APIResponse "Failed to submit answer"
// @Router /exam/{userID}/answer [post]
func (h *GinExamHandler) SubmitAnswer(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var request dto.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Get exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Error:   "Exam session not found",
		})
		return
	}

	// Submit answer
	err = h.examService.SubmitAnswer(c.Request.Context(), examSession.ID, request.ExamQuestionID, request.QuestionOptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to submit answer: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
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
// @Success 200 {object} dto.APIResponse{data=map[string]interface{}} "Exam completed successfully"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 404 {object} dto.APIResponse "Exam session not found"
// @Failure 500 {object} dto.APIResponse "Failed to complete exam"
// @Router /exam/{userID}/complete [post]
func (h *GinExamHandler) CompleteExam(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	// Get exam session
	examSession, err := h.examService.GetExamSession(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Error:   "Exam session not found",
		})
		return
	}

	// Complete exam
	if err := h.examService.CompleteExam(c.Request.Context(), examSession.ID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to complete exam: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
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
// @Success 200 {object} dto.APIResponse{data=dto.ExamResultsResponse} "Exam results retrieved"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 404 {object} dto.APIResponse "Exam results not found"
// @Router /exam/{userID}/results [get]
func (h *GinExamHandler) GetExamResults(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	summary, results, err := h.examService.GetExamResults(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Error:   "Exam results not found",
		})
		return
	}

	// Convert to response format using mapper
	response := dto.ToExamResultsResponse(summary, results)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Exam results retrieved",
		Data:    response,
	})
}

// GetDashboard gets dashboard information for a user
// @Summary Get dashboard information
// @Description Gets comprehensive dashboard data including exam status, progress, and results if completed
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} dto.APIResponse{data=dto.DashboardResponse} "Dashboard data retrieved"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 500 {object} dto.APIResponse "Failed to get dashboard data"
// @Router /exam/{userID}/dashboard [get]
func (h *GinExamHandler) GetDashboard(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	dashboard, err := h.examService.GetUserDashboard(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to get dashboard data: " + err.Error(),
		})
		return
	}

	// Convert to response format
	response := dto.ToDashboardResponse(dashboard)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Dashboard data retrieved",
		Data:    response,
	})
}

// GetAllUsersDashboard gets dashboard information for all users
// @Summary Get all users dashboard
// @Description Gets dashboard data for all users who have taken exams, including their status and results
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.UserListDashboardResponse} "All users dashboard data retrieved"
// @Failure 500 {object} dto.APIResponse "Failed to get dashboard data"
// @Router /dashboard/users [get]
func (h *GinExamHandler) GetAllUsersDashboard(c *gin.Context) {
	users, err := h.examService.GetAllUsersExamStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to get users dashboard data: " + err.Error(),
		})
		return
	}

	response := dto.UserListDashboardResponse{
		TotalUsers: len(users),
		Users:      users,
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "All users dashboard data retrieved",
		Data:    response,
	})
}

// GetUserAnswers gets existing answers for user's active exam session
// @Summary Get user's existing answers
// @Description Gets existing answers for user's active exam session to repopulate on page reload
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} dto.APIResponse{data=map[string]interface{}} "User answers retrieved"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 500 {object} dto.APIResponse "Failed to get user answers"
// @Router /exam/{userID}/answers [get]
func (h *GinExamHandler) GetUserAnswers(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	answers, err := h.examService.GetUserAnswers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to get user answers: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "User answers retrieved",
		Data:    answers,
	})
}

// GetDetailedUserAnswers gets detailed answers with questions and scores for completed exam
// @Summary Get detailed user answers
// @Description Gets detailed answers with question text, selected options, and scores for completed exam
// @Tags exam
// @Accept json
// @Produce json
// @Param userID path string true "User ID" example("1234")
// @Success 200 {object} dto.APIResponse{data=map[string][]dto.DetailedAnswer} "Detailed answers retrieved"
// @Failure 400 {object} dto.APIResponse "Invalid user ID"
// @Failure 404 {object} dto.APIResponse "No completed exam found"
// @Failure 500 {object} dto.APIResponse "Failed to get detailed answers"
// @Router /exam/{userID}/detailed-answers [get]
func (h *GinExamHandler) GetDetailedUserAnswers(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	answers, err := h.examService.GetDetailedUserAnswers(c.Request.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no completed exam session found") {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Error:   "No completed exam found for user",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error:   "Failed to get detailed answers: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Detailed answers retrieved",
		Data:    answers,
	})
}
