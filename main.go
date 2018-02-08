package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	token  = flag.String("token", os.Getenv("GITLAB_TOKEN"), "api token with admin access")
	dryRun = flag.Bool("dryRun", true, "if set to true only prints jobs")
)

func main() {
	flag.Parse()

	git := gitlab.NewClient(nil, *token)

	fmt.Printf("Dry run set to %t\n", *dryRun)

	runners, _, err := git.Runners.ListRunners(nil, nil)

	if err != nil {
		log.Fatal(err)
	}

	inactiveRunners := []*gitlab.Runner{}

	for _, runner := range runners {
		if !(*runner).Online {
			inactiveRunners = append(inactiveRunners, runner)
		}
	}

	fmt.Println("total runners", len(runners))
	fmt.Println("inactive runners ", len(inactiveRunners))

	visibility := gitlab.PrivateVisibility
	options := gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{Page: 0, PerPage: 10}, Visibility: &visibility}
	projects, _, err := git.Projects.ListProjects(&options, nil)

	for _, project := range projects {
		p := *project
		fmt.Printf("disabling runners for project %s\n", p.Name)
		for _, r := range inactiveRunners {
			if *dryRun {
				fmt.Printf("would disable runner %d for project %d\n", r.ID, p.ID)
			} else {
				git.Runners.DisableProjectRunner(p.ID, r.ID, nil)
				if err != nil {
					fmt.Printf("failed to disable runner %s due to %s\n", r.Description, err)
				} else {
					fmt.Printf("Runner %d disabled for project %d [%s]\n", r.ID, p.ID, p.Name)
				}
			}
		}
	}

	for _, r := range inactiveRunners {
		if *dryRun {
			fmt.Printf("would remove runner %d with description %s\n", r.ID, r.Description)
		} else {
			fmt.Printf("removing runner with id %d\n", r.ID)
			_, err := git.Runners.RemoveRunner(r.ID, nil)
			if err != nil {
				fmt.Printf("failed to delete runner %s due to %s\n", r.Description, err)
			} else {
				fmt.Printf("deleted runner %s\n", r.Description)
			}
		}
	}
}
