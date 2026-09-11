package round1

type Attempt struct {
	ID               string                 `json:"id"`
	User             string                 `json:"user"`
	Competition      string                 `json:"competition"`
	Language         string                 `json:"language"`
	StartedAt        string                 `json:"started_at"`
	ServerDeadline   string                 `json:"server_deadline"`
	LastSavedAt      string                 `json:"last_saved_at"`
	CurrentQuestion  float64                `json:"current_question"`
	Answers          map[string]interface{} `json:"answers"`
	Status           string                 `json:"status"`
	CorrectAnswers   float64                `json:"correct_answers"`
	Score            float64                `json:"score"`
	TimeTakenSeconds float64                `json:"time_taken_seconds"`
	SubmittedAt      string                 `json:"submitted_at"`
}
