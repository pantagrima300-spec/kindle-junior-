package round2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"competition-backend/internal/auth"
	"competition-backend/internal/storage"
)

type Handler struct {
	Storage *storage.Storage
}

func NewHandler(store *storage.Storage) *Handler {
	return &Handler{
		Storage: store,
	}
}

type StartAttemptRequest struct {
	CompetitionID   string `json:"competition_id"`
	InitialLanguage string `json:"initial_language"`
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

func (h *Handler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	var request StartAttemptRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.CompetitionID == "" {
		jsonError(w, "competition_id is required", http.StatusBadRequest)
		return
	}

	if request.InitialLanguage == "" {
		request.InitialLanguage = string(LangC)
	}

	if request.InitialLanguage != string(LangC) &&
		request.InitialLanguage != string(LangPython) {
		jsonError(w, "invalid language", http.StatusBadRequest)
		return
	}

	competition, err := h.Storage.PocketBase.GetCompetition(
		r.Context(),
		request.CompetitionID,
	)

	if err != nil {
		jsonError(w, "competition not found", http.StatusNotFound)
		return
	}

	if competition.Status != "scheduled" &&
		competition.Status != "auth" &&
		competition.Status != "active" {
		jsonError(
			w,
			"competition is not available",
			http.StatusForbidden,
		)
		return
	}

	existingAttempt, err :=
		h.Storage.PocketBase.GetRound2Attempt(
			r.Context(),
			token,
			user.ID,
			competition.ID,
		)

	if err != nil {
		jsonError(
			w,
			"failed to check existing attempt",
			http.StatusInternalServerError,
		)
		return
	}

	if existingAttempt != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingAttempt)
		return
	}

	// Round 2 starts when the attempt is created.
	startedAt := time.Now().UTC()

	// Round duration = 70 minutes.
	serverDeadline := startedAt.Add(4200 * time.Second)

	attempt, err :=
		h.Storage.PocketBase.CreateRound2Attempt(
			r.Context(),
			token,
			user.ID,
			competition.ID,
			startedAt.Format(time.RFC3339Nano),
			serverDeadline.Format(time.RFC3339Nano),
			request.InitialLanguage,
		)

	if err != nil {
		http.Error(
			w,
			"failed to create attempt: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(attempt)
}

type Round2StateResponse struct {
	AttemptID        string  `json:"attempt_id"`
	StartedAt        string  `json:"started_at"`
	ServerDeadline   string  `json:"server_deadline"`
	CurrentQuestion  int     `json:"current_question"`
	TotalQuestions   int     `json:"total_questions"`
	SelectedLanguage string  `json:"selected_language"`
	Score            float64 `json:"score"`
	ProblemsSolved   int     `json:"problems_solved"`
	Status           string  `json:"status"`
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	attemptID := r.URL.Query().Get("attempt_id")

	if attemptID == "" {
		jsonError(w, "attempt_id is required", http.StatusBadRequest)
		return
	}

	attempt, err :=
		h.Storage.PocketBase.GetRound2AttemptByID(
			r.Context(),
			token,
			attemptID,
		)

	if err != nil {
		jsonError(w, "attempt not found", http.StatusNotFound)
		return
	}

	if attempt.User != user.ID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}

	/*
		Use the stored server_deadline when available.

		This keeps the deadline tied to the original attempt
		and prevents changing the timer when questions change.
	*/
	serverDeadline, err :=
		parsePocketBaseTime(attempt.ServerDeadline)

	if err != nil || serverDeadline.IsZero() {
		// Backward-compatible fallback for older attempts.
		startedAt, startErr :=
			parsePocketBaseTime(attempt.StartedAt)

		if startErr != nil {
			jsonError(
				w,
				"invalid attempt timing data",
				http.StatusInternalServerError,
			)
			return
		}

		serverDeadline = startedAt.Add(4200 * time.Second)
	}

	if time.Now().UTC().After(serverDeadline) &&
		attempt.Status != "expired" &&
		attempt.Status != "submitted" {

		attempt.Status = "expired"
		attempt.TimeTakenSeconds = 4200

		updateFields := map[string]interface{}{
			"status":             "expired",
			"time_taken_seconds": 4200,
		}

		_, updateErr :=
			h.Storage.PocketBase.UpdateRound2Attempt(
				r.Context(),
				attempt.ID,
				updateFields,
			)

		if updateErr != nil {
			log.Printf(
				"GetState: failed to mark attempt %s expired: %v",
				attempt.ID,
				updateErr,
			)
		}
	}

	resp := Round2StateResponse{
		AttemptID:        attempt.ID,
		StartedAt:        attempt.StartedAt,
		ServerDeadline:   serverDeadline.Format(time.RFC3339Nano),
		CurrentQuestion:  int(attempt.CurrentProblem),
		TotalQuestions:   len(QuestionBank),
		SelectedLanguage: attempt.SelectedLanguage,
		Score:            attempt.Score,
		ProblemsSolved:   int(attempt.ProblemsSolved),
		Status:           attempt.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(resp)
}

type Round2QuestionResponse struct {
	AttemptID        string                  `json:"attempt_id"`
	StartedAt        string                  `json:"started_at"`
	QuestionNumber   int                     `json:"question_number"`
	TotalQuestions   int                     `json:"total_questions"`
	QuestionID       int                     `json:"question_id"`
	Type             QuestionType            `json:"type"`
	Title            string                  `json:"title"`
	Description      string                  `json:"description"`
	SelectedLanguage string                  `json:"selected_language"`
	LanguageContent  LanguageSpecificContent `json:"language_content"`
}

func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	attemptID := r.URL.Query().Get("attempt_id")

	if attemptID == "" {
		jsonError(w, "attempt_id is required", http.StatusBadRequest)
		return
	}

	attempt, err :=
		h.Storage.PocketBase.GetRound2AttemptByID(
			r.Context(),
			token,
			attemptID,
		)

	if err != nil {
		jsonError(w, "attempt not found", http.StatusNotFound)
		return
	}

	if attempt.User != user.ID {
		jsonError(w, "forbidden", http.StatusForbidden)
		return
	}

	if attempt.Status != "in_progress" {
		jsonError(
			w,
			"attempt is no longer active",
			http.StatusForbidden,
		)
		return
	}

	questionNumber := int(attempt.CurrentProblem)

	question := GetQuestion(questionNumber)

	if question == nil {
		jsonError(
			w,
			"invalid question number",
			http.StatusInternalServerError,
		)
		return
	}

	langContent, ok :=
		question.Languages[
			LanguageCode(attempt.SelectedLanguage),
		]

	if !ok {
		langContent = LanguageSpecificContent{}
	}

	response := Round2QuestionResponse{
		AttemptID:        attempt.ID,
		StartedAt:        attempt.StartedAt,
		QuestionNumber:   questionNumber,
		TotalQuestions:   len(QuestionBank),
		QuestionID:       question.ID,
		Type:             question.Type,
		Title:            question.Title,
		Description:      question.Description,
		SelectedLanguage: attempt.SelectedLanguage,
		LanguageContent:  langContent,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

type RunCodeRequest struct {
    AttemptID  string `json:"attempt_id"`
    QuestionID int    `json:"question_id"`
    Language   string `json:"language"`
    Code       string `json:"code"`
}

type SubmitAnswerRequest struct {
	AttemptID  string `json:"attempt_id"`
	QuestionID int    `json:"question_id"`
	Language   string `json:"language"`
	Code       string `json:"code"`
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(
		map[string]string{
			"error": msg,
		},
	)
}

func (h *Handler) RunCode(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	var raw []byte

	raw, _ = io.ReadAll(r.Body)

	log.Printf(
		"RunCode raw request: %s",
		string(raw),
	)

	r.Body = io.NopCloser(bytes.NewBuffer(raw))

	var req RunCodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	attempt, err :=
		h.Storage.PocketBase.GetRound2AttemptByID(
			r.Context(),
			token,
			req.AttemptID,
		)

	if err != nil {
		jsonError(
			w,
			"attempt not found",
			http.StatusNotFound,
		)
		return
	}

	if attempt.User != user.ID ||
		attempt.Status != "in_progress" {

		jsonError(
			w,
			"forbidden or inactive attempt",
			http.StatusForbidden,
		)
		return
	}

	/*
		Server-side timer check.
	*/
	serverDeadline, deadlineErr :=
		parsePocketBaseTime(attempt.ServerDeadline)

	if deadlineErr != nil || serverDeadline.IsZero() {
		startedAt, startErr :=
			parsePocketBaseTime(attempt.StartedAt)

		if startErr != nil {
			jsonError(
				w,
				"invalid attempt timing data",
				http.StatusInternalServerError,
			)
			return
		}

		serverDeadline = startedAt.Add(4200 * time.Second)
	}

	if time.Now().UTC().After(serverDeadline) {
		jsonError(w, "time up", http.StatusForbidden)
		return
	}

	/*
		The selected language is controlled by the attempt.
	*/
	req.Language = attempt.SelectedLanguage

	// if float64(req.QuestionID) != attempt.CurrentProblem {
	// 	jsonError(
	// 		w,
	// 		"invalid question id",
	// 		http.StatusBadRequest,
	// 	)
	// 	return
	// }

	// ✅ REPLACE WITH THIS:
	question := GetQuestion(req.QuestionID)

	if question == nil {
    	jsonError(
        	w,
        	"question not found",
        	http.StatusNotFound,
    	)
    return
}

	/*
		Run Code only checks visible test cases.
		It does NOT advance the question.
	*/
	res :=
		JudgeSubmission(
			r.Context(),
			question,
			LanguageCode(req.Language),
			req.Code,
			false,
		)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUser(r.Context())
	token := r.Header.Get("Authorization")

	if !ok || token == "" {
		jsonError(w, "authentication required", http.StatusUnauthorized)
		return
	}

	var req SubmitAnswerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	attempt, err :=
		h.Storage.PocketBase.GetRound2AttemptByID(
			r.Context(),
			token,
			req.AttemptID,
		)

	if err != nil {
		jsonError(
			w,
			"attempt not found",
			http.StatusNotFound,
		)
		return
	}

	if attempt.User != user.ID ||
		attempt.Status != "in_progress" {

		jsonError(
			w,
			"forbidden or inactive attempt",
			http.StatusForbidden,
		)
		return
	}

	/*
		Use the stored server deadline.
		This ensures the 70-minute timer starts only once,
		when the attempt is created.
	*/
	startedAtSub, startedAtErr :=
		parsePocketBaseTime(attempt.StartedAt)

	if startedAtErr != nil {
		jsonError(
			w,
			"invalid attempt timing data",
			http.StatusInternalServerError,
		)
		return
	}

	serverDeadlineSub, deadlineErr :=
		parsePocketBaseTime(attempt.ServerDeadline)

	if deadlineErr != nil || serverDeadlineSub.IsZero() {
		serverDeadlineSub =
			startedAtSub.Add(4200 * time.Second)
	}

	if time.Now().UTC().After(serverDeadlineSub) {
		updateFields := map[string]interface{}{
			"status":             "expired",
			"time_taken_seconds": 4200,
		}

		_, updateErr :=
			h.Storage.PocketBase.UpdateRound2Attempt(
				r.Context(),
				attempt.ID,
				updateFields,
			)

		if updateErr != nil {
			log.Printf(
				"SubmitAnswer: failed to mark attempt %s expired: %v",
				attempt.ID,
				updateErr,
			)
		}

		jsonError(w, "time up", http.StatusForbidden)
		return
	}

	/*
		Language is controlled by the attempt.
	*/
	req.Language = attempt.SelectedLanguage

	/*
		If the question has already been passed,
		don't process it again.
	*/
	if float64(req.QuestionID) < attempt.CurrentProblem {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"status": "already_submitted",
				"score":  attempt.Score,
			},
		)

		return
	}

	/*
		The participant must submit the current question.
	*/
	if float64(req.QuestionID) != attempt.CurrentProblem {
		jsonError(
			w,
			"invalid question id",
			http.StatusBadRequest,
		)
		return
	}

	question := GetQuestion(req.QuestionID)

	if question == nil {
		jsonError(
			w,
			"question not found",
			http.StatusNotFound,
		)
		return
	}

	/*
		Judge all test cases during actual submission.

		IMPORTANT:
		res.Success determines SCORE only.

		It no longer determines whether the participant
		can move to the next question.
	*/
	res :=
		JudgeSubmission(
			r.Context(),
			question,
			LanguageCode(req.Language),
			req.Code,
			true,
		)

	/*
		Advance the question regardless of whether
		the submitted code passed all test cases.

		Passed:
		    +10 score
		    +1 solved

		Failed:
		    +0 score
		    solved count unchanged
		    BUT question still advances.
	*/
	newTotalScore := attempt.Score
	newProblemsSolved := attempt.ProblemsSolved

	if res.Success {
		newTotalScore += 10
		newProblemsSolved++
	}

	newCurrentProblem := attempt.CurrentProblem + 1

	/*
		If this was the last question, mark the attempt
		as submitted. Otherwise keep it in progress.
	*/
	status := "in_progress"

	if int(newCurrentProblem) > len(QuestionBank) {
		status = "submitted"
	}

	now := time.Now().UTC()

	timeTaken := time.Since(startedAtSub).Seconds()

	if timeTaken < 0 {
		timeTaken = 0
	}

	if timeTaken > 4200 {
		timeTaken = 4200
	}

	updateFields := map[string]interface{}{
		"problems_solved": int(newProblemsSolved),
		"current_problem": int(newCurrentProblem),
		"score":           newTotalScore,
		"status":          status,
		"last_saved_at":   now.Format(time.RFC3339Nano),
	}

	/*
		Only the final submission stores completion time.
	*/
	if status == "submitted" {
		updateFields["time_taken_seconds"] = timeTaken
		updateFields["submitted_at"] = now.Format(time.RFC3339Nano)
	}

	log.Printf(
		"SubmitAnswer: attempt=%s question=%d success=%v score=%.0f solved=%.0f next_question=%d status=%s",
		attempt.ID,
		req.QuestionID,
		res.Success,
		newTotalScore,
		newProblemsSolved,
		int(newCurrentProblem),
		status,
	)

	_, err =
		h.Storage.PocketBase.UpdateRound2Attempt(
			r.Context(),
			attempt.ID,
			updateFields,
		)

	if err != nil {
		log.Printf(
			"PB Error updating attempt %s: %v",
			attempt.ID,
			err,
		)

		jsonError(
			w,
			"failed to update attempt: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	log.Printf(
		"SubmitAnswer: successfully updated attempt %s to status=%s",
		attempt.ID,
		status,
	)

	/*
		Return the judge result along with navigation information.
		This makes the frontend response easier to handle.
	*/
	response := map[string]interface{}{
		"success":         res.Success,
		"status":          res.Status,
		"language":        res.Language,
		"question_id":     res.QuestionID,
		"tests":           res.Tests,
		"passed_tests":    res.PassedTests,
		"total_tests":     res.TotalTests,
		"score":            newTotalScore,
		"problems_solved": int(newProblemsSolved),
		"current_problem": int(newCurrentProblem),
		"next_question":   int(newCurrentProblem),
		"round_status":    status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}