package round2

import (
	"time"

	"context"
	"strings"
	"testing"
)

func TestExecution_C_Success(t *testing.T) {
	code := `#include <stdio.h>
int main() { printf("Hello"); return 0; }`
	out, err, errType := ExecuteCode(func() context.Context {
		ctx, _ := context.WithTimeout(func() context.Context {
			ctx, _ := context.WithTimeout(func() context.Context {
				ctx, _ := context.WithTimeout(func() context.Context {
					ctx, _ := context.WithTimeout(func() context.Context {
						ctx, _ := context.WithTimeout(func() context.Context {
							ctx, _ := context.WithTimeout(func() context.Context {
								ctx, _ := context.WithTimeout(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), 5*time.Second)
								return ctx
							}(), 5*time.Second)
							return ctx
						}(), 5*time.Second)
						return ctx
					}(), 5*time.Second)
					return ctx
				}(), 5*time.Second)
				return ctx
			}(), 5*time.Second)
			return ctx
		}(), 5*time.Second)
		return ctx
	}(), LangC, code, "")
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
	if out != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", out)
	}
	if errType != "" {
		t.Errorf("Expected empty errType, got %s", errType)
	}
}

func TestExecution_Python_Success(t *testing.T) {
	code := `print("Hello")`
	out, err, errType := ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangPython, code, "")
	if err != nil {
		t.Fatalf("Expected success, got %v", err)
	}
	if strings.TrimSpace(out) != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", out)
	}
	if errType != "" {
		t.Errorf("Expected empty errType, got %s", errType)
	}
}

func TestExecution_CompileFailure(t *testing.T) {
	// C compilation failure
	code := `int main() { printf("missing include and syntax error) }`
	_, err, errType := ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangC, code, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if errType != "COMPILE_ERROR" {
		t.Errorf("Expected COMPILE_ERROR, got %s", errType)
	}

	// Python compilation failure
	pyCode := `print("missing quote)`
	_, err, errType = ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangPython, pyCode, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if errType != "COMPILE_ERROR" {
		t.Errorf("Expected COMPILE_ERROR, got %s", errType)
	}
}

func TestExecution_RuntimeFailure(t *testing.T) {
	code := `
import sys
sys.exit(1)
`
	_, err, errType := ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangPython, code, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if errType != "RUNTIME_ERROR" {
		t.Errorf("Expected RUNTIME_ERROR, got %s", errType)
	}
}

func TestExecution_Timeout(t *testing.T) {
	code := `
import time
while True:
    time.sleep(1)
`
	_, err, errType := ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangPython, code, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if errType != "TIME_LIMIT_EXCEEDED" {
		t.Errorf("Expected TIME_LIMIT_EXCEEDED, got %s", errType)
	}
}

func TestExecution_OutputLimit(t *testing.T) {
	// 2 MB of output should hit the 1MB limit
	code := `
print("A" * 2000000)
`
	_, err, errType := ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangPython, code, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if errType != "OUTPUT_LIMIT_EXCEEDED" {
		t.Errorf("Expected OUTPUT_LIMIT_EXCEEDED, got %s", errType)
	}
}

func TestExecution_NetworkDisabled(t *testing.T) {
	code := `
import urllib.request
try:
    urllib.request.urlopen("http://8.8.8.8", timeout=1)
except Exception as e:
    print(e)
`
	out, err, _ := ExecuteCode(func() context.Context { ctx, _ := context.WithTimeout(context.Background(), 5*time.Second); return ctx }(), LangPython, code, "")
	if err != nil {
		// It might return runtime error if network fetch is not caught properly,
		// but since we catch it, it should succeed and print Network is unreachable.
	}
	if !strings.Contains(out, "Network unreachable") && !strings.Contains(out, "urllib.error.URLError") {
		t.Errorf("Expected network unreachable error in output, got %s", out)
	}
}
