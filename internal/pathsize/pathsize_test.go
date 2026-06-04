package pathsize

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	tmpDir := t.TempDir()

	const (
		visibleFileSize = int64(5)
		hiddenFileSize  = int64(4)
		nestedFileSize  = int64(7)
	)

	visibleEntrySize := visibleFileSize
	mustWriteFile(t, filepath.Join(tmpDir, "file.txt"), "hello")
	mustWriteFile(t, filepath.Join(tmpDir, ".hidden.txt"), "hide")
	mustMkdir(t, filepath.Join(tmpDir, "empty"))
	mustMkdir(t, filepath.Join(tmpDir, "nested"))
	mustWriteFile(t, filepath.Join(tmpDir, "nested", "inside.txt"), "1234567")

	if runtime.GOOS != "windows" {
		mustSymlink(t, "file.txt", filepath.Join(tmpDir, "file.link"))
		visibleEntrySize += int64(len("file.txt"))
	}

	tests := []struct {
		name          string
		path          string
		includeHidden bool
		recursive     bool
		want          int64
	}{
		{
			name:          "regular file",
			path:          filepath.Join(tmpDir, "file.txt"),
			includeHidden: false,
			recursive:     false,
			want:          visibleFileSize,
		},
		{
			name:          "empty directory",
			path:          filepath.Join(tmpDir, "empty"),
			includeHidden: false,
			recursive:     false,
			want:          0,
		},
		{
			name:          "directory without recursive and hidden",
			path:          tmpDir,
			includeHidden: false,
			recursive:     false,
			want:          visibleEntrySize,
		},
		{
			name:          "directory with recursive without hidden",
			path:          tmpDir,
			includeHidden: false,
			recursive:     true,
			want:          visibleEntrySize + nestedFileSize,
		},
		{
			name:          "directory with recursive and hidden",
			path:          tmpDir,
			includeHidden: true,
			recursive:     true,
			want:          visibleEntrySize + nestedFileSize + hiddenFileSize,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Calculate(tt.path, tt.includeHidden, tt.recursive)
			if err != nil {
				t.Fatalf("Calculate() error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("Calculate() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalculateMissingPath(t *testing.T) {
	_, err := Calculate(filepath.Join(t.TempDir(), "missing"), false, false)
	if err == nil {
		t.Fatal("expected error for missing path")
	}

	if !strings.Contains(err.Error(), "stat path") {
		t.Fatalf("expected wrapped error context, got %q", err)
	}
}

func TestCalculateSymlinkAsArgument(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require elevated privileges on Windows")
	}

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target.txt")
	link := filepath.Join(tmpDir, "target.link")

	mustWriteFile(t, target, "123456789")
	mustSymlink(t, "target.txt", link)

	linkInfo, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("Lstat() error = %v", err)
	}

	got, err := Calculate(link, false, false)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	if got != linkInfo.Size() {
		t.Fatalf("Calculate() = %d, want symlink size %d", got, linkInfo.Size())
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()

	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("Mkdir(%q) error = %v", path, err)
	}
}

func mustSymlink(t *testing.T, target, link string) {
	t.Helper()

	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("Symlink(%q -> %q) error = %v", target, link, err)
	}
}
