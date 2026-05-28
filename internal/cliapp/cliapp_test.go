package cliapp

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrintsFormattedPathSize(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte(strings.Repeat("a", 2048)), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	err := New(&stdout).Run(context.Background(), []string{"hexlet-path-size", "--human", filePath})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := "2.0KB\t" + filePath + "\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunRejectsInvalidPathArgumentCount(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "missing path", args: []string{"hexlet-path-size"}},
		{name: "multiple paths", args: []string{"hexlet-path-size", ".", ".."}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			err := New(&stdout).Run(context.Background(), tt.args)
			if err == nil {
				t.Fatal("expected error")
			}

			if !strings.Contains(err.Error(), pathArgumentError) {
				t.Fatalf("error = %q, want %q", err, pathArgumentError)
			}

			if stdout.String() != "" {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
		})
	}
}

func TestRunHelpShowsRequiredPath(t *testing.T) {
	var stdout bytes.Buffer
	err := New(&stdout).Run(context.Background(), []string{"hexlet-path-size", "--help"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(stdout.String(), "<path>") {
		t.Fatalf("stdout = %q, want help with %q", stdout.String(), "<path>")
	}
}
