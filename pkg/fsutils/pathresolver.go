package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetRelativePath returns the relative path of 'targetPath' with respect to 'basePath'.
// It supports both relative and absolute paths for 'targetPath'.
// If 'targetPath' is already an absolute path, it validates its existence.
// If 'targetPath' is relative, it computes the absolute path considering 'basePath' and validates.
// Returns an error if paths do not exist or are not relative to each other.
func GetRelativePath(targetPath, basePath string) (string, error) {
	absBasePath, err := filepath.Abs(filepath.Clean(basePath))
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute base path: %w", err)
	}
	if _, err := os.Stat(absBasePath); os.IsNotExist(err) {
		return "", fmt.Errorf("base path does not exist: %s", absBasePath)
	}

	absTargetPath := filepath.Clean(targetPath)
	if !filepath.IsAbs(targetPath) {
		absTargetPath = filepath.Join(absBasePath, targetPath)
	}

	if _, err := os.Stat(absTargetPath); os.IsNotExist(err) {
		return "", fmt.Errorf("target path does not exist: %s", absTargetPath)
	}

	relativePath, err := filepath.Rel(absBasePath, absTargetPath)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path: %w", err)
	}

	return relativePath, nil
}
