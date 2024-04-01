package fsutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRelativePath(t *testing.T) {
	tempDir := t.TempDir()

	basePath := filepath.Join(tempDir, "base")
	targetPath := filepath.Join(basePath, "target")
	_ = os.MkdirAll(targetPath, os.ModePerm)

	tests := []struct {
		name        string
		targetPath  string
		basePath    string
		wantRelPath string
		expectErr   bool
		setupFunc   func() (cleanupFunc func(), err error)
	}{
		{
			name:        "absolute path to target",
			targetPath:  targetPath,
			basePath:    basePath,
			wantRelPath: "target",
			expectErr:   false,
		},
		{
			name:        "relative path to target",
			targetPath:  "target",
			basePath:    basePath,
			wantRelPath: "target",
			expectErr:   false,
		},
		{
			name:       "non-existent target path",
			targetPath: filepath.Join(tempDir, "nonexistent"),
			basePath:   basePath,
			expectErr:  true,
		},
		{
			name:       "non-existent base path",
			targetPath: targetPath,
			basePath:   filepath.Join(tempDir, "nonexistent"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				cleanupFunc, err := tt.setupFunc()
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				defer cleanupFunc()
			}

			gotRelPath, err := GetRelativePath(tt.targetPath, tt.basePath)
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got none")
			} else if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectErr && gotRelPath != tt.wantRelPath {
				t.Errorf("expected relative path %s, got %s", tt.wantRelPath, gotRelPath)
			}
		})
	}
}
