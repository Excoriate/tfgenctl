package cmd

import (
	"fmt"
	"github.com/alecthomas/kong"
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println("Version: ", vars["version"])
	return nil
}

type Globals struct {
	Version VersionFlag `name:"version" help:"Print version information and quit"`
	Debug   bool        `short:"D" help:"Enable debug mode"`
}

type CLI struct {
	Globals

	ModuleCmd ModuleCmd `cmd:"" name:"module" help:"🚀 Generate a terraform module. Here's how you can use it:\n\nUse this command to scaffold a new Terraform module with ease. You can specify the path where the module should be created and optionally include examples to kickstart your development.\n\nExample:\n  tfgenctl module my-awesome-module --path=./modules --with-examples\n\nThis will create a new module named 'my-awesome-module' in the './modules' directory and include example usage of the module.\n\nFlags:\n  --path          The directory where the new module will be created. This is required.\n  --with-examples Include a set of example files demonstrating how to use the module. Optional but highly recommended to get started quickly.\n\n💡 Tip: Use the --debug flag if you're encountering issues to print detailed log messages during execution."`
}
