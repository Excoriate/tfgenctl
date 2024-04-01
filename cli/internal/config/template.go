package config

import (
	"github.com/Excoriate/tftest/pkg/utils"
	"path/filepath"
)

// GetTemplateTerraformDir returns the path to the terraform template directory based on the version.
func GetTemplateTerraformDir(version string) (string, error) {
	if version == "" || version == "v1" {
		return filepath.Join(TemplateDir, TerraformTemplateV1), nil
	}

	templateDirPath := filepath.Join(TemplateDir, version)

	if err := utils.DirExistAndHasContent(templateDirPath); err != nil {
		return "", err
	}

	return templateDirPath, nil
}
