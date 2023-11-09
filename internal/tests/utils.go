package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func CreateTempDir(t *testing.T) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	_, err = os.Create(filepath.Join(tempDir, "test1.gitignore"))
	if err != nil {
		t.Fatalf("Failed to create test1.gitignore file: %v", err)
	}
	_, err = os.Create(filepath.Join(tempDir, "test2.gitignore"))
	if err != nil {
		t.Fatalf("Failed to create test2.gitignore file: %v", err)
	}

	return tempDir, func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("Failed to remove temp directory: %v", err)
		}
	}
}
