package storage

import (
	"context"
	"net/http"
	"net/url"
)

type CompetitionRecord struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	AuthStart       string `json:"auth_start"`
	AuthEnd         string `json:"auth_end"`
	DurationSeconds int    `json:"duration_seconds"`
	MaxAttempts     int    `json:"max_attempts"`
}

type CompetitionListResponse struct {
	Page       int                 `json:"page"`
	PerPage    int                 `json:"perPage"`
	TotalItems int                 `json:"totalItems"`
	TotalPages int                 `json:"totalPages"`
	Items      []CompetitionRecord `json:"items"`
}

func (pb *PocketBaseClient) GetCompetitions(
	ctx context.Context,
) (*CompetitionListResponse, error) {

	var response CompetitionListResponse

	err := pb.do(
		ctx,
		http.MethodGet,
		"/api/collections/Competition/records?"+url.Values{
			"perPage": {"50"},
			"sort":    {"auth_start"},
		}.Encode(),
		"",
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (pb *PocketBaseClient) GetCompetition(
	ctx context.Context,
	id string,
) (*CompetitionRecord, error) {

	var response CompetitionRecord

	err := pb.do(
		ctx,
		http.MethodGet,
		"/api/collections/Competition/records/"+url.PathEscape(id),
		"",
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
