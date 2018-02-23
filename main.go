package main

import (
	"log"
	"os"

	"github.com/FlakM/gitlab-runner-cleaner/cleaner"
	"github.com/urfave/cli"
)

var (
	App *cli.App
)

func init() {
	App = cli.NewApp()

	App.Name = "gitlab cleaner"
	App.Usage = `Check -h`
	App.Author = "Maciej Flak"

	App.Version = ""

	cleaner.InitializeConfig()
}

func main() {
	AddCommands()
	if err := App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func AddCommands() {
	AddCommand(cleaner.CleanCommand())
	AddCommand(cleaner.RegisterCommand())
}

// AddCommand adds a child command.
func AddCommand(cmd cli.Command) {
	App.Commands = append(App.Commands, cmd)
}
