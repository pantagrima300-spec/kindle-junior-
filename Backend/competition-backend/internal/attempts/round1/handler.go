package round1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"competition-backend/internal/auth"
	"competition-backend/internal/storage"
)

const QuestionDuration = 60 * time.Second

type Handler struct {
	Storage *storage.Storage
}

func NewHandler(store *storage.Storage) *Handler {
	return &Handler{
		Storage: store,
	}
}

type StartAttemptRequest struct {
	CompetitionID string `json:"competition_id"`
	Language      string `json:"language"`
}

func parsePocketBaseTime(value string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999Z",
		"2006-01-02 15:04:05Z",
		"2006-01-02 15:04:05.999999Z",
		"2006-01-02 15:04:05.999999999Z",
	}

	for _, format := range formats {
		t, err := time.Parse(format, value)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"unsupported PocketBase time format: %q",
		value,
	)
}

/*
---------------------------------------------------------
START / RESUME ATTEMPT
---------------------------------------------------------
*/

func (h *Handler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		http.Error(
			w,
			"authentication required",
			http.StatusUnauthorized,
		)
		return
	}

	var request StartAttemptRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if request.CompetitionID == "" {
		http.Error(
			w,
			"competition_id is required",
			http.StatusBadRequest,
		)
		return
	}

	competition, err := h.Storage.PocketBase.GetCompetition(
		r.Context(),
		request.CompetitionID,
	)

	if err != nil {
		http.Error(
			w,
			"competition not found",
			http.StatusNotFound,
		)
		return
	}

	if competition.Status != "scheduled" &&
		competition.Status != "auth" &&
		competition.Status != "active" {

		http.Error(
			w,
			"competition is not available",
			http.StatusForbidden,
		)
		return
	}

	/*
		-----------------------------------------------------
		STEP 1:
		CHECK IF USER ALREADY HAS AN ATTEMPT
		-----------------------------------------------------
	*/

	existingAttempt, err := h.Storage.PocketBase.GetRound1Attempt(
		r.Context(),
		token,
		user.ID,
		competition.ID,
	)

	if err != nil {
		http.Error(
			w,
			"failed to check existing attempt",
			http.StatusInternalServerError,
		)
		return
	}

	if existingAttempt != nil {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		json.NewEncoder(w).Encode(existingAttempt)
		return
	}

	// Validate language for new attempt
	if request.Language != "C" && request.Language != "Python" {
		http.Error(
			w,
			"invalid language",
			http.StatusBadRequest,
		)
		return
	}

	/*
		-----------------------------------------------------
		STEP 2:
		CREATE NEW ATTEMPT
		-----------------------------------------------------
	*/

	startedAt := time.Now().UTC()

	serverDeadline := startedAt.Add(
		time.Duration(
			competition.DurationSeconds,
		) * time.Second,
	)

	attempt, err := h.Storage.PocketBase.CreateRound1Attempt(
		r.Context(),
		token,
		user.ID,
		competition.ID,
		request.Language,
		startedAt.Format(time.RFC3339Nano),
		serverDeadline.Format(time.RFC3339Nano),
	)

	if err != nil {
		/*
			-------------------------------------------------
			IMPORTANT FIX:
			Handle race condition / duplicate request.
			-------------------------------------------------
		*/

		existingAttempt, getErr :=
			h.Storage.PocketBase.GetRound1Attempt(
				r.Context(),
				token,
				user.ID,
				competition.ID,
			)

		if getErr == nil && existingAttempt != nil {
			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(existingAttempt)

			return
		}

		http.Error(
			w,
			"failed to create attempt: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(attempt)
}

/*
---------------------------------------------------------
QUESTION RESPONSE
---------------------------------------------------------
*/

type Round1QuestionResponse struct {
	AttemptID        string   `json:"attempt_id"`
	QuestionNumber   int      `json:"question_number"`
	TotalQuestions   int      `json:"total_questions"`
	QuestionID       int      `json:"question_id"`
	Domain           string   `json:"domain"`
	Question         string   `json:"question"`
	Options          []string `json:"options"`
	ServerDeadline   string   `json:"server_deadline"`
	QuestionDeadline string   `json:"question_deadline"`
}

/*
---------------------------------------------------------
GET CURRENT QUESTION
---------------------------------------------------------
*/

func (h *Handler) GetQuestion(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		http.Error(
			w,
			"authentication required",
			http.StatusUnauthorized,
		)
		return
	}

	attemptID := r.URL.Query().Get("attempt_id")

	if attemptID == "" {
		http.Error(
			w,
			"attempt_id is required",
			http.StatusBadRequest,
		)
		return
	}

	attempt, err :=
		h.Storage.PocketBase.GetRound1AttemptByID(
			r.Context(),
			token,
			attemptID,
		)

	if err != nil {
		http.Error(
			w,
			"attempt not found",
			http.StatusNotFound,
		)
		return
	}

	if attempt.User != user.ID {
		http.Error(
			w,
			"forbidden",
			http.StatusForbidden,
		)
		return
	}

	if attempt.Status != "active" {
		http.Error(
			w,
			"attempt is no longer active",
			http.StatusForbidden,
		)
		return
	}

	serverDeadline, err :=
		parsePocketBaseTime(
			attempt.ServerDeadline,
		)

	if err != nil {
		http.Error(
			w,
			"invalid server deadline",
			http.StatusInternalServerError,
		)
		return
	}

	now := time.Now().UTC()

	/*
		-----------------------------------------------------
		OVERALL ATTEMPT DEADLINE
		-----------------------------------------------------
	*/

	if !now.Before(serverDeadline) {
		h.expireAttempt(
			r,
			attempt,
			serverDeadline,
			now,
		)

		http.Error(
			w,
			"attempt has expired",
			http.StatusForbidden,
		)

		return
	}

	/*
		-----------------------------------------------------
		CURRENT QUESTION
		-----------------------------------------------------
	*/

	questionNumber := int(attempt.CurrentQuestion)

	// current_question is stored as a 0-based question number.
	question := GetQuestion(attempt.Language, questionNumber+1)

	if question == nil {
		http.Error(w, "invalid question number", http.StatusInternalServerError)
		return
	}

	var questionDeadline time.Time

	if attempt.QuestionStartedAt != "" {

		questionStartedAt, err :=
			parsePocketBaseTime(
				attempt.QuestionStartedAt,
			)

		if err != nil {
			http.Error(
				w,
				"invalid question start time",
				http.StatusInternalServerError,
			)
			return
		}

		/*
			-------------------------------------------------
			EXACTLY 60 SECONDS PER QUESTION
			-------------------------------------------------
		*/

		questionDeadline =
			questionStartedAt.Add(
				QuestionDuration,
			)

		if questionDeadline.After(serverDeadline) {
			questionDeadline = serverDeadline
		}

		/*
			-------------------------------------------------
			QUESTION EXPIRED
			-------------------------------------------------
		*/

		if !now.Before(questionDeadline) {

			elapsed :=
				now.Sub(questionStartedAt).Seconds()

			if elapsed > QuestionDuration.Seconds() {
				elapsed =
					QuestionDuration.Seconds()
			}

			newTimeTaken :=
				attempt.TimeTakenSeconds + elapsed

			if questionNumber >= GetTotalQuestions(attempt.Language)-1 {

				_, updateErr :=
					h.Storage.PocketBase.UpdateRound1Attempt(
						r.Context(),
						attempt.ID,
						map[string]interface{}{
							"status": "submitted",

							"time_taken_seconds": newTimeTaken,

							"submitted_at": now.Format(
								time.RFC3339Nano,
							),
						},
					)

				if updateErr != nil {
					http.Error(
						w,
						"failed to finish attempt",
						http.StatusInternalServerError,
					)
					return
				}

				http.Error(
					w,
					"attempt completed",
					http.StatusForbidden,
				)

				return
			}

			/*
				-------------------------------------------------
				ADVANCE TO NEXT QUESTION
				-------------------------------------------------
			*/

			nextQuestionNumber :=
				questionNumber + 1

			nextQuestionStartedAt :=
				now.UTC()

			updatedAttempt, updateErr :=
				h.Storage.PocketBase.UpdateRound1Attempt(
					r.Context(),
					attempt.ID,
					map[string]interface{}{
						"current_question": nextQuestionNumber,

						"question_started_at": nextQuestionStartedAt.Format(
							time.RFC3339Nano,
						),

						"last_saved_at": nextQuestionStartedAt.Format(
							time.RFC3339Nano,
						),

						"time_taken_seconds": newTimeTaken,
					},
				)

			if updateErr != nil {
				http.Error(
					w,
					"failed to advance question",
					http.StatusInternalServerError,
				)
				return
			}

			attempt = updatedAttempt

			questionNumber =
				int(attempt.CurrentQuestion)

			// current_question is stored as a 0-based question number.
			question =
				GetQuestion(attempt.Language, questionNumber+1)

			if question == nil {
				http.Error(
					w,
					"invalid next question number",
					http.StatusInternalServerError,
				)
				return
			}

			questionStartedAt, err :=
				parsePocketBaseTime(
					attempt.QuestionStartedAt,
				)

			if err != nil {
				http.Error(
					w,
					"invalid next question start time",
					http.StatusInternalServerError,
				)
				return
			}

			questionDeadline =
				questionStartedAt.Add(
					QuestionDuration,
				)

			if questionDeadline.After(serverDeadline) {
				questionDeadline = serverDeadline
			}
		}
	}

	/*
		-----------------------------------------------------
		RESPONSE
		-----------------------------------------------------
	*/

	response := Round1QuestionResponse{
		AttemptID: attempt.ID,

		QuestionNumber: questionNumber,

		TotalQuestions: GetTotalQuestions(attempt.Language),

		QuestionID: question.ID,

		Domain: question.Domain,

		Question: question.Question,

		Options: question.Options,

		ServerDeadline: attempt.ServerDeadline,
	}

	if !questionDeadline.IsZero() {
		response.QuestionDeadline =
			questionDeadline.Format(
				time.RFC3339Nano,
			)
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

/*
---------------------------------------------------------
SUBMIT ANSWER
---------------------------------------------------------
*/

type SubmitAnswerRequest struct {
	AttemptID  string `json:"attempt_id"`
	QuestionID int    `json:"question_id"`
	Answer     int    `json:"answer"`
}

func (h *Handler) SubmitAnswer(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		http.Error(
			w,
			"authentication required",
			http.StatusUnauthorized,
		)
		return
	}

	var req SubmitAnswerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.AttemptID == "" {
		http.Error(
			w,
			"attempt_id is required",
			http.StatusBadRequest,
		)
		return
	}

	attempt, err :=
		h.Storage.PocketBase.GetRound1AttemptByID(
			r.Context(),
			token,
			req.AttemptID,
		)

	if err != nil {
		http.Error(
			w,
			"attempt not found",
			http.StatusNotFound,
		)
		return
	}

	if attempt.User != user.ID {
		http.Error(
			w,
			"forbidden",
			http.StatusForbidden,
		)
		return
	}

	if attempt.Status != "active" {
		http.Error(
			w,
			"attempt is no longer active",
			http.StatusForbidden,
		)
		return
	}

	serverDeadline, err :=
		parsePocketBaseTime(
			attempt.ServerDeadline,
		)

	if err != nil {
		http.Error(
			w,
			"invalid server deadline",
			http.StatusInternalServerError,
		)
		return
	}

	now := time.Now().UTC()

	if !now.Before(serverDeadline) {
		h.expireAttempt(
			r,
			attempt,
			serverDeadline,
			now,
		)

		http.Error(
			w,
			"attempt has expired",
			http.StatusForbidden,
		)

		return
	}

	questionNumber :=
		int(attempt.CurrentQuestion)

	// current_question is stored as a 0-based question number.
	question :=
		GetQuestion(attempt.Language, questionNumber+1)

	if question == nil {
		http.Error(
			w,
			"invalid current question",
			http.StatusInternalServerError,
		)
		return
	}

	if req.QuestionID != question.ID {
		http.Error(
			w,
			"answer does not match current question",
			http.StatusBadRequest,
		)
		return
	}

	var newTimeTaken float64 = attempt.TimeTakenSeconds

	if attempt.QuestionStartedAt != "" {

		questionStartedAt, err :=
			parsePocketBaseTime(
				attempt.QuestionStartedAt,
			)

		if err != nil {
			http.Error(
				w,
				"invalid question start time",
				http.StatusInternalServerError,
			)
			return
		}

		/*
			-------------------------------------------------
			EXACTLY 60 SECONDS
			-------------------------------------------------
		*/

		questionDeadline :=
			questionStartedAt.Add(
				QuestionDuration,
			)

		if questionDeadline.After(serverDeadline) {
			questionDeadline = serverDeadline
		}

		if !now.Before(questionDeadline) {
			http.Error(
				w,
				"question has expired",
				http.StatusForbidden,
			)
			return
		}

		elapsed :=
			now.Sub(questionStartedAt).Seconds()

		if elapsed < 0 {
			elapsed = 0
		}

		if elapsed > QuestionDuration.Seconds() {
			elapsed =
				QuestionDuration.Seconds()
		}

		newTimeTaken += elapsed
	}

	/*
		-----------------------------------------------------
		SAVE ANSWER
		-----------------------------------------------------
	*/

	if attempt.Answers == nil {
		attempt.Answers =
			make(map[string]interface{})
	}

	attempt.Answers[fmt.Sprintf("%d", question.ID)] = req.Answer

	/*
		-----------------------------------------------------
		RECALCULATE SCORE
		-----------------------------------------------------
	*/

	correctAnswersCount := 0
	score := 0.0

	for qIDStr, ansRaw := range attempt.Answers {

		var qID int

		fmt.Sscanf(
			qIDStr,
			"%d",
			&qID,
		)

		var ans int

		if ansFloat, ok :=
			ansRaw.(float64); ok {

			ans = int(ansFloat)

		} else if ansInt, ok :=
			ansRaw.(int); ok {

			ans = ansInt
		}

		var bank []Question
		if attempt.Language == "Python" {
			bank = PythonQuestionBank
		} else {
			bank = CQuestionBank
		}

		for _, q := range bank {

			if q.ID == qID {

				if ans == q.CorrectAnswer {
					correctAnswersCount++
					score += 1.0
				}

				break
			}
		}
	}

	/*
		-----------------------------------------------------
		UPDATE ATTEMPT
		-----------------------------------------------------
	*/

	updateFields :=
		map[string]interface{}{
			"answers": attempt.Answers,

			"correct_answers": correctAnswersCount,

			"score": score,

			"last_saved_at": now.Format(
				time.RFC3339Nano,
			),

			"time_taken_seconds": newTimeTaken,
		}

	if questionNumber >= GetTotalQuestions(attempt.Language)-1 {

		updateFields["status"] =
			"submitted"

		updateFields["submitted_at"] =
			now.Format(time.RFC3339Nano)

	} else {

		updateFields["current_question"] =
			questionNumber + 1

		updateFields["question_started_at"] =
			now.Format(time.RFC3339Nano)
	}

	_, updateErr :=
		h.Storage.PocketBase.UpdateRound1Attempt(
			r.Context(),
			attempt.ID,
			updateFields,
		)

	if updateErr != nil {
		http.Error(
			w,
			"failed to save answer",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	w.Write(
		[]byte(`{"status":"success"}`),
	)
}

/*
---------------------------------------------------------
EXPIRE ATTEMPT
---------------------------------------------------------
*/

func (h *Handler) expireAttempt(
	r *http.Request,
	attempt *storage.Round1AttemptRecord,
	serverDeadline time.Time,
	now time.Time,
) {
	startedAt, err :=
		parsePocketBaseTime(
			attempt.StartedAt,
		)

	var timeTaken float64 = attempt.TimeTakenSeconds

	if err == nil {

		totalDuration :=
			serverDeadline.
				Sub(startedAt).
				Seconds()

		if timeTaken > totalDuration {
			timeTaken = totalDuration
		}
	}

	h.Storage.PocketBase.UpdateRound1Attempt(
		r.Context(),
		attempt.ID,
		map[string]interface{}{
			"status": "expired",

			"time_taken_seconds": timeTaken,

			"submitted_at": now.Format(
				time.RFC3339Nano,
			),
		},
	)
}
