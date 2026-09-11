package round2

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type JudgeResult struct {
	Passed         bool   `json:"passed"`
	ActualOutput   string `json:"actual"`
	ExpectedOutput string `json:"expected"`
	Error          string `json:"error,omitempty"`
}

type ExecutionResult struct {
	Success     bool          `json:"success"`
	Status      string        `json:"status"` // "completed" etc
	Language    LanguageCode  `json:"language"`
	QuestionID  int           `json:"question_id"`
	Tests       []JudgeResult `json:"tests"`
	PassedTests int           `json:"passed_tests"`
	TotalTests  int           `json:"total_tests"`
}

// LimitedWriter stops writing and discards data after MaxBytes
type LimitedWriter struct {
	buf      bytes.Buffer
	MaxBytes int
	written  int
}

func (l *LimitedWriter) Write(p []byte) (n int, err error) {
	if l.written >= l.MaxBytes {
		return len(p), nil // pretend we wrote it to avoid breaking the pipe
	}
	toWrite := len(p)
	if l.written+toWrite > l.MaxBytes {
		toWrite = l.MaxBytes - l.written
	}
	n, err = l.buf.Write(p[:toWrite])
	l.written += n
	if err != nil {
		return n, err
	}
	return len(p), nil
}

func (l *LimitedWriter) String() string {
	return l.buf.String()
}

func ExecuteCode(ctx context.Context, lang LanguageCode, code string, input string) (string, error, string) {
	encodedCode := base64.StdEncoding.EncodeToString([]byte(code))

	var shellCmd string
	var containerImage string

	switch lang {
	case LangC:
		containerImage = "gcc:13"
		shellCmd = fmt.Sprintf("echo %s | base64 -d > /tmp/solution.c && gcc -O2 /tmp/solution.c -o /tmp/prog; if [ $? -ne 0 ]; then exit 42; fi; timeout 3 /tmp/prog", encodedCode)
	case LangPython:
		containerImage = "python:3.10-alpine"
		shellCmd = fmt.Sprintf("echo %s | base64 -d > /tmp/solution.py && python3 -m py_compile /tmp/solution.py; if [ $? -ne 0 ]; then exit 42; fi; timeout 3 python3 /tmp/solution.py", encodedCode)
	default:
		return "", fmt.Errorf("unsupported language: %s", lang), "RUNTIME_ERROR"
	}

	// Security flags:
	// --network none
	// --cpus=0.5
	// -m 128m
	// --pids-limit 64
	// --user nobody (we use nobody user for isolation, but ensure /tmp is writable)
	cmd := exec.CommandContext(ctx, "podman", "run", "--rm", "-i", "--net", "none", "--cpus=0.5", "-m", "128m", "--pids-limit", "64", "--workdir", "/tmp", containerImage, "sh", "-c", shellCmd)

	cmd.Stdin = strings.NewReader(input)

	stdout := &LimitedWriter{MaxBytes: 1024 * 1024} // 1MB limit
	stderr := &LimitedWriter{MaxBytes: 1024 * 1024} // 1MB limit
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	outStr := stdout.String() + stderr.String()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return outStr, fmt.Errorf("TIME_LIMIT_EXCEEDED"), "TIME_LIMIT_EXCEEDED"
		}

		exitCode := 1
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}

		if exitCode == 42 {
			return outStr, fmt.Errorf("COMPILE_ERROR"), "COMPILE_ERROR"
		} else if exitCode == 143 || exitCode == 124 || strings.Contains(stderr.String(), "Terminated") {
			return outStr, fmt.Errorf("TIME_LIMIT_EXCEEDED"), "TIME_LIMIT_EXCEEDED"
		} else if exitCode == 137 { // OOM kill
			return outStr, fmt.Errorf("RUNTIME_ERROR: OOM"), "RUNTIME_ERROR"
		} else if stdout.written >= stdout.MaxBytes {
			return outStr, fmt.Errorf("OUTPUT_LIMIT_EXCEEDED"), "OUTPUT_LIMIT_EXCEEDED"
		}

		return outStr, fmt.Errorf("RUNTIME_ERROR: %v", err), "RUNTIME_ERROR"
	}

	if stdout.written >= stdout.MaxBytes {
		return outStr, fmt.Errorf("OUTPUT_LIMIT_EXCEEDED"), "OUTPUT_LIMIT_EXCEEDED"
	}

	return outStr, nil, ""
}

func normalizeOutput(out string) string {
	out = strings.ReplaceAll(out, "\r\n", "\n")
	out = strings.TrimSpace(out)
	return out
}

func JudgeSubmission(ctx context.Context, question *Question, lang LanguageCode, code string, isSubmit bool) ExecutionResult {
	res := ExecutionResult{
		Success:    false,
		Status:     "failed",
		Language:   lang,
		QuestionID: question.ID,
		Tests:      []JudgeResult{},
	}

	total := 0
	passed := 0

	for _, tc := range question.TestCases {
		if tc.IsHidden && !isSubmit {
			continue
		}

		total++

		runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		actualRaw, err, errType := ExecuteCode(runCtx, lang, code, tc.Input)
		cancel()

		actual := normalizeOutput(actualRaw)
		expected := normalizeOutput(tc.ExpectedOutput)

		tcRes := JudgeResult{
			ExpectedOutput: expected,
			ActualOutput:   actual,
		}

		if err != nil {
			tcRes.Passed = false
			tcRes.Error = errType
			if errType == "COMPILE_ERROR" || errType == "RUNTIME_ERROR" {
				tcRes.ActualOutput = actualRaw
			}
		} else {
			if actual == expected {
				tcRes.Passed = true
				passed++
			} else {
				tcRes.Passed = false
				tcRes.Error = "WRONG_ANSWER"
			}
		}

		if tc.IsHidden {
			tcRes.ExpectedOutput = "HIDDEN"
			tcRes.ActualOutput = "HIDDEN"
		}

		res.Tests = append(res.Tests, tcRes)
	}

	res.TotalTests = total
	res.PassedTests = passed

	if passed == total && total > 0 {
		res.Success = true
		res.Status = "completed"
	}

	return res
}
