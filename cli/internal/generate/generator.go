package generate

import (
	"fmt"
	"github.com/Excoriate/tfgenctl/cli/internal/config"
	"github.com/Excoriate/tftest/pkg/tmplutils"
	"github.com/Excoriate/tftest/pkg/utils"
	"os"
	"path/filepath"
	"text/template"
)

type ModuleGenerator interface {
	// GenerateModule GenerateCanonical generates a canonical module from the template directory.
	// 'templatedir' is the targetModulePath to the template directory.
	// 'destpath' is the targetModulePath to the destination directory.
	GenerateModule(templatedir string, paths *TargetPaths) error
	GenerateModuleExample(templatedir string, paths *TargetPaths) error
	// Validate validates the module generation.
	// 'module' is the name of the module.
	// 'group' is the group name of the module. If the group is passed, normally it refers
	// to the 'service' (cloud provider) the module is for. E.g.: /modules/s3/s3_bucket
	Validate(module, group string, options *GenOptions) error
	// ResolvePaths ResolveModuleDestPath resolves the destination targetModulePath for the module.
	// 'module' is the name of the module.
	// 'group' is the group name of the module. If the group is passed, normally it refers
	// to the 'service' (cloud provider) the module is for. E.g.: /modules/s3/s3_bucket
	ResolvePaths(module, group string) (*TargetPaths, error)
	// CreateBaseDirsIfNotExist CreateModulesDirectoryIfNotExists creates the modules directory if it does not exist.
	// 'repoRoot' is the root directory of the repository.
	CreateBaseDirsIfNotExist(repoRoot string) error
	CreateTargetPath(targetPath string) error
	GenerateTerraformFiles(sourcePath string, destPath string) error
	GenerateFile(sourcePath string, destPath string, funcMap template.FuncMap, data interface{}) error
}

type GenOptions struct {
	FailIfNotGitRepo bool
}

type TargetPaths struct {
	ModulePath  string
	ExamplePath string
	TestPath    string
}

func (c *Client) ResolvePaths(module, group string) (*TargetPaths, error) {
	if module == "" {
		return nil, fmt.Errorf("failed to resolve destination targetModulePath, module name is empty")
	}

	module = ConvertNameToCanonical(module)
	var targetModulePath string
	var targetExamplePath string
	var targetTestPath string

	repoRoot, err := c.Config.GetRepoRoot()
	if err != nil {
		return nil, err
	}

	tp := &TargetPaths{}

	if group != "" {
		group = ConvertNameToCanonical(group)
		targetModulePath = filepath.Join(repoRoot, "modules", group, module)
		targetExamplePath = filepath.Join(repoRoot, "examples", group, module)
		targetTestPath = filepath.Join(repoRoot, "tests", group, module)
		c.Log.Info("the group is not empty, the module will be placed in the group %s", group)
	} else {
		targetModulePath = filepath.Join(repoRoot, "modules", module)
		targetExamplePath = filepath.Join(repoRoot, "examples", module)
		targetTestPath = filepath.Join(repoRoot, "tests", module)
	}

	c.Log.Info("the module will be placed in the targetModulePath %s", targetModulePath)

	// If dir exist, and has content, fail with an error
	if err := utils.DirExistAndHasContent(targetModulePath); err == nil {
		return nil, fmt.Errorf("module destination targetModulePath already exists: %s", targetModulePath)
	}

	// Target module path
	tp.ModulePath = targetModulePath
	tp.ExamplePath = targetExamplePath
	tp.TestPath = targetTestPath

	return tp, nil
}

func (c *Client) Validate(module, group string, options *GenOptions) error {
	if options == nil {
		return fmt.Errorf("generation failed, options are nil")
	}

	if _, err := c.ResolvePaths(module, group); err != nil {
		return err
	}

	_, isGitRepo := c.Config.IsGitRepo()

	if options.FailIfNotGitRepo && !isGitRepo {
		return fmt.Errorf("generation failed, the current directory is not a git repository")
	}

	return nil
}

func (c *Client) CreateBaseDirsIfNotExist(repoRoot string) error {
	modulesDir := filepath.Join(repoRoot, config.TargetModuleDir)
	examplesDir := filepath.Join(repoRoot, config.TargetModuleExampleDir)
	testsDir := filepath.Join(repoRoot, config.TargetModuleTestDir)

	if err := utils.IsValidDirE(repoRoot); err != nil {
		return fmt.Errorf("failed to validate repo root: %w", err)
	}

	if err := utils.IsValidDirE(modulesDir); err != nil {
		// Create modules directory
		if err := os.MkdirAll(modulesDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create modules directory: %w", err)
		}
	}

	if err := utils.IsValidDirE(examplesDir); err != nil {
		// Create examples directory
		if err := os.MkdirAll(examplesDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create examples directory: %w", err)
		}
	}

	if err := utils.IsValidDirE(testsDir); err != nil {
		// Create tests directory
		if err := os.MkdirAll(testsDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create tests directory: %w", err)
		}
	}

	return nil
}

func (c *Client) CreateTargetPath(targetPath string) error {
	if targetPath == "" {
		return fmt.Errorf("failed to create target targetModulePath, targetModulePath is empty")
	}

	if err := utils.IsValidDirE(targetPath); err != nil {
		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create target targetModulePath: %w", err)
		}
	}

	return nil
}

func (c *Client) createAllTargetPaths(paths *TargetPaths) error {
	// Create the necessary paths in the root repository.
	if err := c.CreateTargetPath(paths.ModulePath); err != nil {
		return err
	}

	if err := c.CreateTargetPath(paths.ExamplePath); err != nil {
		return err
	}

	if err := c.CreateTargetPath(paths.TestPath); err != nil {
		return err
	}

	return nil
}

// GenerateModule generates a canonical module from the template directory.
func (c *Client) GenerateModule(templateVersion string, paths *TargetPaths) error {
	templateSourcePath, err := config.GetTerraformModuleTemplate(templateVersion)
	if err != nil {
		return err
	}

	if paths == nil {
		return fmt.Errorf("failed to generate canonical module, target paths are nil")
	}

	// Generate the target paths.
	if err := c.createAllTargetPaths(paths); err != nil {
		return err
	}

	// Generate the terraform files
	if err := c.GenerateTerraformFiles(templateSourcePath, paths.ModulePath); err != nil {
		return err
	}

	return nil
}

// GenerateModuleExample generates an example module from the template directory.
func (c *Client) GenerateModuleExample(templateVersion string, paths *TargetPaths) error {
	templateSourcePath, err := config.GetTerraformModuleExampleTemplate(templateVersion)
	if err != nil {
		return err
	}

	if paths == nil {
		return fmt.Errorf("failed to generate example module, target paths are nil")
	}

	// Generate the target paths.
	if err := c.createAllTargetPaths(paths); err != nil {
		return err
	}

	// Generate the terraform files
	if err := c.GenerateTerraformFiles(templateSourcePath, paths.ExamplePath); err != nil {
		return err
	}

	return nil
}

// GenerateTerraformFiles generates files from templates, allowing individual file treatment.
func (c *Client) GenerateTerraformFiles(sourcePath, targetPath string) error {
	funcMap := template.FuncMap{
		"include": func(string) string { return "" }, // No-op for the example
	}

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(sourcePath, path)
			if err != nil {
				return err
			}

			// Modify the file name if needed (e.g., remove .tmpl extension)
			relPath = CleanBeforeCopy(relPath)

			destPath := filepath.Join(targetPath, relPath)

			// Process the template file
			if err := c.GenerateFile(path, destPath, funcMap, nil); err != nil {
				return err
			}
		}
		return nil
	})
}

// GenerateFile generates a file from a template.
// It's a wrapper around the ProcessTemplFile function that's used from the TFTest library.
func (c *Client) GenerateFile(sourcePath, destPath string, funcMap template.FuncMap, data interface{}) error {
	return tmplutils.ProcessTemplFile(sourcePath, destPath, funcMap, data)
}
