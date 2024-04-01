package cmd

import (
	"github.com/Excoriate/tfgenctl/cli/internal/cliutils"
	"github.com/Excoriate/tfgenctl/cli/internal/generate"
	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}

type ModuleCmd struct {
	Name         string `arg:"" required:"" help:"The name of the module to operate on."`
	Group        string `name:"group" help:"The group name of the module. If the group is passed, normally it refers to the 'service' (cloud provider) the module is for. E.g.: /modules/s3/s3_bucket."`
	Path         string `name:"path" help:"The directory where the new module will be created. If not set, it'll be created in the current directory where the modules/ resides."`
	WithExamples bool   `name:"with-examples" help:"Generate examples for the module."`
	WithTests    bool   `name:"with-tests" help:"Generate tests for the module."`
}

func (c *ModuleCmd) Run(ctx *kong.Context) error {
	g := generate.NewGeneratorDefault()

	if err := g.Validate(c.Name, "",
		&generate.GenOptions{
			FailIfNotGitRepo: true,
		}); err != nil {
		return err
	}

	// Form the path or the destination directory for the module
	paths, err := g.ResolvePaths(c.Name, c.Group)
	if err != nil {
		return err
	}

	repoRoot, _ := g.Config.GetRepoRoot()

	if err := g.CreateBaseDirsIfNotExist(repoRoot); err != nil {
		return err
	}

	if err := g.GenerateModule("v1", paths); err != nil {
		return err
	}

	cliutils.PrintSuccess("Module generated!", paths.ModulePath)

	return nil
}
