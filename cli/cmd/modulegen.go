package cmd

import (
	"github.com/Excoriate/tfgenctl/cli/internal/cliutils"
	"github.com/Excoriate/tfgenctl/cli/internal/generate"
	"github.com/alecthomas/kong"
)

type Context struct {
	Debug bool
}

type TfCmd struct {
	ModuleCmd ModuleCmd `cmd:"" name:"module" help:"🚀 Generate a terraform module. Here's how you can use it:\n\nUse this command to scaffold a new Terraform module with ease. You can specify the path where the module should be created and optionally include examples to kickstart your development.\n\nExample:\n  tfgenctl module my-awesome-module --path=./modules --with-examples\n\nThis will create a new module named 'my-awesome-module' in the './modules' directory and include example usage of the module.\n\nFlags:\n  --path          The directory where the new module will be created. This is required.\n  --with-examples Include a set of example files demonstrating how to use the module. Optional but highly recommended to get started quickly.\n\n💡 Tip: Use the --debug flag if you're encountering issues to print detailed log messages during execution."`
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

	if c.WithExamples {
		if err := g.GenerateModuleExample("v1", paths); err != nil {
			return err
		}
	}

	cliutils.PrintSuccess("Module generated!", paths.ModulePath)

	return nil
}
