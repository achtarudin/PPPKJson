package dto

// SubmitAnswerRequest represents the request payload for submitting an answer
type SubmitAnswerRequest struct {
	ExamQuestionID   uint `json:"exam_question_id" binding:"required" example:"1"`
	QuestionOptionID uint `json:"question_option_id" binding:"required" example:"59"`
}
