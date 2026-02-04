package handlers

import (
	"cutbray/pppk-json/internal/dto"
	"cutbray/pppk-json/internal/repositories/models"
	"cutbray/pppk-json/internal/repositories/question_service"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ginQuestionHandler struct {
	questionRepo question_service.QuestionService
}

func NewGinQuestionHandler(db *gorm.DB) *ginQuestionHandler {
	return &ginQuestionHandler{
		questionRepo: question_service.NewQuestionService(db),
	}
}

// RegisterRoutes registers all question management routes
func (h *ginQuestionHandler) RegisterRoutes(router *gin.Engine) {
	// Add the download endpoint at root level for easy access
	router.GET("/questions", h.DownloadQuestionsJSON)

	// Use the existing /api/v1 group from gin adapter
	v1 := router.Group("/api/v1")
	questionGroup := v1.Group("/questions")
	{
		questionGroup.GET("/management", h.GetQuestionsByCategory)
		questionGroup.PUT("/:questionID/option/:optionID/score", h.UpdateOptionScore)
		questionGroup.GET("/categories", h.GetCategories)
	}
}

// GetQuestionsByCategory retrieves questions filtered by category and search text with pagination
// @Summary Get questions by category and search text with pagination
// @Description Retrieves questions with their options, filtered by category and question text, with pagination support
// @Tags questions
// @Accept json
// @Produce json
// @Param category query string false "Category filter (TEKNIS, MANAJERIAL, SOSIAL KULTURAL, WAWANCARA)"
// @Param search query string false "Search by question text"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 10, use 0 for all)" minimum(0)
// @Success 200 {object} dto.APIResponse{data=dto.PaginatedQuestionResponse}
// @Router /questions/management [get]
func (h *ginQuestionHandler) GetQuestionsByCategory(c *gin.Context) {
	category := c.Query("category")
	searchText := c.Query("search")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse pagination parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	// Get total count using repository
	totalCount, err := h.questionRepo.CountQuestionsWithFilters(category, searchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to count questions",
			Error:   err.Error(),
		})
		return
	}

	// Get questions using repository
	var questions []models.Question
	if limit > 0 {
		offset := (page - 1) * limit
		questions, err = h.questionRepo.GetQuestionsWithFilters(category, searchText, offset, limit)
	} else {
		questions, err = h.questionRepo.GetAllQuestionsWithFilters(category, searchText)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to fetch questions",
			Error:   err.Error(),
		})
		return
	}

	// Convert to management DTOs
	questionResponses := dto.ToQuestionManagementResponses(questions)

	// Calculate pagination metadata
	totalPages := 1
	if limit > 0 {
		totalPages = int(math.Ceil(float64(totalCount) / float64(limit)))
	} else {
		page = 1 // Reset page to 1 when showing all
	}

	response := dto.PaginatedQuestionResponse{
		Questions: questionResponses,
		Pagination: dto.PaginationMetadata{
			CurrentPage:  page,
			ItemsPerPage: limit,
			TotalItems:   int(totalCount),
			TotalPages:   totalPages,
		},
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Questions retrieved successfully",
		Data:    response,
	})
}

// UpdateOptionScore updates the score of a specific question option
// @Summary Update question option score
// @Description Updates the score value for a specific question option
// @Tags questions
// @Accept json
// @Produce json
// @Param questionID path int true "Question ID"
// @Param optionID path int true "Option ID"
// @Param body body dto.UpdateScoreRequest true "New score value"
// @Success 200 {object} dto.APIResponse{data=dto.QuestionOptionManagementResponse}
// @Router /questions/{questionID}/option/{optionID}/score [put]
func (h *ginQuestionHandler) UpdateOptionScore(c *gin.Context) {
	questionIDStr := c.Param("questionID")
	optionIDStr := c.Param("optionID")

	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid question ID",
			Error:   err.Error(),
		})
		return
	}

	optionID, err := strconv.ParseUint(optionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid option ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate score range (0-10)
	finalScore := *req.Score
	if finalScore < 0 || finalScore > 10 {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Score must be between 0 and 10",
		})
		return
	}

	// Get the question option using repository
	option, err := h.questionRepo.GetQuestionOptionByID(uint(questionID), uint(optionID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "Question option not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to find question option",
			Error:   err.Error(),
		})
		return
	}

	// Update the score
	option.Score = finalScore
	if err := h.questionRepo.UpdateQuestionOption(option); err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to update score",
			Error:   err.Error(),
		})
		return
	}

	// Convert to management DTO
	optionResponse := dto.ToQuestionOptionManagementResponse(option)

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Score updated successfully",
		Data:    optionResponse,
	})
}

// GetCategories returns all available question categories
// @Summary Get question categories
// @Description Returns list of all available question categories
// @Tags questions
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]string}
// @Router /questions/categories [get]
func (h *ginQuestionHandler) GetCategories(c *gin.Context) {
	categories, err := h.questionRepo.GetDistinctCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// DownloadQuestionsJSON downloads questions as JSON file based on search parameters
// @Summary Download questions as JSON file
// @Description Downloads questions in JSON format based on category and search text filters
// @Tags questions
// @Accept json
// @Produce application/json
// @Param category query string false "Category filter (TEKNIS, MANAJERIAL, SOSIAL KULTURAL, WAWANCARA)"
// @Param search query string false "Search by question text"
// @Success 200 {array} dto.ExportQuestionResponse
// @Router /questions [get]
func (h *ginQuestionHandler) DownloadQuestionsJSON(c *gin.Context) {
	category := c.Query("category")
	searchText := c.Query("search")

	// Get questions using repository
	questions, err := h.questionRepo.GetAllQuestionsWithFilters(category, searchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to fetch questions",
			Error:   err.Error(),
		})
		return
	}

	// Convert to export format
	exportQuestions := dto.ToExportQuestionResponses(questions)

	// Return JSON without download headers (let frontend handle download)
	c.JSON(http.StatusOK, exportQuestions)
}
