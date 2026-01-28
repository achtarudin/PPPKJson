package handlers

import (
	"cutbray/pppk-json/internal/dto"
	"cutbray/pppk-json/internal/repositories/exam_service"
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
