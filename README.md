# forscore-setlist

Tools that produce forScore **Open Setlist** `.4ss` files (XML body) for **forScore 12+**.

| Command | Purpose |
|---------|---------|
| `4ss2open4ss` | Native `.4ss` setlist (gzip / binary plist) → Open Setlist |
| `txt2open4ss` | Text file (one title per non-blank line) → Open Setlist placeholders |

If you only want the command-line tools, you do **not** need to download this
project’s files. Follow [Install with Go](docs/install-go.md); `go install …@latest`
fetches what it needs without a local copy of the repository.

## Install

- **With Go:** [Install with Go](docs/install-go.md) — Git, Go, `go install` for
  both commands, PATH, and optional project files.
- **Without Go:** [INSTALL.md](INSTALL.md) — install from a pre-built release
  package.

## Try it

```bash
4ss2open4ss <file.4ss>
txt2open4ss <file.txt>
```

Each writes a sibling file named `* (forScore12plus).4ss`. Full usage is in
[Using the tools](docs/usage.md).

## Docs

- [Using the tools](docs/usage.md) — commands, exit codes, element mapping
- [Install with Go](docs/install-go.md) — Git, Go, `go install`, PATH, update, optional clone
- [INSTALL.md](INSTALL.md) — install from a release package (no Go required)
- [Developing](docs/developing.md) — layout, local builds, tests
- [Backlog](BACKLOG.md) — outstanding work
- [Development history](DEVELOPMENT_HISTORY.md) — completed work and decisions

## Licence

[MIT](LICENSE)
