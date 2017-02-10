# architect

`architect` is a tool for managing Giant Swarm release engineering

To fetch the latest build of `architect`, run:
```
wget -q $(curl -sS https://api.github.com/repos/giantswarm/architect/releases/latest | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
```