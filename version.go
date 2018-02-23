package main

import (
	"fmt"

	"github.com/urfave/cli"
)

// NB: Please, DO NOT set these manually!
// These variables are meant to be set through ldflags.
var buildTag, buildDate string

// VersionCommand returns a version command.
func VersionCommand() cli.Command {
	return cli.Command{
		Name:  "version",
		Usage: "Print the version number of the bag service",
		Action: func(ctx *cli.Context) {
			if buildTag != "" && buildDate != "" {
				fmt.Printf("version: \t\t%s\n", buildTag)
				fmt.Printf("build date: \t\t%s\n", buildDate)
				return
			}
			fmt.Print("undefined")
		},
	}
}
