package cliapp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
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

func TestRunRequiresExactlyOnePath(t *testing.T) {
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
			var stderr bytes.Buffer
			exitCode := runWithTestExitHandler(t, &stdout, &stderr, tt.args)

			if stdout.String() != "" {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}

			if exitCode != 1 {
				t.Fatalf("exit code = %d, want 1", exitCode)
			}

			if !strings.Contains(stderr.String(), pathArgumentMessage) {
				t.Fatalf("stderr = %q, want message with %q", stderr.String(), pathArgumentMessage)
			}
		})
	}
}

func TestRunAllAndRecursiveFlags(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmpDir, "nested"), 0o755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("1234"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "nested", "inner.txt"), []byte("abcd"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "non-recursive visible", args: []string{"hexlet-path-size", tmpDir}, want: "0B\t" + tmpDir + "\n"},
		{name: "recursive visible", args: []string{"hexlet-path-size", "--recursive", tmpDir}, want: "4B\t" + tmpDir + "\n"},
		{name: "recursive all", args: []string{"hexlet-path-size", "--recursive", "--all", tmpDir}, want: "8B\t" + tmpDir + "\n"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			err := New(&stdout).Run(context.Background(), tt.args)
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if stdout.String() != tt.want {
				t.Fatalf("stdout = %q, want %q", stdout.String(), tt.want)
			}
		})
	}
}

func runWithTestExitHandler(t *testing.T, stdout, stderr *bytes.Buffer, args []string) int {
	t.Helper()

	exitCode := 0
	app := New(stdout)
	app.ErrWriter = stderr
	app.ExitErrHandler = func(_ context.Context, cmd *cli.Command, err error) {
		var exitCoder cli.ExitCoder
		if errors.As(err, &exitCoder) {
			exitCode = exitCoder.ExitCode()
		} else {
			exitCode = 1
		}

		_, _ = fmt.Fprintln(cmd.Root().ErrWriter, err)
	}

	if err := app.Run(context.Background(), args); err == nil {
		t.Fatal("Run() error = nil, want error")
	}

	return exitCode
}
