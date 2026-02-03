package dto

import (
	"cutbray/pppk-json/internal/repositories/models"
	"strconv"
)

// ToExamSessionResponse converts domain model to DTO
func ToExamSessionResponse(examSession *models.ExamSession) ExamSessionResponse {
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

// ToExamSummaryResponse converts domain model to DTO
func ToExamSummaryResponse(summary *models.ExamSummary) ExamSummaryResponse {
	return ExamSummaryResponse{
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
}

// ToExamResultResponses converts domain models to DTOs
func ToExamResultResponses(results []models.ExamResult) []ExamResultResponse {
	responses := make([]ExamResultResponse, len(results))
	for i, result := range results {
		responses[i] = ExamResultResponse{
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
	return responses
}

// ToExamResultsResponse combines summary and results into single response
func ToExamResultsResponse(summary *models.ExamSummary, results []models.ExamResult) ExamResultsResponse {
	return ExamResultsResponse{
		Summary:           ToExamSummaryResponse(summary),
		ResultsByCategory: ToExamResultResponses(results),
	}
}

// ToDashboardResponse converts dashboard data to response DTO
func ToDashboardResponse(dashboard *DashboardData) DashboardResponse {
	response := DashboardResponse{
		UserID:     dashboard.UserID,
		HasExam:    dashboard.HasExam,
		ExamStatus: dashboard.ExamStatus,
	}

	// Convert ExamSession if exists
	if dashboard.ExamSession != nil {
		if examSession, ok := dashboard.ExamSession.(*models.ExamSession); ok {
			sessionResponse := ToExamSessionResponse(examSession)
			response.ExamSession = &sessionResponse
		}
	}

	// Convert ExamResults if completed
	if dashboard.ExamStatus == "COMPLETED" && dashboard.ExamSummary != nil && dashboard.ExamResults != nil {
		if summary, ok := dashboard.ExamSummary.(*models.ExamSummary); ok {
			if results, ok := dashboard.ExamResults.([]models.ExamResult); ok {
				resultsResponse := ToExamResultsResponse(summary, results)
				response.ExamResults = &resultsResponse
			}
		}
	}

	// Convert ProgressInfo if in progress
	if dashboard.ProgressInfo != nil {
		response.ProgressInfo = &ProgressInfoResponse{
			TotalQuestions:    dashboard.ProgressInfo.TotalQuestions,
			AnsweredQuestions: dashboard.ProgressInfo.AnsweredQuestions,
			RemainingTime:     dashboard.ProgressInfo.RemainingTime,
		}
	}

	return response
}

// ToQuestionManagementResponse converts question model to management DTO
func ToQuestionManagementResponse(question *models.Question) QuestionManagementResponse {
	options := make([]QuestionOptionManagementResponse, len(question.Options))
	for i, opt := range question.Options {
		options[i] = QuestionOptionManagementResponse{
			ID:         opt.ID,
			QuestionID: opt.QuestionID,
			OptionText: opt.OptionText,
			Score:      opt.Score,
			CreatedAt:  opt.CreatedAt,
			UpdatedAt:  opt.UpdatedAt,
		}
	}

	return QuestionManagementResponse{
		ID:           question.ID,
		Category:     question.Category,
		QuestionText: question.QuestionText,
		Options:      options,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}
}

// ToQuestionManagementResponses converts slice of question models to management DTOs
func ToQuestionManagementResponses(questions []models.Question) []QuestionManagementResponse {
	responses := make([]QuestionManagementResponse, len(questions))
	for i, question := range questions {
		responses[i] = ToQuestionManagementResponse(&question)
	}
	return responses
}

// ToQuestionOptionManagementResponse converts question option model to management DTO
func ToQuestionOptionManagementResponse(option *models.QuestionOption) QuestionOptionManagementResponse {
	return QuestionOptionManagementResponse{
		ID:         option.ID,
		QuestionID: option.QuestionID,
		OptionText: option.OptionText,
		Score:      option.Score,
		CreatedAt:  option.CreatedAt,
		UpdatedAt:  option.UpdatedAt,
	}
}

// ToExportQuestionResponses converts slice of question models to export format
func ToExportQuestionResponses(questions []models.Question) []ExportQuestionResponse {
	responses := make([]ExportQuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = ToExportQuestionResponse(&question, i+1)
	}
	return responses
}

// ToExportQuestionResponse converts question model to export DTO format
func ToExportQuestionResponse(question *models.Question, index int) ExportQuestionResponse {
	options := make([]ExportQuestionOptionResponse, len(question.Options))
	for i, opt := range question.Options {
		options[i] = ExportQuestionOptionResponse{
			OptionText: opt.OptionText,
			Score:      opt.Score,
		}
	}

	return ExportQuestionResponse{
		ID:           strconv.Itoa(index), // Convert index to string (1,2,3...)
		Category:     question.Category,
		QuestionText: question.QuestionText,
		Options:      options,
	}
}
