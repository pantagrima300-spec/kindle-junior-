package competition

type Competition struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	AuthStart       string `json:"auth_start"`
	AuthEnd         string `json:"auth_end"`
	DurationSeconds int    `json:"duration_seconds"`
	MaxAttempts     int    `json:"max_attempts"`
}
