package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIRequiresPath(t *testing.T) {
	stdout, stderr, err := runCLI(t)
	if err == nil {
		t.Fatal("expected non-zero exit status")
	}

	var exitErr *exec.ExitError
	if !asExitError(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}

	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1", exitErr.ExitCode())
	}

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}

	if !strings.Contains(stderr, "exactly one path argument is required") {
		t.Fatalf("expected stderr to contain %q, got %q", "exactly one path argument is required", stderr)
	}
}

func TestCLIRejectsMultiplePaths(t *testing.T) {
	stdout, stderr, err := runCLI(t, ".", "..")
	if err == nil {
		t.Fatal("expected non-zero exit status")
	}

	var exitErr *exec.ExitError
	if !asExitError(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}

	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code = %d, want 1", exitErr.ExitCode())
	}

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}

	if !strings.Contains(stderr, "exactly one path argument is required") {
		t.Fatalf("expected stderr to contain %q, got %q", "exactly one path argument is required", stderr)
	}
}

func TestCLIPrintsHumanReadableSize(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	content := strings.Repeat("a", 2048)

	if err := osWriteFile(filePath, []byte(content)); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	stdout, stderr, err := runCLI(t, "--human", filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr: %s", err, stderr)
	}

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	if !strings.Contains(stdout, "2.0KB") {
		t.Fatalf("expected output to contain %q, got %q", "2.0KB", stdout)
	}
}

func TestCLIAllAndRecursiveFlags(t *testing.T) {
	tmpDir := t.TempDir()
	if err := osMkdir(filepath.Join(tmpDir, "nested")); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := osWriteFile(filepath.Join(tmpDir, ".hidden"), []byte("1234")); err != nil {
		t.Fatalf("write hidden file: %v", err)
	}
	if err := osWriteFile(filepath.Join(tmpDir, "nested", "inner.txt"), []byte("abcd")); err != nil {
		t.Fatalf("write nested file: %v", err)
	}

	stdout, stderr, err := runCLI(t, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "0B") {
		t.Fatalf("expected non-recursive output 0B, got %q", stdout)
	}

	stdout, stderr, err = runCLI(t, "--recursive", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "4B") {
		t.Fatalf("expected recursive output 4B, got %q", stdout)
	}

	stdout, stderr, err = runCLI(t, "--recursive", "--all", tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v, stderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "8B") {
		t.Fatalf("expected recursive+all output 8B, got %q", stdout)
	}
}

func TestCLIHelpShowsRequiredPath(t *testing.T) {
	stdout, stderr, err := runCLI(t, "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	if !strings.Contains(stdout, "<path>") {
		t.Fatalf("expected help output to contain %q, got %q", "<path>", stdout)
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
	if err != nil {
		return stdout.String(), stderr.String(), err
	}

	return stdout.String(), stderr.String(), nil
}

func projectCmdDir(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current file path")
	}

	return filepath.Dir(currentFile)
}

func asExitError(err error, target **exec.ExitError) bool {
	if err == nil {
		return false
	}

	var exitErr *exec.ExitError
	ok := errors.As(err, &exitErr)
	if !ok {
		return false
	}

	*target = exitErr
	return true
}

func osWriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

func osMkdir(path string) error {
	return os.Mkdir(path, 0o755)
}
