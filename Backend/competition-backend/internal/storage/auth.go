package storage

import (
	"context"
)

type AuthRecord struct {
	ID             string `json:"id"`
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	Email          string `json:"email"`
	Verified       bool   `json:"verified"`
	Name           string `json:"name"`
}

type AuthRefreshResponse struct {
	Token  string     `json:"token"`
	Record AuthRecord `json:"record"`
}

func (pb *PocketBaseClient) ValidateToken(
	ctx context.Context,
	token string,
) (*AuthRefreshResponse, error) {

	var response AuthRefreshResponse

	err := pb.do(
		ctx,
		"POST",
		"/api/collections/users/auth-refresh",
		token,
		nil,
		&response,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
