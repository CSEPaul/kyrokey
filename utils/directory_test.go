package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetDirectory(t *testing.T) {
	// Create a temporary base directory (auto-cleanup)
	baseDir := t.TempDir()

	// Target directory inside temp
	dirPath := filepath.Join(baseDir, "testdir")

	// 1. Should create directory successfully
	err := SetDirectory(dirPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify directory exists
	info, statErr := os.Stat(dirPath)
	if statErr != nil {
		t.Fatalf("expected directory to exist, got error %v", statErr)
	}
	if !info.IsDir() {
		t.Fatalf("expected a directory, got something else")
	}

	// 2. Calling again should fail (since os.Mkdir errors if exists)
	err = SetDirectory(dirPath)
	if err == nil {
		t.Fatalf("expected error when directory already exists, got nil")
	}
}

func TestDirectoryExists(t *testing.T) {
	baseDir := t.TempDir()

	existingDir := filepath.Join(baseDir, "existing")
	nonExistingDir := filepath.Join(baseDir, "missing")

	// Create a real directory
	if err := os.Mkdir(existingDir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// 1. Should return true for existing directory
	if !DirectoryExists(existingDir) {
		t.Errorf("expected true for existing directory, got false")
	}

	// 2. Should return false for non-existing directory
	if DirectoryExists(nonExistingDir) {
		t.Errorf("expected false for non-existing directory, got true")
	}

	// 3. Optional: test with a file (your function returns true)
	filePath := filepath.Join(baseDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if !DirectoryExists(filePath) {
		t.Errorf("expected true for existing file, got false")
	}
}
