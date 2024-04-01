package generate

import "github.com/Excoriate/tftest/pkg/utils"

// IsValidDir checks if the given path is a valid directory.
// It returns an error if the path is not a valid directory.
// Or if the directory isn't a directory.
func IsValidDir(path string) error {
	if err := utils.IsValidDirE(path); err != nil {
		return err
	}

	return nil
}
