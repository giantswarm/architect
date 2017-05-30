# architect

`architect` is a tool for managing Giant Swarm release engineering

`architect` is used as part of the Giant Swarm release workflow, to build services.
The latest release is fetched automatically during builds (running on CircleCI),
and then executed to perform the build. This allows us to update one tool,
and affect all builds.

On master merges, `architect` is also used for the actual deployment.
This is likely to change in the future, as we cannot push from the build servers to all installations.

`architect` runs all build steps as Docker containers, to allow for portability and reproducibility.

To download the latest build of `architect`, run:
```
wget -q $(curl -sS https://api.github.com/repos/giantswarm/architect/releases/latest | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
```
