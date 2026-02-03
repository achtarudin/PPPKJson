package question_service

import (
	"cutbray/pppk-json/internal/repositories/models"

	"gorm.io/gorm"
)

type QuestionService interface {
	GetQuestionsWithFilters(category, searchText string, offset, limit int) ([]models.Question, error)
	CountQuestionsWithFilters(category, searchText string) (int64, error)
	GetAllQuestionsWithFilters(category, searchText string) ([]models.Question, error)
	GetQuestionOptionByID(questionID, optionID uint) (*models.QuestionOption, error)
	UpdateQuestionOption(option *models.QuestionOption) error
	GetDistinctCategories() ([]string, error)
}

type questionService struct {
	db *gorm.DB
}

func NewQuestionService(db *gorm.DB) QuestionService {
	return &questionService{
		db: db,
	}
}

func (r *questionService) GetQuestionsWithFilters(category, searchText string, offset, limit int) ([]models.Question, error) {
	var questions []models.Question
	query := r.db.Preload("Options").Order("id ASC")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if searchText != "" {
		query = query.Where("LOWER(question_text) LIKE LOWER(?)", "%"+searchText+"%")
	}

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Find(&questions).Error
	return questions, err
}

func (r *questionService) CountQuestionsWithFilters(category, searchText string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Question{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if searchText != "" {
		query = query.Where("LOWER(question_text) LIKE LOWER(?)", "%"+searchText+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *questionService) GetAllQuestionsWithFilters(category, searchText string) ([]models.Question, error) {
	var questions []models.Question
	query := r.db.Preload("Options").Order("id ASC")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if searchText != "" {
		query = query.Where("LOWER(question_text) LIKE LOWER(?)", "%"+searchText+"%")
	}

	err := query.Find(&questions).Error
	return questions, err
}

func (r *questionService) GetQuestionOptionByID(questionID, optionID uint) (*models.QuestionOption, error) {
	var option models.QuestionOption
	err := r.db.Where("id = ? AND question_id = ?", optionID, questionID).First(&option).Error
	return &option, err
}

func (r *questionService) UpdateQuestionOption(option *models.QuestionOption) error {
	return r.db.Save(option).Error
}

func (r *questionService) GetDistinctCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Question{}).Distinct("category").Pluck("category", &categories).Error
	return categories, err
}
