package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIPrintsHumanReadableSize(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")

	if err := os.WriteFile(filePath, []byte(strings.Repeat("a", 2048)), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	stdout, stderr, err := runCLI(t, "--human", filePath)
	if err != nil {
		t.Fatalf("runCLI() error = %v, stderr = %q", err, stderr)
	}

	want := "2.0KB\t" + filePath + "\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}

	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestCLIReportsUnknownFlagOnce(t *testing.T) {
	stdout, stderr, err := runCLI(t, "--unknown")
	if err == nil {
		t.Fatal("runCLI() error = nil, want non-zero exit")
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("runCLI() error = %T, want *exec.ExitError", err)
	}

	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1", exitErr.ExitCode())
	}

	if !strings.Contains(stdout, "<path>") {
		t.Fatalf("stdout = %q, want help with required path", stdout)
	}

	if got := strings.Count(stderr, "flag provided but not defined: -unknown"); got != 1 {
		t.Fatalf("unknown flag message count = %d, want 1; stderr = %q", got, stderr)
	}
}

func runCLI(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	cmdArgs := append([]string{"run", "."}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = projectCmdDir(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func projectCmdDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current file path")
	}

	return filepath.Dir(currentFile)
}
