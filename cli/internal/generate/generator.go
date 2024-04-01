package generate

import (
	"fmt"
	"github.com/Excoriate/tfgenctl/cli/internal/config"
	"github.com/Excoriate/tftest/pkg/utils"
	"os"
	"path/filepath"
	"text/template"
)

type ModuleGenerator interface {
	// GenerateCanonical generates a canonical module from the template directory.
	// 'templatedir' is the path to the template directory.
	// 'destpath' is the path to the destination directory.
	GenerateCanonical(templatedir, destpath string) (string, error)
	// Validate validates the module generation.
	// 'module' is the name of the module.
	// 'group' is the group name of the module. If the group is passed, normally it refers
	// to the 'service' (cloud provider) the module is for. E.g.: /modules/s3/s3_bucket
	Validate(module, group string, options *GenOptions) error
	// ResolveModuleDestPath resolves the destination path for the module.
	// 'module' is the name of the module.
	// 'group' is the group name of the module. If the group is passed, normally it refers
	// to the 'service' (cloud provider) the module is for. E.g.: /modules/s3/s3_bucket
	ResolveModuleDestPath(module, group string) (string, error)
	// CreateModulesDirectoryIfNotExists creates the modules directory if it does not exist.
	// 'repoRoot' is the root directory of the repository.
	CreateModulesDirectoryIfNotExists(repoRoot string) error
}

type GenOptions struct {
	FailIfNotGitRepo bool
}

func (c *Client) ResolveModuleDestPath(module, group string) (string, error) {
	if module == "" {
		return "", fmt.Errorf("failed to resolve destination path, module name is empty")
	}

	module = ConvertNameToCanonical(module)
	var path string

	repoRoot, err := c.Config.GetRepoRoot()
	if err != nil {
		return "", err
	}

	if group != "" {
		group = ConvertNameToCanonical(group)
		path = filepath.Join(repoRoot, "modules", group, module)
		c.Log.Info("the group is not empty, the module will be placed in the group %s", group)
	} else {
		path = filepath.Join(repoRoot, "modules", module)
	}

	c.Log.Info("the module will be placed in the path %s", path)

	if err := utils.DirExistAndHasContent(path); err != nil {
		return path, nil
	}

	return "", fmt.Errorf("module destination path already exists: %s", path)
}

func (c *Client) Validate(module, group string, options *GenOptions) error {
	if options == nil {
		return fmt.Errorf("generation failed, options are nil")
	}

	if _, err := c.ResolveModuleDestPath(module, group); err != nil {
		return err
	}

	_, isGitRepo := c.Config.IsGitRepo()

	if options.FailIfNotGitRepo && !isGitRepo {
		return fmt.Errorf("generation failed, the current directory is not a git repository")
	}

	return nil
}

func (c *Client) CreateModulesDirectoryIfNotExists(repoRoot string) error {
	modulesDir := filepath.Join(repoRoot, "modules")

	if err := utils.IsValidDirE(repoRoot); err != nil {
		return fmt.Errorf("failed to validate repo root: %w", err)
	}

	if err := utils.IsValidDirE(modulesDir); err != nil {
		// Create modules directory
		if err := os.MkdirAll(modulesDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create modules directory: %w", err)
		}
	}

	return nil
}

// GenerateCanonical generates a canonical module from the template directory.
// 'templatedir' is the path to the template directory.
// 'destpath' is the path to the destination directory.
func (c *Client) GenerateCanonical(templateVersion, targetPath string) (string, error) {
	templateSourcePath, err := config.GetTemplateTerraformDir(templateVersion)
	if err != nil {
		return "", err
	}

	if targetPath == "" {
		return "", fmt.Errorf("failed to generate canonical module, target path is empty")
	}

	err = filepath.Walk(templateSourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(templateSourcePath, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(targetPath, relPath)

			// Read template content
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template file: %w", err)
			}

			// Process template
			tmpl, err := template.New("template").Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse template: %w", err)
			}

			// Write to destination with processed content
			file, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file: %w", err)
			}

			defer file.Close()

			if err := tmpl.Execute(file, nil); // Replace `nil` with data struct if you need to pass data to templates
			err != nil {
				return fmt.Errorf("failed to execute template: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to process templates: %w", err)
	}

	return targetPath, nil
}
