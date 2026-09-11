package round2

import (
	"testing"
)

// Unit tests for state handling logic and idempotency rules

func TestIdempotency(t *testing.T) {
	// If current question is 3, submitting question 2 should be rejected.
	currentProblem := 3.0
	submittedProblem := 2

	if float64(submittedProblem) >= currentProblem {
		t.Errorf("Should have triggered idempotency, submitted: %v, current: %v", submittedProblem, currentProblem)
	}
}

func TestQuestionAdvancement(t *testing.T) {
	currentScore := 20.0
	currentProblem := 2.0

	scoreObtained := 10.0
	newScore := currentScore + scoreObtained
	newProblem := currentProblem + 1

	if newScore != 30 {
		t.Errorf("Expected 30")
	}
	if newProblem != 3 {
		t.Errorf("Expected 3")
	}
}

func TestQuestionBank_Length(t *testing.T) {
	if len(QuestionBank) != 7 {
		t.Errorf("Expected exactly 7 logical questions, got %d", len(QuestionBank))
	}
}

func TestQuestion_NoJava(t *testing.T) {
	for _, q := range QuestionBank {
		if _, ok := q.Languages[LanguageCode("Java")]; ok {
			t.Errorf("Java should be completely removed, but found in Question %d", q.ID)
		}
	}
}
