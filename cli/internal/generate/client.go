package generate

import (
	"github.com/Excoriate/tfgenctl/cli/internal/config"
	"github.com/Excoriate/tfgenctl/pkg/logger"
	"github.com/Excoriate/tftest/pkg/utils"
	"os"
)

const canonicalModuleDestDir = "modules"

type Config interface {
	GetDestDir() (string, error)
	GetTFTemplateDirV1() (string, error)
	GetModulesDir() string
	GetHomeDir() string
	GetCurrentDir() string
	IsGitRepo() (string, bool)
	GetRepoRoot() (string, error)
}

type Options struct {
	destDir           string
	tfTemplateVersion string
	homeDir           string
	currentDir        string
	isGitRepo         bool
	repoRoot          string
}

type Generator struct{}

func (o *Options) GetDestDir() (string, error) {
	return o.destDir, nil
}

func (o *Options) GetTFTemplateDirV1() (string, error) {
	templateDir, err := config.GetTemplateTerraformDir("v1")
	if err != nil {
		return "", err
	}

	return templateDir, nil
}

func (o *Options) GetModulesDir() string {
	return canonicalModuleDestDir
}

func (o *Options) GetHomeDir() string {
	if o.homeDir == "" {
		o.homeDir, _ = os.UserHomeDir()
	}

	return o.homeDir
}

func (o *Options) GetCurrentDir() string {
	if o.currentDir == "" {
		o.currentDir, _ = os.Getwd()
	}

	return o.currentDir
}

func (o *Options) IsGitRepo() (string, bool) {
	currentDir := o.GetCurrentDir()
	repoRoot, _, err := utils.IsAGitRepository(currentDir, 3)
	if err != nil {
		return "", false
	}

	if o.repoRoot == "" {
		o.repoRoot = repoRoot
	}

	return repoRoot, true
}

func (o *Options) GetRepoRoot() (string, error) {
	if o.repoRoot == "" {
		repoRoot, isGitRepo := o.IsGitRepo()
		if isGitRepo {
			o.repoRoot = repoRoot
		}
	}

	return o.repoRoot, nil
}

type Client struct {
	Config Config
	Log    logger.Log
}

func NewGeneratorDefault() *Client {
	return &Client{
		Config: &Options{
			destDir:           canonicalModuleDestDir,
			tfTemplateVersion: "v1",
		},
		Log: logger.NewLogger().Logger,
	}
}
