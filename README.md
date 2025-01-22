[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/architect/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/architect/tree/main)

# architect

A highly opinionated tool used at Giant Swarm for building services.

Architect is used as part of the Giant Swarm release workflow, to **build services**.
The latest release is fetched automatically during builds (running on CircleCI),
and then executed to perform the build. This allows us to update one tool,
and affect all builds.

architect runs all build steps in containers, to allow for portability and reproducibility.
