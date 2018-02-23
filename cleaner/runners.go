package cleaner

import (
	"strconv"

	gitlab "github.com/xanzy/go-gitlab"
)

func GetRunnerByName(git *gitlab.Client, name string) (r *gitlab.Runner, err error) {
	return recursiveRunnerByDescription(git, 1, name)
}

func recursiveRunnerByDescription(git *gitlab.Client, page int, desc string) (r *gitlab.Runner, err error) {
	scope := gitlab.String("active")
	opt := gitlab.ListRunnersOptions{Scope: scope}
	opt.Page = page
	runners, rsp, err := git.Runners.ListRunners(&opt, nil)
	if err != nil {
		return
	}

	for _, runner := range runners {
		if runner.Description == desc {
			return runner, nil
		}
	}

	next := rsp.Header.Get("X-Next-Page")
	if next != "" {
		n, err := strconv.Atoi(next)
		if err != nil { // github returned some bullshit
			return nil, err
		}
		return recursiveRunnerByDescription(git, n, desc)
	}
	return
}

func ListAllInactiveRunners(git *gitlab.Client) (r []*gitlab.Runner, err error) {
	all, err := recursiveAllRunners(git, 1)

	if err != nil {
		return nil, err
	}

	inactiveRunners := []*gitlab.Runner{}

	for _, runner := range all {
		if !(*runner).Online {
			inactiveRunners = append(inactiveRunners, runner)
		}
	}
	return inactiveRunners, nil
}

func recursiveAllRunners(git *gitlab.Client, page int) (r []*gitlab.Runner, err error) {
	opt := gitlab.ListRunnersOptions{Scope: gitlab.String("active")}
	opt.Page = page
	runners, rsp, err := git.Runners.ListRunners(&opt, nil)
	if err != nil {
		return
	}

	next := rsp.Header.Get("X-Next-Page")
	if next != "" {
		n, err := strconv.Atoi(next)
		if err != nil { // github returned some bullshit
			return nil, err
		}
		rest, err := recursiveAllRunners(git, n)
		if err != nil {
			return nil, err
		}

		runners = append(runners, rest...)
	}
	return runners, nil
}

func ListAllProjects(git *gitlab.Client) (r []*gitlab.Project, err error) {
	return recursiveAllProjects(git, 1)
}

func recursiveAllProjects(git *gitlab.Client, page int) (r []*gitlab.Project, err error) {
	visibility := gitlab.PrivateVisibility
	options := gitlab.ListProjectsOptions{ListOptions: gitlab.ListOptions{Page: page}, Visibility: &visibility}
	projects, rsp, err := git.Projects.ListProjects(&options, nil)

	if err != nil {
		return
	}

	next := rsp.Header.Get("X-Next-Page")
	if next != "" {
		n, err := strconv.Atoi(next)
		if err != nil { // github returned some bullshit
			return nil, err
		}
		rest, err := recursiveAllProjects(git, n)
		if err != nil {
			return nil, err
		}

		projects = append(projects, rest...)
	}
	return projects, nil
}
