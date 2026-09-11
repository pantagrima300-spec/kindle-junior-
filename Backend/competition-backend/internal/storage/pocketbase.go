package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type PocketBaseClient struct {
	BaseURL    string
	HTTPClient *http.Client

	AdminEmail    string
	AdminPassword string

	adminMu    sync.Mutex
	adminToken string
}

func (pb *PocketBaseClient) HealthCheck(ctx context.Context) error {
	return pb.do(
		ctx,
		"GET",
		"/api/health",
		"",
		nil,
		nil,
	)
}

func NewPocketBaseClient(
	baseURL string,
	adminEmail string,
	adminPassword string,
) *PocketBaseClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,

		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		ForceAttemptHTTP2: true,
	}

	return &PocketBaseClient{
		BaseURL: strings.TrimRight(baseURL, "/"),

		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   5 * time.Second,
		},
		AdminEmail:    adminEmail,
		AdminPassword: adminPassword,
	}
}

func (pb *PocketBaseClient) AdminToken(ctx context.Context) (string, error) {
	pb.adminMu.Lock()
	defer pb.adminMu.Unlock()

	return pb.adminTokenLocked(ctx, false)
}

func (pb *PocketBaseClient) adminTokenLocked(
	ctx context.Context,
	forceRefresh bool,
) (string, error) {
	if pb.adminToken != "" && !forceRefresh {
		return pb.adminToken, nil
	}

	if pb.AdminEmail == "" || pb.AdminPassword == "" {
		return "", fmt.Errorf("PocketBase admin credentials are not configured")
	}

	payload := map[string]string{
		"identity": pb.AdminEmail,
		"password": pb.AdminPassword,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	var response struct {
		Token string `json:"token"`
	}

	err = pb.do(
		ctx,
		http.MethodPost,
		"/api/collections/_superusers/auth-with-password",
		"",
		body,
		&response,
	)
	if err != nil {
		return "", fmt.Errorf(
			"PocketBase admin authentication failed: %w",
			err,
		)
	}

	if response.Token == "" {
		return "", fmt.Errorf(
			"PocketBase admin authentication returned empty token",
		)
	}

	pb.adminToken = response.Token

	return pb.adminToken, nil
}

func (pb *PocketBaseClient) clearAdminToken() {
	pb.adminMu.Lock()
	defer pb.adminMu.Unlock()

	pb.adminToken = ""
}

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("PocketBase API error: status=%d message=%s", e.Status, e.Message)
}

func (pb *PocketBaseClient) do(
	ctx context.Context,
	method string,
	path string,
	token string,
	body []byte,
	result any,
) error {
	url := pb.BaseURL + path

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		strings.NewReader(string(body)),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")

	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	resp, err := pb.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}

		_ = json.NewDecoder(resp.Body).Decode(&apiErr)

		message := apiErr.Message

		if len(apiErr.Data) > 0 {
			dataBytes, err := json.Marshal(apiErr.Data)
			if err == nil {
				message += " data=" + string(dataBytes)
			}
		}

		return &APIError{
			Status:  resp.StatusCode,
			Message: message,
		}
	}

	if result == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(result)
}
