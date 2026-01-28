package dto

import (
	"cutbray/pppk-json/internal/repositories/models"
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
