package config

import (
	"github.com/Excoriate/tftest/pkg/utils"
	"path/filepath"
)

// GetTerraformModuleTemplate returns the path to the terraform template directory based on the version.
func GetTerraformModuleTemplate(version string) (string, error) {
	if version == "" || version == "v1" {
		return filepath.Join(TemplateDir, TerraformModuleTemplateV1), nil
	}

	templateDirPath := filepath.Join(TemplateDir, version)

	if err := utils.DirExistAndHasContent(templateDirPath); err != nil {
		return "", err
	}

	return templateDirPath, nil
}

// GetTerraformModuleExampleTemplate returns the path to the terraform example template directory based on the version.
func GetTerraformModuleExampleTemplate(version string) (string, error) {
	if version == "" || version == "v1" {
		return filepath.Join(TemplateDir, TerraformModuleExampleTemplateV1), nil
	}

	templateDirPath := filepath.Join(TemplateDir, version)

	if err := utils.DirExistAndHasContent(templateDirPath); err != nil {
		return "", err
	}

	return templateDirPath, nil
}
