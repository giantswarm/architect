# Integrating app-build-suite (abs) into the architect image

## Goal

Bundle `app-build-suite` (abs) directly into the `architect` image so that:

- The `architect-orb` `push-to-app-catalog` job can run on the **`architect`
  executor** instead of the dedicated `app-build-suite` executor.
- We stop maintaining the extra **`circleci.Dockerfile`** in the
  `app-build-suite` repo (it only layered CI tools — conftest, gh, cosign,
  gitsemver, etc. — that the architect image already ships).
- The redundant `src/executors/app-build-suite.yaml` in `architect-orb` can be
  deleted.

Historically abs shipped as a Docker image (`app-build-suite:<ver>` +
`app-build-suite:<ver>-circleci`). abs is now installable via `uv`
(`uv tool install app-build-suite`, invoked as `abs`), which makes bundling it
into the architect image practical.

## Decisions (locked with the requester)

- **Target the existing Alpine/musl architect image** (single-image goal), but
  **validate a real chart build before relying on it** (see Validation).
- **Install abs from PyPI at the latest stable release** (`uv tool install
  app-build-suite`, unpinned). Pinning is a possible later hardening.

## What changed in this repo

`Dockerfile`:

- Added a `uv` build stage (`ghcr.io/astral-sh/uv:0.11.19`, matching abs's own
  Dockerfile).
- Added `cairo`, `pango`, `font-dejavu` to the `apk` install — abs uses
  `cairosvg`/`pillow` for chart **icon (SVG) validation** (the abs image
  installs the Debian equivalent `libpangocairo-1.0-0`).
- Added `kube-linter` (`v0.8.3`, matching abs) — an abs runtime dependency for
  chart linting. Arch-aware download (`kube-linter-linux.tar.gz` for amd64,
  `kube-linter-linux_arm64.tar.gz` for arm64).
- Installed abs via uv onto a **managed Python 3.13**:
  `uv tool install --python 3.13 app-build-suite`. The system Python is 3.11,
  which is too old (abs requires `>=3.13`), so a uv-managed 3.13 interpreter is
  used. The `abs` entry point is linked into `/usr/local/bin`; the managed
  interpreter and tool venv live under `/opt/uv`. The build runs `abs --version`
  as a smoke test.

`CHANGELOG.md`: added an `Added` entry under `[Unreleased]`.

What was **already present** (so it was not re-added): helm, chart-testing
(`ct`), conftest, cosign, gitsemver, gh / gh-token, jq, openssh-client.

The existing `gs_metadata_chart_schema.yaml` copy from the `app-build-suite:2.1.2`
build stage (`/etc/ct/chart_schema.yaml`) was **left as-is** — it is the regular
abs image, not the `-circleci` variant, and extracting the schema reliably from
the PyPI wheel is unconfirmed. See Open questions.

## Validation — this is the go/no-go for staying on Alpine/musl

The Dockerfile edits cannot be proven correct without an actual multi-arch
build. abs and its native deps are upstream-validated on **Debian/glibc +
Python 3.13**; this image is **Alpine/musl**.

Build and test **both arches**:

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t architect-abs-test .
```

Then, in a container, run a real chart package (not just `abs --version`),
mirroring what `architect-orb`'s `push-to-app-catalog` does:

```bash
abs --chart-dir ./helm/<chart> \
    --destination build \
    --generate-metadata \
    --catalog-base-url "https://giantswarm.github.io/<catalog>/" \
    --keep-chart-changes \
    --override-chart-version 0.0.0 --override-app-version 0.0.0
```

Watch for:

- uv being able to download a **musl** managed Python 3.13 for both arches.
- **musllinux cp313 wheels** resolving for `pillow` and `cffi` (used by
  `cairocffi`/`cairosvg`). If they don't, the build will fail and we'd need to
  add build deps (`build-base`, `python3-dev`, `libffi-dev`, `cairo-dev`,
  jpeg/zlib/freetype dev libs) — or fall back to a glibc image for this tool.
- `cairosvg` being able to `dlopen` `libcairo`/`libpango` at runtime (the
  reason for the `cairo`/`pango` apk packages).

If musl proves too costly, the fallback is a **glibc** image variant carrying
abs, with only the `push-to-app-catalog` job's executor pointing at it.

## Follow-up in architect-orb (separate repo, after a new image tag ships)

Once this image is built and a new tag is published, in `architect-orb`:

1. `src/executors/architect.yaml` — bump the image tag to the abs-bundled one.
2. `src/jobs/push-to-app-catalog.yaml` — change `executor: "<< parameters.executor >>"`
   to `executor: "architect"`; change the "Execute App Build Suite" step from
   `python -m app_build_suite ...` to `abs ...`. Keep the (already deprecated)
   `executor` enum param as an accepted no-op for backward compatibility.
3. `src/commands/tools-info.yaml` — change `python -m app_build_suite --version`
   to `abs --version`.
4. Delete `src/executors/app-build-suite.yaml`.
5. Refresh the stale comment in `src/commands/push-helm.yaml` that says
   "executors without a docker CLI (e.g. app-build-suite)" — the architect image
   has docker; the manual `~/.docker/config.json` write stays (harmless).
6. Update `docs/job/push-to-app-catalog.md` + `CHANGELOG.md`.
7. Dev-publish the orb (`@dev:<branch>`) and run `push-to-app-catalog`
   end-to-end against a consuming repo (package + metadata + push + cosign).

Note: `run-tests-with-ats` is **out of scope** — it uses **app-test-suite** (a
different tool, `machine` executor), not abs.

## Final cleanup (in app-build-suite repo, after rollout)

- Remove `circleci.Dockerfile` and its image build/publish.
- Optionally stop publishing `*-circleci` abs image tags.

## Open questions

- Can `gs_metadata_chart_schema.yaml` be sourced from the installed PyPI wheel
  (dropping the `app-build-suite:2.1.2` build stage entirely)? Needs confirming
  whether abs's `resources/` dir ships in the wheel.
- Should the PyPI abs version be pinned for reproducible image builds (currently
  unpinned per the "latest stable" decision)?
