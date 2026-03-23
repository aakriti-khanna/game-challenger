package models

// AnswerPayload represents the JSON payload exactly as the API server needs
type AnswerPayload struct {
	UserID string `json:"user_id"`
	Answer string `json:"answer"`
}
