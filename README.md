# Description

This is a simple tool to manage runners for gitlab projects.

We use ephemeral docker auto registering instances of gitlab-runner [image here](https://github.com/FlakM/gitlab-runner-auto-register).
As runner is able to register itself into one project we added some environment variables that make this container register itself in given projects.

If the docker engine dies from some reason before runner manages to unregister it becomes offline. 
To clean up from time to time you have to unregister runners from all projects and then remove it.

# Usage


Basic usage: 

```bash
# https://github.com/FlakM/gitlab-runner-cleaner/releases go here to check latest releases:
VER=0.0.1
GITLAB_TOKEN=<my secret token here>
curl -LJO https://github.com/FlakM/gitlab-runner-cleaner/releases/download/$VER/gitlab-runner-cleaner.sh
chmod +x gitlab-runner-cleaner.sh
./gitlab-runner-cleaner.sh -h # this requires GITLAB_TOKEN enviroment variable
```


