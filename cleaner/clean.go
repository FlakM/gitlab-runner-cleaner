package cleaner

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
	gitlab "github.com/xanzy/go-gitlab"
)

func CleanCommand() cli.Command {
	return cli.Command{
		Name:  "clean",
		Usage: "cleans offline runners for all available projects",
		Action: func(ctx *cli.Context) (err error) {
			git := gitlab.NewClient(nil, Config.GitlabToken)

			projects, err := ListAllProjects(git)
			if err != nil {
				return
			}

			runners, err := ListAllInactiveRunners(git)

			if err != nil {
				return
			}

			c1 := make(chan string, len(projects))

			for _, p := range projects {
				go UnregisterRunners(git, p, runners, c1)
			}

			for _ = range projects {
				log.Printf("%s", <-c1)
			}

			c2 := make(chan string, len(runners))
			for _, r := range runners {
				go RemoveRunner(git, r, c2)
			}

			for _ = range runners {
				log.Printf("%s", <-c2)
			}

			log.Println("Finished!")

			return
		},
	}

}

func UnregisterRunners(git *gitlab.Client, p *gitlab.Project, runners []*gitlab.Runner, c chan string) {
	counter := 0
	for _, r := range runners {
		_, err := git.Runners.DisableProjectRunner(p.ID, r.ID, nil)
		if err != nil {
			log.Printf("failed to disable runner %s due to %s\n", r.Description, err)
		} else {
			counter = counter + 1
		}
	}
	c <- fmt.Sprintf("removed %d runners for project %s", counter, p.Name)
}

func RemoveRunner(git *gitlab.Client, r *gitlab.Runner, c chan string) {

	_, err := git.Runners.RemoveRunner(r.ID, nil)
	if err != nil {
		log.Fatalf("failed to remove runner %s\n", r.Description)
	}

	c <- fmt.Sprintf("removed runner %s", r.Description)
}
