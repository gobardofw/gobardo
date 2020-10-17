package main

import (
	"github.com/gobardofw/cli"
	"github.com/gobardofw/gobardo/commands"
)

func main() {
	cli := cli.NewCLI("gobardo", "gobardo framework cli tools")
	cli.AddCommand(commands.VersionCommand)
	cli.AddCommand(commands.NewCommand)
	cli.Run()
}
