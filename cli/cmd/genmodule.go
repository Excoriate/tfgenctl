package cmd

import "fmt"

type Context struct {
	Debug bool
}

type GenModuleCmd struct {
	Path           []string `arg:"" name:"path" help:"Path where the 'terraform module' is going to be generated." type:"path"`
	IncludeExample bool     `help:"Include example module."`
	Force          bool     `help:"Force overwriting existing files and/or creating canonical directories if they aren't present'."`
}

func (c *GenModuleCmd) Run(ctx *Context) error {
	fmt.Println("running", c.Path)
	return nil
}
