[![CircleCI](https://circleci.com/gh/giantswarm/architect.svg?style=shield)](https://circleci.com/gh/giantswarm/architect)

# architect

architect is a highly opinionated tool used at Giant Swarm for building and deploying services.

architect is used as part of the Giant Swarm release workflow, to **build services**.
The latest release is fetched automatically during builds (running on CircleCI),
and then executed to perform the build. This allows us to update one tool,
and affect all builds.

On master merges, architect is also used to **trigger a deployment** of the built
service. It creates an event that is picked up by a companion tool
[draughtsman](https://github.com/giantswarm/draughtsman), which runs inside an
installation and pulls and deploys the service.

architect runs all build steps in Docker containers, to allow for portability and reproducibility.

To download the latest build of architect, run:

```nohighlight
wget -q $(curl -sS https://api.github.com/repos/giantswarm/architect/releases/latest | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
```
