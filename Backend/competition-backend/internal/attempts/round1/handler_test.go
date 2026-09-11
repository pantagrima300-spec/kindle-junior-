package round1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"competition-backend/internal/attempts/round1"
	"competition-backend/internal/auth"
	"competition-backend/internal/storage"
)

const (
	pbURL         = "http://127.0.0.1:8090"
	adminEmail    = "dhananjaysharma20.2006@gmail.com"
	adminPassword = "adolf_hitler"
)

func getAdminToken(t *testing.T) string {
	payload := map[string]string{
		"identity": adminEmail,
		"password": adminPassword,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(pbURL+"/api/collections/_superusers/auth-with-password", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to auth admin: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Token
}

func createTestUser(t *testing.T, token, email string) (string, string) {
	payload := map[string]string{
		"email":           email,
		"password":        "password123",
		"passwordConfirm": "password123",
		"name":            "Test User",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", pbURL+"/api/collections/users/records", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Logf("User creation response: %s", string(b))
	}

	authPayload := map[string]string{
		"identity": email,
		"password": "password123",
	}
	authBody, _ := json.Marshal(authPayload)
	authResp, err := http.Post(pbURL+"/api/collections/users/auth-with-password", "application/json", bytes.NewReader(authBody))
	if err != nil {
		t.Fatalf("failed to auth user: %v", err)
	}
	defer authResp.Body.Close()
	var result struct {
		Token  string `json:"token"`
		Record struct {
			ID string `json:"id"`
		} `json:"record"`
	}
	json.NewDecoder(authResp.Body).Decode(&result)
	return result.Record.ID, result.Token
}

func createTestCompetition(t *testing.T, token string, durationSeconds int) string {
	payload := map[string]interface{}{
		"name":             fmt.Sprintf("Test Comp %d", time.Now().UnixNano()),
		"status":           "active",
		"duration_seconds": durationSeconds,
		"max_attempts":     1,
		"auth_start":       time.Now().Add(-1 * time.Hour).Format(time.RFC3339Nano),
		"auth_end":         time.Now().Add(1 * time.Hour).Format(time.RFC3339Nano),
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", pbURL+"/api/collections/Competition/records", bytes.NewReader(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to create comp: %v", err)
	}
	defer resp.Body.Close()
	var result struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.ID
}

func setupAuthContext(req *http.Request, userID string) *http.Request {
	user := &storage.AuthRecord{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test",
	}
	ctx := auth.WithUser(req.Context(), user)
	return req.WithContext(ctx)
}

func TestRound1AttemptCreationAndFlow(t *testing.T) {
	adminToken := getAdminToken(t)
	userID, userToken := createTestUser(t, adminToken, fmt.Sprintf("test%d@example.com", time.Now().UnixNano()))
	compID := createTestCompetition(t, adminToken, 1800) // 30 minutes

	store := storage.New(pbURL, adminEmail, adminPassword)
	handler := round1.NewHandler(store)

	// 1. Start Attempt
	startReqBody, _ := json.Marshal(map[string]string{"competition_id": compID, "language": "C"})
	req := httptest.NewRequest("POST", "/api/round1/attempt", bytes.NewReader(startReqBody))
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	req = setupAuthContext(req, userID)

	rr := httptest.NewRecorder()
	handler.StartAttempt(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d: %s", rr.Code, rr.Body.String())
	}

	var attempt storage.Round1AttemptRecord
	json.NewDecoder(rr.Body).Decode(&attempt)

	if attempt.Status != "active" {
		t.Errorf("Expected active status, got %s", attempt.Status)
	}

	attemptID := attempt.ID

	// 2. Fetch Question 1
	req = httptest.NewRequest("GET", "/api/round1/question?attempt_id="+attemptID, nil)
	req.Header.Set("Authorization", userToken)
	req = setupAuthContext(req, userID)
	rr = httptest.NewRecorder()
	handler.GetQuestion(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}

	var qResp round1.Round1QuestionResponse
	json.NewDecoder(rr.Body).Decode(&qResp)

	if qResp.QuestionNumber != 0 {
		t.Errorf("Expected question 0, got %d", qResp.QuestionNumber)
	}

	// 3. Submit valid answer to Q1
	ansBody, _ := json.Marshal(round1.SubmitAnswerRequest{
		AttemptID:  attemptID,
		QuestionID: qResp.QuestionID,
		Answer:     1, // Assuming 1 is some option
	})
	req = httptest.NewRequest("POST", "/api/round1/answer", bytes.NewReader(ansBody))
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	req = setupAuthContext(req, userID)
	rr = httptest.NewRecorder()
	handler.SubmitAnswer(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK on submit, got %d: %s", rr.Code, rr.Body.String())
	}

	// 4. Fetch Question 2
	req = httptest.NewRequest("GET", "/api/round1/question?attempt_id="+attemptID, nil)
	req.Header.Set("Authorization", userToken)
	req = setupAuthContext(req, userID)
	rr = httptest.NewRecorder()
	handler.GetQuestion(rr, req)

	json.NewDecoder(rr.Body).Decode(&qResp)
	if qResp.QuestionNumber != 1 {
		t.Errorf("Expected question 1 after answering, got %d", qResp.QuestionNumber)
	}
}

func TestRound1ExpiredAttempt(t *testing.T) {
	adminToken := getAdminToken(t)
	userID, userToken := createTestUser(t, adminToken, fmt.Sprintf("test%d@example.com", time.Now().UnixNano()))

	// Competition expires instantly basically by simulating server_deadline in past
	compID := createTestCompetition(t, adminToken, -10)

	store := storage.New(pbURL, adminEmail, adminPassword)
	handler := round1.NewHandler(store)

	// 1. Start Attempt
	startReqBody, _ := json.Marshal(map[string]string{"competition_id": compID, "language": "C"})
	req := httptest.NewRequest("POST", "/api/round1/attempt", bytes.NewReader(startReqBody))
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	req = setupAuthContext(req, userID)

	rr := httptest.NewRecorder()
	handler.StartAttempt(rr, req)

	var attempt storage.Round1AttemptRecord
	json.NewDecoder(rr.Body).Decode(&attempt)
	attemptID := attempt.ID

	// Wait 1 second just in case
	time.Sleep(1 * time.Second)

	// 2. Fetch Question should fail due to expiry
	req = httptest.NewRequest("GET", "/api/round1/question?attempt_id="+attemptID, nil)
	req.Header.Set("Authorization", userToken)
	req = setupAuthContext(req, userID)
	rr = httptest.NewRecorder()
	handler.GetQuestion(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden for expired attempt, got %d", rr.Code)
	}

	// 3. Try submitting an answer, should fail
	ansBody, _ := json.Marshal(round1.SubmitAnswerRequest{
		AttemptID:  attemptID,
		QuestionID: 1,
		Answer:     1,
	})
	req = httptest.NewRequest("POST", "/api/round1/answer", bytes.NewReader(ansBody))
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	req = setupAuthContext(req, userID)
	rr = httptest.NewRecorder()
	handler.SubmitAnswer(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden on submit for expired attempt, got %d", rr.Code)
	}

	// 4. Verify in DB it is marked expired
	dbAttempt, _ := store.PocketBase.GetRound1AttemptByID(context.Background(), userToken, attemptID)
	if dbAttempt.Status != "expired" {
		t.Errorf("Expected status to be expired, got %s", dbAttempt.Status)
	}
}

func TestRound1DirectPocketBaseBypass(t *testing.T) {
	adminToken := getAdminToken(t)
	userID, userToken := createTestUser(t, adminToken, fmt.Sprintf("test%d@example.com", time.Now().UnixNano()))
	compID := createTestCompetition(t, adminToken, 1800)

	// Try to create an attempt directly using user token
	payload := map[string]interface{}{
		"user":            userID,
		"competition":     compID,
		"started_at":      time.Now().Format(time.RFC3339Nano),
		"server_deadline": time.Now().Add(100 * time.Hour).Format(time.RFC3339Nano),
		"status":          "active",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", pbURL+"/api/collections/attempts_round1/records", bytes.NewReader(body))
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		t.Fatalf("User was able to create an attempt directly via PocketBase! Security rules are missing.")
	}
}
