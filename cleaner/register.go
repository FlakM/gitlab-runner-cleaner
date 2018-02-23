package cleaner

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"github.com/xanzy/go-gitlab"
)

func RegisterCommand() cli.Command {

	return cli.Command{
		Name:  "register",
		Usage: "register private runners by id or name",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "projects", Usage: "project's id to register to separated by comma"},
			cli.BoolTFlag{Name: "byId", Usage: "should register using id (otherwise will use with name lookup))"},
		},
		Action: func(ctx *cli.Context) (err error) {
			if len(ctx.Args()) != 1 {
				return errors.New("register takes only one argument: [name or ID]")
			}
			projects := strings.Split(ctx.String("projects"), ",")

			if len(projects) == 0 {
				return fmt.Errorf("you didn't specify any projects use -projects flag")
			}

			name := ctx.Args().Get(0)

			git := gitlab.NewClient(nil, Config.GitlabToken)

			var runnerID int

			if ctx.BoolT("byId") {
				runnerID, err = strconv.Atoi(name)
				if err != nil {
					log.Printf("maybe you wanted to use by name call and not id (int!)")
					return err
				}
			} else {
				r, err := GetRunnerByName(git, name)
				if err != nil { // error fetching runners
					return err
				}
				if r == nil {
					return fmt.Errorf("runner %s not found", name)
				}
				runnerID = r.ID
			}

			for _, proj := range projects {
				runner, rsp, err := git.Runners.EnableProjectRunner(proj, &gitlab.EnableProjectRunnerOptions{RunnerID: runnerID}, nil)
				if err != nil && rsp.StatusCode != 409 {
					return err
				}

				if rsp.StatusCode != 409 {
					log.Printf("Runner {id: %d, name: %s} enabled for project %s\n", runner.ID, runner.Description, proj)
				} else {
					log.Printf("Runner {id: %d} was already enabled for project %s\n", runnerID, proj)
				}

			}

			return nil
		},
	}

}
