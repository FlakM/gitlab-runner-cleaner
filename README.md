# Description

This is a simple tool to manage runners for gitlab projects.

We use ephemeral docker auto registering instances of gitlab-runner [image here](https://github.com/FlakM/gitlab-runner-auto-register).
As runner is able to register itself into one project we added some environment variables that make this container register itself in given projects.

If the docker engine dies from some reason before runner manages to unregister it becomes offline. 
To clean up from time to time you have to unregister runners from all projects and then remove it.

# Usage



```bash
go run main.go 
```
