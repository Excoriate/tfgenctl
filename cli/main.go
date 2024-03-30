package main

import (
	"github.com/Excoriate/tfgenctl/cli/cmd"
	"github.com/Excoriate/tfgenctl/cli/internal/cliutils"
	"github.com/alecthomas/kong"
)

func main() {
	cli := cmd.CLI{
		Globals: cmd.Globals{
			Version: cmd.VersionFlag("0.0.1"),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("tfgenctl"),
		kong.Description("A CLI tool to generate terraform modules."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
	err := ctx.Run(&cli.Globals)
	if err != nil {
		cliutils.PrintErrorFatal("Ups, something went wrong!", err)
	}
}
