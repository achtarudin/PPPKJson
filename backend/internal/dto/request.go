package dto

// SubmitAnswerRequest represents the request payload for submitting an answer
type SubmitAnswerRequest struct {
	ExamQuestionID   uint `json:"exam_question_id" binding:"required" example:"1"`
	QuestionOptionID uint `json:"question_option_id" binding:"required" example:"59"`
}

// UpdateScoreRequest represents the request payload for updating question option score
type UpdateScoreRequest struct {
	Score int `json:"score" binding:"required,min=0,max=10" example:"5"`
}
