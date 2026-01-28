package models

import (
	"time"

	"gorm.io/gorm"
)

// Question represents a test question in the PPPK exam
type Question struct {
	ID           uint             `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Category     string           `gorm:"column:category;type:varchar(100);not null" json:"category"`
	QuestionText string           `gorm:"column:question_text;type:text;not null" json:"question_text"`
	Options      []QuestionOption `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"options"`
	CreatedAt    time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"deleted_at"`
}

// TableName specifies the table name for Question model
func (Question) TableName() string {
	return "questions"
}

// QuestionOption represents an answer option for a question
type QuestionOption struct {
	ID         uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	QuestionID uint           `gorm:"column:question_id;not null;index" json:"question_id"`
	OptionText string         `gorm:"column:option_text;type:text;not null" json:"option_text"`
	Score      int            `gorm:"column:score;not null" json:"score"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName specifies the table name for QuestionOption model
func (QuestionOption) TableName() string {
	return "question_options"
}
