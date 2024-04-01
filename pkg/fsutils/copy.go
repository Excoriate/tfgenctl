package fsutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDirTo copies a directory from src to dest, including all nested directories and files.
func CopyDirTo(dest, src string) error {
	// Get properties of source dir
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source directory info: %w", err)
	}

	// Create the destination directory with the same permissions as the source
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("error creating destination directory: %w", err)
	}

	directory, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source directory: %w", err)
	}
	defer directory.Close()

	objects, err := directory.Readdir(-1) // Read all objects in the directory

	for _, obj := range objects {
		srcFullPath := filepath.Join(src, obj.Name())
		destFullPath := filepath.Join(dest, obj.Name())

		if obj.IsDir() {
			// Create sub-directories - recursively call CopyDirTo
			err = CopyDirTo(destFullPath, srcFullPath)
			if err != nil {
				return err
			}
		} else {
			// Perform the file copy
			err = CopyFile(destFullPath, srcFullPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copies a single file from src to dest
func CopyFile(dest, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // Copy the file
	if err != nil {
		return fmt.Errorf("error copying file content: %w", err)
	}

	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source file info: %w", err)
	}
	err = os.Chmod(dest, srcFileInfo.Mode()) // Set the same permissions
	if err != nil {
		return fmt.Errorf("error setting file permissions: %w", err)
	}

	return nil
}
