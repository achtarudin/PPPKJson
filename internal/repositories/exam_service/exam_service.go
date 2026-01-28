package exam_service

import (
	"context"
	"cutbray/pppk-json/internal/repositories/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ExamService struct {
	db *gorm.DB
}

func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{db: db}
}

// CreateExamSession creates a new exam session for a user and assigns random questions
func (s *ExamService) CreateExamSession(ctx context.Context, userID string) (*models.ExamSession, error) {
	sessionCode := fmt.Sprintf("EXAM_%s_%d", userID, time.Now().Unix())

	examSession := &models.ExamSession{
		UserID:      userID,
		SessionCode: sessionCode,
		Status:      "NOT_STARTED",
		ExpiresAt:   time.Now().Add(2 * time.Hour), // 2 hours from now
		Duration:    120,                           // 120 minutes
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create exam session
		if err := tx.Create(examSession).Error; err != nil {
			return fmt.Errorf("failed to create exam session: %w", err)
		}

		// Assign random questions for each category
		categories := []string{"MANAJERIAL", "SOSIAL KULTURAL", "TEKNIS", "WAWANCARA"}
		questionsPerCategory := 5
		orderNumber := 1

		for _, category := range categories {
			// Get random 5 questions from this category
			var questions []models.Question
			if err := tx.Where("category = ?", category).
				Order("RANDOM()").
				Limit(questionsPerCategory).
				Find(&questions).Error; err != nil {
				return fmt.Errorf("failed to get random questions for category %s: %w", category, err)
			}

			if len(questions) < questionsPerCategory {
				return fmt.Errorf("not enough questions in category %s: need %d, got %d",
					category, questionsPerCategory, len(questions))
			}

			// Assign questions to exam session
			for _, question := range questions {
				examQuestion := models.ExamQuestion{
					ExamSessionID: examSession.ID,
					QuestionID:    question.ID,
					Category:      category,
					OrderNumber:   orderNumber,
				}

				if err := tx.Create(&examQuestion).Error; err != nil {
					return fmt.Errorf("failed to assign question to exam: %w", err)
				}

				orderNumber++
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return examSession, nil
}

// GetExamSession retrieves an exam session with assigned questions
func (s *ExamService) GetExamSession(ctx context.Context, userID string) (*models.ExamSession, error) {
	var examSession models.ExamSession

	err := s.db.WithContext(ctx).
		Preload("ExamQuestions", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("ExamQuestions.Question").
		Preload("ExamQuestions.Question.Options").
		Where("user_id = ? AND status IN (?)", userID, []string{"NOT_STARTED", "IN_PROGRESS"}).
		Order("created_at DESC").
		First(&examSession).Error

	if err != nil {
		return nil, fmt.Errorf("exam session not found for user %s: %w", userID, err)
	}

	return &examSession, nil
}

// StartExam starts the exam session
func (s *ExamService) StartExam(ctx context.Context, sessionID uint) error {
	now := time.Now()
	return s.db.WithContext(ctx).
		Model(&models.ExamSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"status":     "IN_PROGRESS",
			"started_at": now,
		}).Error
}

// SubmitAnswer submits an answer for a question
func (s *ExamService) SubmitAnswer(ctx context.Context, examSessionID, examQuestionID, questionOptionID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the question option to get the score
		var option models.QuestionOption
		if err := tx.First(&option, questionOptionID).Error; err != nil {
			return fmt.Errorf("question option not found: %w", err)
		}

		// Get the exam question
		var examQuestion models.ExamQuestion
		if err := tx.First(&examQuestion, examQuestionID).Error; err != nil {
			return fmt.Errorf("exam question not found: %w", err)
		}

		// Check if answer already exists
		var existingAnswer models.UserAnswer
		err := tx.Where("exam_session_id = ? AND exam_question_id = ?",
			examSessionID, examQuestionID).First(&existingAnswer).Error

		if err == gorm.ErrRecordNotFound {
			// Create new answer
			userAnswer := models.UserAnswer{
				ExamSessionID:    examSessionID,
				ExamQuestionID:   examQuestionID,
				QuestionID:       examQuestion.QuestionID,
				QuestionOptionID: questionOptionID,
				Score:            option.Score,
				AnsweredAt:       time.Now(),
			}

			if err := tx.Create(&userAnswer).Error; err != nil {
				return fmt.Errorf("failed to create user answer: %w", err)
			}
		} else if err == nil {
			// Update existing answer
			existingAnswer.QuestionOptionID = questionOptionID
			existingAnswer.Score = option.Score
			existingAnswer.AnsweredAt = time.Now()

			if err := tx.Save(&existingAnswer).Error; err != nil {
				return fmt.Errorf("failed to update user answer: %w", err)
			}
		} else {
			return fmt.Errorf("error checking existing answer: %w", err)
		}

		return nil
	})
}

// CompleteExam completes the exam and calculates results
func (s *ExamService) CompleteExam(ctx context.Context, examSessionID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Update exam session status
		if err := tx.Model(&models.ExamSession{}).
			Where("id = ?", examSessionID).
			Updates(map[string]interface{}{
				"status":       "COMPLETED",
				"completed_at": now,
			}).Error; err != nil {
			return fmt.Errorf("failed to update exam session: %w", err)
		}

		// Calculate results per category
		categories := []string{"MANAJERIAL", "SOSIAL KULTURAL", "TEKNIS", "WAWANCARA"}
		totalScore := 0
		totalAnswered := 0

		for _, category := range categories {
			var categoryStats struct {
				TotalAnswered int
				TotalScore    int
			}

			err := tx.Table("user_answers ua").
				Select("COUNT(ua.id) as total_answered, COALESCE(SUM(ua.score), 0) as total_score").
				Joins("JOIN exam_questions eq ON eq.id = ua.exam_question_id").
				Where("ua.exam_session_id = ? AND eq.category = ?", examSessionID, category).
				Scan(&categoryStats).Error

			if err != nil {
				return fmt.Errorf("failed to calculate stats for category %s: %w", category, err)
			}

			percentage := float64(categoryStats.TotalScore) / 20.0 * 100.0 // Max score per category is 20
			grade := calculateGrade(percentage)
			isPassed := percentage >= 60.0

			examResult := models.ExamResult{
				ExamSessionID:  examSessionID,
				Category:       category,
				TotalQuestions: 5,
				TotalAnswered:  categoryStats.TotalAnswered,
				TotalScore:     categoryStats.TotalScore,
				MaxScore:       20,
				Percentage:     percentage,
				Grade:          grade,
				IsPassed:       isPassed,
			}

			if err := tx.Create(&examResult).Error; err != nil {
				return fmt.Errorf("failed to create exam result for category %s: %w", category, err)
			}

			totalScore += categoryStats.TotalScore
			totalAnswered += categoryStats.TotalAnswered
		}

		// Create overall exam summary
		overallPercentage := float64(totalScore) / 80.0 * 100.0 // Max total score is 80
		overallGrade := calculateGrade(overallPercentage)
		overallPassed := overallPercentage >= 60.0

		// Get exam session for user ID
		var examSession models.ExamSession
		if err := tx.First(&examSession, examSessionID).Error; err != nil {
			return fmt.Errorf("failed to get exam session: %w", err)
		}

		examSummary := models.ExamSummary{
			ExamSessionID:     examSessionID,
			UserID:            examSession.UserID,
			TotalQuestions:    20,
			TotalAnswered:     totalAnswered,
			TotalScore:        totalScore,
			MaxScore:          80,
			OverallPercentage: overallPercentage,
			OverallGrade:      overallGrade,
			IsPassed:          overallPassed,
			CompletedAt:       now,
		}

		if err := tx.Create(&examSummary).Error; err != nil {
			return fmt.Errorf("failed to create exam summary: %w", err)
		}

		return nil
	})
}

// calculateGrade calculates grade based on percentage
func calculateGrade(percentage float64) string {
	if percentage >= 90 {
		return "A"
	} else if percentage >= 80 {
		return "B"
	} else if percentage >= 70 {
		return "C"
	} else if percentage >= 60 {
		return "D"
	} else {
		return "E"
	}
}

// GetExamResults gets the exam results for a user
func (s *ExamService) GetExamResults(ctx context.Context, userID string) (*models.ExamSummary, []models.ExamResult, error) {
	var examSummary models.ExamSummary
	var examResults []models.ExamResult

	// Get exam summary
	err := s.db.WithContext(ctx).
		Preload("ExamSession").
		Where("user_id = ?", userID).
		Order("completed_at DESC").
		First(&examSummary).Error

	if err != nil {
		return nil, nil, fmt.Errorf("exam summary not found for user %s: %w", userID, err)
	}

	// Get detailed results per category
	err = s.db.WithContext(ctx).
		Where("exam_session_id = ?", examSummary.ExamSessionID).
		Order("category ASC").
		Find(&examResults).Error

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get exam results: %w", err)
	}

	return &examSummary, examResults, nil
}
