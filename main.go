package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

var (
	token *string = flag.String("token", os.Getenv("GITLAB_TOKEN"), "api token with admin access")
)

func main() {
	flag.Parse()

	git := gitlab.NewClient(nil, *token)

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

	visibility := gitlab.PrivateVisibility
	options := gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{Page: 0, PerPage: 10}, Visibility: &visibility}
	projects, _, err := git.Projects.ListProjects(&options, nil)

	for _, project := range projects {
		p := *project
		fmt.Printf("clearing runners for project %s\n", p.Name)
		for _, r := range inactiveRunners {
			_, err := disableRunner(git, r, project)
			if err != nil {
				fmt.Printf("failed to disable runner %s due to %s\n", r.Description, err)
			}
		}
	}

	for _, r := range inactiveRunners {
		_, err := deleteRunner(git, r)
		if err != nil {
			fmt.Printf("failed to delete runner %s due to %s\n", r.Description, err)
		}
	}

	fmt.Println("total runners len:", len(runners))
	fmt.Println("inactive runners len:", len(inactiveRunners))

}

func disableRunner(c *gitlab.Client, r *gitlab.Runner, p *gitlab.Project) (bool, error) {
	// /projects/:id/runners/:runner_id
	uri := fmt.Sprintf("/projects/%d/runners/%d", p.ID, r.ID)
	req, err := c.NewRequest("DELETE", uri, nil, nil)
	var rs interface{}
	rsp, err := c.Do(req, &rs)
	if err != nil {
		return false, err
	}
	log.Printf("Got rsp %d with msg: %s for project %s and runner %s\n", rsp.StatusCode, rsp.Status, p.Name, r.Description)
	return rsp.StatusCode == 200, nil
}

func deleteRunner(c *gitlab.Client, r *gitlab.Runner) (bool, error) {
	id := fmt.Sprintf("runners/%d", r.ID)
	req, err := c.NewRequest("DELETE", id, nil, nil)
	var rs interface{}
	rsp, err := c.Do(req, &rs)
	if err != nil {
		return false, err
	}

	return rsp.StatusCode == 200, nil
}
