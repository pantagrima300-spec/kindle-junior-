package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Round1AttemptRecord struct {
	ID                string                 `json:"id"`
	User              string                 `json:"user"`
	Competition       string                 `json:"competition"`
	Language          string                 `json:"language"`
	StartedAt         string                 `json:"started_at"`
	ServerDeadline    string                 `json:"server_deadline"`
	QuestionStartedAt string                 `json:"question_started_at"`
	LastSavedAt       string                 `json:"last_saved_at"`
	CurrentQuestion   float64                `json:"current_question"`
	Answers           map[string]interface{} `json:"answers"`
	Status            string                 `json:"status"`
	CorrectAnswers    float64                `json:"correct_answers"`
	Score             float64                `json:"score"`
	TimeTakenSeconds  float64                `json:"time_taken_seconds"`
	SubmittedAt       string                 `json:"submitted_at"`
}

type Round1AttemptListResponse struct {
	Page       int                   `json:"page"`
	PerPage    int                   `json:"perPage"`
	TotalItems int                   `json:"totalItems"`
	TotalPages int                   `json:"totalPages"`
	Items      []Round1AttemptRecord `json:"items"`
}

func (pb *PocketBaseClient) GetRound1Attempt(
	ctx context.Context,
	token string,
	userID string,
	competitionID string,
) (*Round1AttemptRecord, error) {

	filter := fmt.Sprintf(
		`user = "%s" && competition = "%s"`,
		userID,
		competitionID,
	)

	path := "/api/collections/attempts_round1/records" +
		"?filter=" + url.QueryEscape(filter) +
		"&perPage=1"

	var response struct {
		Items []Round1AttemptRecord `json:"items"`
	}

	err := pb.do(
		ctx,
		"GET",
		path,
		token,
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, nil
	}

	return &response.Items[0], nil
}

func (pb *PocketBaseClient) CreateRound1Attempt(
	ctx context.Context,
	token string,
	userID string,
	competitionID string,
	language string,
	startedAt string,
	serverDeadline string,
) (*Round1AttemptRecord, error) {

	payload := map[string]interface{}{
		"user":                userID,
		"competition":         competitionID,
		"language":            language,
		"started_at":          startedAt,
		"server_deadline":     serverDeadline,
		"question_started_at": startedAt,
		"last_saved_at":       startedAt,
		"current_question":    0,
		"answers":             map[string]interface{}{},
		"status":              "active",
		"correct_answers":     0,
		"score":               0,
		"time_taken_seconds":  0,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var response Round1AttemptRecord

	adminToken, err := pb.AdminToken(ctx)
	if err != nil {
		return nil, err
	}

	err = pb.do(
		ctx,
		http.MethodPost,
		"/api/collections/attempts_round1/records",
		adminToken,
		body,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (pb *PocketBaseClient) GetRound1AttemptByID(
	ctx context.Context,
	token string,
	attemptID string,
) (*Round1AttemptRecord, error) {

	path := "/api/collections/attempts_round1/records/" +
		url.PathEscape(attemptID)

	var response Round1AttemptRecord

	err := pb.do(
		ctx,
		http.MethodGet,
		path,
		token,
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (pb *PocketBaseClient) UpdateRound1Attempt(
	ctx context.Context,
	attemptID string,
	fields map[string]interface{},
) (*Round1AttemptRecord, error) {
	token, err := pb.AdminToken(ctx)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	path := "/api/collections/attempts_round1/records/" +
		url.PathEscape(attemptID)

	var response Round1AttemptRecord

	err = pb.do(
		ctx,
		http.MethodPatch,
		path,
		token,
		body,
		&response,
	)

	if err != nil {
		if apiErr, ok := err.(*APIError); ok &&
			apiErr.Status == http.StatusUnauthorized {

			pb.clearAdminToken()

			token, err = pb.AdminToken(ctx)
			if err != nil {
				return nil, err
			}

			err = pb.do(
				ctx,
				http.MethodPatch,
				path,
				token,
				body,
				&response,
			)
		}
	}

	if err != nil {
		return nil, err
	}

	return &response, nil
}
