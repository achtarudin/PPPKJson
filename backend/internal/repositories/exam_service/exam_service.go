package exam_service

import (
	"context"
	"cutbray/pppk-json/internal/dto"
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
		ExpiresAt:   time.Now().Add(130 * time.Minute), // 130 minutes from now
		Duration:    130,                               // 130 minutes
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create exam session
		if err := tx.Create(examSession).Error; err != nil {
			return fmt.Errorf("failed to create exam session: %w", err)
		}

		// Assign random questions for each category with specific counts in order
		categoryOrder := []string{"TEKNIS", "MANAJERIAL", "SOSIAL KULTURAL", "WAWANCARA"}
		categoryQuestions := map[string]int{
			"TEKNIS":          90,
			"MANAJERIAL":      25,
			"SOSIAL KULTURAL": 20,
			"WAWANCARA":       10,
		}
		orderNumber := 1

		for _, category := range categoryOrder {
			questionCount := categoryQuestions[category]
			// Get random questions from this category
			var questions []models.Question
			if err := tx.Where("category = ?", category).
				Order("RANDOM()").
				Limit(questionCount).
				Find(&questions).Error; err != nil {
				return fmt.Errorf("failed to get random questions for category %s: %w", category, err)
			}

			if len(questions) < questionCount {
				return fmt.Errorf("not enough questions in category %s: need %d, got %d",
					category, questionCount, len(questions))
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

	// First check and update any expired sessions
	s.CheckAndUpdateExpiredSessions(ctx)

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
		// First, validate exam session status and expiry
		var examSession models.ExamSession
		if err := tx.First(&examSession, examSessionID).Error; err != nil {
			return fmt.Errorf("exam session not found: %w", err)
		}

		// Check if exam session is still in progress
		if examSession.Status == "COMPLETED" {
			return fmt.Errorf("exam has already been completed, cannot submit more answers")
		}

		if examSession.Status == "EXPIRED" {
			return fmt.Errorf("exam has expired, cannot submit answers")
		}

		if examSession.Status == "NOT_STARTED" {
			return fmt.Errorf("exam has not been started yet")
		}

		// Check if exam has expired
		if time.Now().After(examSession.ExpiresAt) {
			// Update exam session to expired
			tx.Model(&examSession).Update("status", "EXPIRED")
			return fmt.Errorf("exam has expired, cannot submit answers")
		}

		// Get the question option to get the score
		var option models.QuestionOption
		if err := tx.First(&option, questionOptionID).Error; err != nil {
			return fmt.Errorf("question option not found: %w", err)
		}

		// Get the exam question and validate it belongs to this exam session
		var examQuestion models.ExamQuestion
		if err := tx.Where("id = ? AND exam_session_id = ?", examQuestionID, examSessionID).
			First(&examQuestion).Error; err != nil {
			return fmt.Errorf("exam question not found or doesn't belong to this exam session: %w", err)
		}

		// Check if answer already exists
		var existingAnswer models.UserAnswer
		err := tx.Where("exam_session_id = ? AND exam_question_id = ?",
			examSessionID, examQuestionID).First(&existingAnswer).Error

		switch err {
		case gorm.ErrRecordNotFound:
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
		case nil:
			// Update existing answer
			existingAnswer.QuestionOptionID = questionOptionID
			existingAnswer.Score = option.Score
			existingAnswer.AnsweredAt = time.Now()

			if err := tx.Save(&existingAnswer).Error; err != nil {
				return fmt.Errorf("failed to update user answer: %w", err)
			}
		default:
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

		// Calculate results per category in order
		categoryOrder := []string{"TEKNIS", "MANAJERIAL", "SOSIAL KULTURAL", "WAWANCARA"}
		categoryMaxScores := map[string]int{
			"TEKNIS":          450, // 90 questions → 450 max (official PPPK spec)
			"MANAJERIAL":      100, // 25 questions → 100 max (from shared 200 total)
			"SOSIAL KULTURAL": 100, // 20 questions → 100 max (from shared 200 total)
			"WAWANCARA":       40,  // 10 questions → 40 max (official PPPK spec)
		}

		// Category passing thresholds - only Grade A (100%) and B (90%) pass
		categoryThresholds := map[string]float64{
			"TEKNIS":          90.0, // Minimum 90% (Grade B)
			"MANAJERIAL":      90.0, // Minimum 90% (Grade B)
			"SOSIAL KULTURAL": 90.0, // Minimum 90% (Grade B)
			"WAWANCARA":       90.0, // Minimum 90% (Grade B)
		}

		categoryQuestionCounts := map[string]int{
			"TEKNIS":          90,
			"MANAJERIAL":      25,
			"SOSIAL KULTURAL": 20,
			"WAWANCARA":       10,
		}
		totalScore := 0
		totalAnswered := 0

		for _, category := range categoryOrder {
			questionCount := categoryQuestionCounts[category]
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

			maxScore := categoryMaxScores[category]
			threshold := categoryThresholds[category]
			percentage := float64(categoryStats.TotalScore) / float64(maxScore) * 100.0
			grade := calculateGrade(percentage)
			// Pass if percentage meets minimum threshold
			isPassed := percentage >= threshold

			examResult := models.ExamResult{
				ExamSessionID:  examSessionID,
				Category:       category,
				TotalQuestions: questionCount,
				TotalAnswered:  categoryStats.TotalAnswered,
				TotalScore:     categoryStats.TotalScore,
				MaxScore:       maxScore,
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
		totalMaxScore := 690  // Official PPPK total: 450 + 100 + 100 + 40
		totalQuestions := 145 // 90 + 25 + 20 + 10
		overallPercentage := float64(totalScore) / float64(totalMaxScore) * 100.0
		overallGrade := calculateGrade(overallPercentage)
		// Overall passing: minimum 90% overall (Grade B or better)
		overallPassed := overallPercentage >= 90.0 // Only Grade A and B pass

		// Get exam session for user ID
		var examSession models.ExamSession
		if err := tx.First(&examSession, examSessionID).Error; err != nil {
			return fmt.Errorf("failed to get exam session: %w", err)
		}

		examSummary := models.ExamSummary{
			ExamSessionID:     examSessionID,
			UserID:            examSession.UserID,
			TotalQuestions:    totalQuestions,
			TotalAnswered:     totalAnswered,
			TotalScore:        totalScore,
			MaxScore:          totalMaxScore,
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
// 100%=A, 90%=B, 80%=C, 70%=D (minimum passing), <70%=E (fail)
func calculateGrade(percentage float64) string {
	if percentage >= 100 {
		return "A"
	} else if percentage >= 90 {
		return "B"
	} else if percentage >= 80 {
		return "C"
	} else if percentage >= 70 {
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

// CheckAndUpdateExpiredSessions updates expired exam sessions
func (s *ExamService) CheckAndUpdateExpiredSessions(ctx context.Context) error {
	now := time.Now()
	return s.db.WithContext(ctx).
		Model(&models.ExamSession{}).
		Where("expires_at < ? AND status IN (?)", now, []string{"NOT_STARTED", "IN_PROGRESS"}).
		Update("status", "EXPIRED").Error
}

// GetUserDashboard gets dashboard data including exam status and results
func (s *ExamService) GetUserDashboard(ctx context.Context, userID string) (*dto.DashboardData, error) {
	dashboard := &dto.DashboardData{
		UserID: userID,
	}

	// First check and update any expired sessions
	s.CheckAndUpdateExpiredSessions(ctx)

	// Get the latest exam session
	var examSession models.ExamSession
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&examSession).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			dashboard.HasExam = false
			dashboard.ExamStatus = "NO_EXAM"
			return dashboard, nil
		}
		return nil, fmt.Errorf("failed to get exam session: %w", err)
	}

	dashboard.HasExam = true
	dashboard.ExamStatus = examSession.Status
	dashboard.ExamSession = &examSession

	// If exam is completed, get the results
	if examSession.Status == "COMPLETED" {
		summary, results, err := s.GetExamResults(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get exam results: %w", err)
		}

		dashboard.ExamSummary = summary
		dashboard.ExamResults = results
	}

	// If exam is in progress, get progress info
	if examSession.Status == "IN_PROGRESS" || examSession.Status == "NOT_STARTED" {
		var answeredCount int64
		s.db.WithContext(ctx).
			Table("user_answers").
			Where("exam_session_id = ?", examSession.ID).
			Count(&answeredCount)

		dashboard.ProgressInfo = &dto.ProgressInfo{
			TotalQuestions:    145, // 90 + 25 + 20 + 10 total questions
			AnsweredQuestions: int(answeredCount),
			RemainingTime:     int(time.Until(examSession.ExpiresAt).Minutes()),
		}
	}

	return dashboard, nil
}

// GetAllUsersExamStatus gets exam status for all users who have taken exams
func (s *ExamService) GetAllUsersExamStatus(ctx context.Context) ([]dto.UserDashboardSummary, error) {
	// First check and update any expired sessions
	s.CheckAndUpdateExpiredSessions(ctx)

	var results []dto.UserDashboardSummary

	// Query to get all exam sessions with their summaries (if completed)
	query := `
		SELECT 
			es.user_id,
			es.status as exam_status,
			es.session_code,
			es.started_at,
			es.completed_at,
			esm.total_score,
			esm.max_score,
			esm.overall_percentage,
			esm.overall_grade,
			esm.is_passed
		FROM exam_sessions es
		LEFT JOIN exam_summaries esm ON es.id = esm.exam_session_id
		WHERE es.id IN (
			SELECT MAX(id) 
			FROM exam_sessions 
			GROUP BY user_id
		)
		ORDER BY es.created_at DESC
	`

	rows, err := s.db.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to get users exam status: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userSummary dto.UserDashboardSummary
		var startedAt, completedAt *time.Time

		err := rows.Scan(
			&userSummary.UserID,
			&userSummary.ExamStatus,
			&userSummary.SessionCode,
			&startedAt,
			&completedAt,
			&userSummary.TotalScore,
			&userSummary.MaxScore,
			&userSummary.Percentage,
			&userSummary.Grade,
			&userSummary.IsPassed,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user summary: %w", err)
		}

		// Convert time pointers to string pointers
		if startedAt != nil {
			startedAtStr := startedAt.Format(time.RFC3339)
			userSummary.StartedAt = &startedAtStr
		}
		if completedAt != nil {
			completedAtStr := completedAt.Format(time.RFC3339)
			userSummary.CompletedAt = &completedAtStr
		}

		results = append(results, userSummary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// GetUserAnswers gets existing answers for a user's active exam session
func (s *ExamService) GetUserAnswers(ctx context.Context, userID string) (map[uint]uint, error) {
	// Get the latest exam session
	var examSession models.ExamSession
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND status IN (?)", userID, []string{"NOT_STARTED", "IN_PROGRESS"}).
		Order("created_at DESC").
		First(&examSession).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return make(map[uint]uint), nil
		}
		return nil, fmt.Errorf("failed to get exam session: %w", err)
	}

	// Get existing answers
	var userAnswers []models.UserAnswer
	err = s.db.WithContext(ctx).
		Where("exam_session_id = ?", examSession.ID).
		Find(&userAnswers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user answers: %w", err)
	}

	// Convert to map[examQuestionID]optionID
	answers := make(map[uint]uint)
	for _, answer := range userAnswers {
		answers[answer.ExamQuestionID] = answer.QuestionOptionID
	}

	return answers, nil
}
