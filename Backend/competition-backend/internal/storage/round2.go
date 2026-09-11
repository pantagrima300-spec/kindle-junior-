package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Round2AttemptRecord struct {
	ID               string  `json:"id"`
	User             string  `json:"user"`
	Competition      string  `json:"competition"`
	StartedAt        string  `json:"started_at"`
	ServerDeadline   string  `json:"server_deadline"`
	LastSavedAt      string  `json:"last_saved_at"`
	CurrentProblem   float64 `json:"current_problem"`
	ProblemsSolved   float64 `json:"problems_solved"`
	Status           string  `json:"status"`
	Score            float64 `json:"score"`
	TimeTakenSeconds float64 `json:"time_taken_seconds"`
	SubmittedAt      string  `json:"submitted_at"`
	SelectedLanguage string  `json:"selected_language"`
}

func (pb *PocketBaseClient) GetRound2Attempt(
	ctx context.Context,
	token string,
	userID string,
	competitionID string,
) (*Round2AttemptRecord, error) {

	filter := fmt.Sprintf(
		`user = "%s" && competition = "%s"`,
		userID,
		competitionID,
	)

	path := "/api/collections/attempts_round2/records" +
		"?filter=" + url.QueryEscape(filter) +
		"&perPage=1"

	var response struct {
		Items []Round2AttemptRecord `json:"items"`
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

func (pb *PocketBaseClient) CreateRound2Attempt(
	ctx context.Context,
	token string,
	userID string,
	competitionID string,
	startedAt string,
	serverDeadline string,
	initialLanguage string,
) (*Round2AttemptRecord, error) {

	payload := map[string]interface{}{
		"user":               userID,
		"competition":        competitionID,
		"started_at":         startedAt,
		"server_deadline":    serverDeadline,
		"last_saved_at":      startedAt,
		"current_problem":    1,
		"problems_solved":    0,
		"status":             "in_progress",
		"score":              0,
		"time_taken_seconds": 0,
		"selected_language":  initialLanguage,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var response Round2AttemptRecord

	adminToken, err := pb.AdminToken(ctx)
	if err != nil {
		return nil, err
	}

	err = pb.do(
		ctx,
		http.MethodPost,
		"/api/collections/attempts_round2/records",
		adminToken,
		body,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (pb *PocketBaseClient) GetRound2AttemptByID(
	ctx context.Context,
	token string,
	attemptID string,
) (*Round2AttemptRecord, error) {

	// Use admin token for reliable direct-ID lookups regardless of PocketBase rules
	adminToken, err := pb.AdminToken(ctx)
	if err != nil {
		return nil, err
	}

	path := "/api/collections/attempts_round2/records/" +
		url.PathEscape(attemptID)

	var response Round2AttemptRecord

	err = pb.do(
		ctx,
		http.MethodGet,
		path,
		adminToken,
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (pb *PocketBaseClient) UpdateRound2Attempt(
	ctx context.Context,
	attemptID string,
	fields map[string]interface{},
) (*Round2AttemptRecord, error) {
	token, err := pb.AdminToken(ctx)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	path := "/api/collections/attempts_round2/records/" +
		url.PathEscape(attemptID)

	var response Round2AttemptRecord

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
