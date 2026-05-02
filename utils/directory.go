package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"go.uber.org/zap"
)

// SetDirectory creates a new directory at the specified file path with
// permissions set to 0755 (read/write/execute for owner, read/execute for group/others).
//
// Parameters:
//   - filePath: The path of the directory to create. If parent directories
//     do not exist, this function will fail.
//
// Returns:
//   - error: Returns any error encountered while creating the directory. If
//     the directory is created successfully, error will be nil.
//
// Behavior:
//   - Prints an error message to stdout using fmt.Println if directory creation fails.
//   - Does not create parent directories recursively (use os.MkdirAll for that).
//
// Example:
//
//	err := SetDirectory("./logs")
//	if err != nil {
//	    log.Fatalf("Failed to create log directory: %v", err)
//	}
func SetDirectory(filePath string) error {
	err := os.Mkdir(filePath, 0755)
	if err != nil {
		fmt.Println("Error creating directory", zap.Error(err))
	}
	return err
}

// DirectoryExists checks whether a directory or file exists at the specified path.
//
// Parameters:
//   - filePath: The path of the directory or file to check.
//
// Returns:
//   - bool: Returns true if the directory/file exists, false otherwise.
//
// Behavior:
//   - Uses os.Stat to determine existence.
//   - If the file or directory does not exist, logs an error using the global Zap logger.
//   - Does not create the directory; only checks for existence.
//
// Example:
//
//	exists := DirectoryExists("./logs")
//	if !exists {
//	    fmt.Println("Directory does not exist")
//	}
func DirectoryExists(filePath string) bool {

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		zap.L().Error("No directory or file found", zap.String("filepath", filePath))
		return false
	}
	return true
}

// SetDirectoryIfNotExists ensures that a directory exists at the specified path.
// If the directory does not exist, it attempts to create it.
//
// Parameters:
//   - filePath: The path of the directory to check or create.
//
// Returns:
//   - error: Returns any error encountered while creating the directory. If the
//     directory already exists, returns nil.
//
// Behavior:
//   - Calls DirectoryExists to check if the directory already exists.
//   - If the directory does not exist, calls SetDirectory to create it.
//   - Errors during creation are returned to the caller.
//
// Example:
//
//	err := SetDirectoryIfNotExists("./logs")
//	if err != nil {
//	    log.Fatalf("Failed to ensure log directory: %v", err)
//	}
func SetDirectoryIfNotExists(filePath string) error {
	if !DirectoryExists(filePath) {
		SetDirectory(filePath)
	}
	return nil
}

// GetDirectory retrieves the current working directory of the application.
//
// Returns:
//   - string: The absolute path of the current working directory.
//   - error:  Any error encountered while retrieving the working directory.
//
// Behavior:
//   - Uses os.Getwd to obtain the current working directory.
//   - Logs an error using the global Zap logger if it fails to get the directory.
//
// Example:
//
//	cwd, err := GetDirectory()
//	if err != nil {
//	    log.Fatalf("Cannot get current directory: %v", err)
//	}
//	fmt.Println("Current working directory:", cwd)
func GetDirectory() (string, error) {
	// Get the current working directory
	path, err := os.Getwd()
	if err != nil {
		zap.L().Error("Error getting current working directory - make sure app permission are set", zap.Error(err))
		return "", err
	}
	return path, nil
}

// DeleteDirectory deletes the specified directory and all its contents
func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		zap.L().Error("Unable to delete the directory", zap.Error(err))
		return fmt.Errorf("failed to delete directory %s: %w", dirPath, err)
	}
	return nil
}

// EnsureExecutable verifies that the binary exists and is runnable.
func EnsureExecutable(binary string) error {
	path, err := exec.LookPath(binary)
	if err != nil {
		zap.L().Error(binary+"not found on PATH", zap.Error(err))
		return fmt.Errorf("%s not found on PATH", binary)
	}

	// Unix-like systems: check executable bit
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			zap.L().Error("cannot stat"+path, zap.Error(err))
			return fmt.Errorf("cannot stat %s: %w", path, err)
		}
		//Prevents permission errors
		if info.Mode()&0111 == 0 {
			zap.L().Error(path + "exists but is not executable")
			return fmt.Errorf("%s exists but is not executable", path)
		}
	}

	// Final check: try to start the process
	cmd := exec.Command(path, "--help")
	if err := cmd.Start(); err != nil {
		zap.L().Error("binary exists but cannot be executed")
		return errors.New("binary exists but cannot be executed")
	}
	_ = cmd.Process.Kill()

	return nil
}
