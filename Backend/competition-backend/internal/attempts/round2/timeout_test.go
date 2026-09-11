package round2

import (
	"testing"
	"time"
)

func TestTimeoutLogic(t *testing.T) {
	// 70 minute exact timer logic
	startedAt := time.Now().UTC()
	serverDeadline := startedAt.Add(4200 * time.Second)

	// Time taken calculation clamped
	completionTime := startedAt.Add(4500 * time.Second)

	timeTaken := completionTime.Sub(startedAt).Seconds()
	if timeTaken > 4200 {
		timeTaken = 4200
	}

	if timeTaken != 4200 {
		t.Errorf("Expected time to be clamped to 4200, got %v", timeTaken)
	}

	// Test expired logic
	now := startedAt.Add(4201 * time.Second)
	if !now.After(serverDeadline) {
		t.Errorf("Expected 4201 seconds to be after deadline")
	}
}

func TestTimeTakenStandard(t *testing.T) {
	startedAt := time.Now().UTC()
	completionTime := startedAt.Add(38*time.Minute + 17*time.Second) // 2297 seconds

	timeTaken := completionTime.Sub(startedAt).Seconds()
	if timeTaken > 4200 {
		timeTaken = 4200
	}

	if int(timeTaken) != 2297 {
		t.Errorf("Expected 2297 seconds, got %v", timeTaken)
	}
}
