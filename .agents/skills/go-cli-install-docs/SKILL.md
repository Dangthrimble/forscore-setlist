---
name: go-cli-install-docs
description: >-
  Installs Go command-line tools with go install and writes README install/update
  docs for copy-paste users (Git, repository, Go, other dependencies, PATH,
  verification). Use when building or installing Go CLIs under cmd/, revising
  install instructions, or when the user mentions go install, GOBIN, GOPATH/bin,
  PATH setup, Git install, Go install/update, dependencies, or cloning/updating
  the repository for a Go utility.
---

# Go CLI install (this repository)

## Policy

- Prefer **Go + `go install`** over release binaries or system package installs.
- Do not suggest unsigned downloads, Gatekeeper workarounds, or copying into
  `/usr/bin` with sudo.
- Audience: can copy-paste Terminal commands; explain PATH explicitly.
- Docs: README must cover Git → (optional) repository → Go → other dependencies →
  Install → Verify → PATH fix → Update.
- **Up front:** the README must state clearly that users who only want the
  command-line tools do **not** need to download or clone the project files;
  after Git (as documented) and Go are ready, `go install …@latest` is enough.
  Put that callout before the Git/repository sections so readers see it first.
- There needs to be a link to or description of how to install and upgrade Git and
  the repository. The description, whether included or linked to, needs to use
  non-technical plain English. Document the repository as **optional** (for people
  who want the source, or local `go install ./cmd/…` / `go build`)—not as a
  prerequisite for installing the CLIs.
- The same plain-English rule applies to **Go** (install and update) and to **any
  other dependencies** required to install or update the command-line utilities.
  When a new dependency appears, add or link a non-technical install-and-update
  description for it before documenting the CLI install steps that need it.

## Discover the repo first

1. Read `go.mod` → module path + minimum Go version.
2. List packages under `cmd/` → executable names (directory names).
3. Find the public clone URL (`git remote get-url origin` or the hosting page).
4. Note every prerequisite needed to install or update the CLIs (Git, Go, and any
   others that exist now or are added later).
5. Install/doc commands use:
   `go install <module>/cmd/<name>@latest`
   For a local checkout: `go install ./cmd/<name>` from the module root.

If there is no `cmd/`, ask which main packages to install.

## Local install

1. Ensure Git is installed (`git --version`); if not, follow the README’s plain-English
   Git section or link.
2. Do **not** require a clone for CLI-only install. Clone or update the repository
   only when the user wants the source files or a local checkout (see README);
   end-user docs must say this up front.
3. Ensure Go meets `go.mod` (`go version`); if not, follow the README’s plain-English
   Go install/update section or link, then open a new terminal.
4. Ensure any other documented dependencies are installed/updated the same way.
5. `go install` each `cmd/...` as above (`…@latest` from the module path; no local
   checkout required).
6. Binaries: `$(go env GOBIN)` or `$(go env GOPATH)/bin`.
7. Verify: `command -v <name>` (Windows: `where <name>`).
8. If missing → PATH section in the template; new terminal; recheck.

## README updates

Fill [readme-install-template.md](readme-install-template.md) using this repo’s
module path, `cmd` names, Go version, clone URL, and dependency list.

Required doc pieces (each install **and** update path in non-technical plain
English, inline or linked):

1. **Upfront callout** — CLI-only users can skip downloading the project; go to
   **Go**, then **Install**.
2. **Git** — install and upgrade.
3. **Repository (optional)** — first-time download/clone and later upgrade
   (`git pull`); for source checkout or developers, not required to install the
   CLIs.
4. **Go** — what Go is in one plain sentence; install the version from `go.mod`;
   how to update Go later; confirm with `go version` after a new terminal.
5. **Other dependencies** — for every extra tool or runtime needed to install or
   update the CLIs, the same install-and-update treatment. If there are none yet,
   omit the section; when one is added, document it before the CLI install steps.
6. **Install the CLIs** — `go install …@latest` with full module paths.
7. **Verify** — `command -v` / `where` and one example run.
8. **PATH fix** — if not found: print bin dir, exact zsh/bash (and short Windows) steps,
   new terminal, recheck.
9. **Update the CLIs** — upgrade Git/Go/other dependencies as needed, update the
   repo copy if used, then re-run the same `go install …@latest` commands.

Official download pages may be linked; surrounding README text must stay plain English.
Do not lead with `go build` for end users.

## Out of scope

Release binaries, notarization, Homebrew bottles—unless the user explicitly
asks to add them.
