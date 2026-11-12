# 4ss24ssxml

Tools that produce forScore **Open Setlist** `.4ss` files (XML body) for **forScore 12+**.

| Command | Purpose |
|---------|---------|
| `4ss2open4ss` | Native `.4ss` setlist (gzip / binary plist) → Open Setlist |
| `txt2open4ss` | Text file (one title per non-blank line) → Open Setlist placeholders |

If you only want the command-line tools, you do **not** need to download this project’s files. Skip **This project’s files**, set up **Git** and **Go** below, then go to **Install**.

## Git

Git is a tool that downloads and updates project files on your computer.

- **Install:** get Git from [git-scm.com/downloads](https://git-scm.com/downloads) (choose the installer for your computer) and run it.
- **Upgrade:** download a newer installer from the same page and run it again, or use your computer’s usual update method.

Open a **new** Terminal window and check:

```bash
git --version
```

## This project’s files (repository)

Optional. Only needed if you want a copy of the source on your computer (or to build from a local checkout). Not required to install the commands.

- **First time:** download a full copy of the project. In Terminal:

  ```bash
  git clone https://github.com/Dangthrimble/forscore-setlist.git
  cd forscore-setlist
  ```

  This creates a folder named `forscore-setlist` with the project files inside.

- **Later:** get the latest project files:

  ```bash
  cd forscore-setlist
  git pull
  ```

Updating these files does **not** by itself update the installed commands. To refresh the commands, follow **Update** below.

## Go

Go is the programming toolkit used to download and install these command-line programs.

- **Install:** get Go **1.24** or newer from [go.dev/dl](https://go.dev/dl/). Choose the installer for your computer, run it, and accept the defaults unless you know you need otherwise.
- **Check:** open a **new** Terminal window and run:

  ```bash
  go version
  ```

  You should see a version of 1.24 or higher.

- **Update:** later, download a newer installer from the same page and run it again (or follow that page’s update steps), then open a new Terminal and run `go version` to confirm.

## Install

In Terminal, run (no local copy of the project is required):

```bash
go install github.com/Dangthrimble/4ss24ssxml/cmd/4ss2open4ss@latest
go install github.com/Dangthrimble/4ss24ssxml/cmd/txt2open4ss@latest
```

## Check that it worked

```bash
command -v 4ss2open4ss
command -v txt2open4ss
```

Each command should print a path (often ending in `go/bin/…`).

Then try:

```bash
4ss2open4ss
```

If you see usage text (or a clear usage error), the program is installed. The same idea applies to `txt2open4ss`.

## If `command -v` prints nothing (PATH)

Go puts the programs in a `bin` folder that your computer may not search yet.

1. Show the folder:

   ```bash
   go env GOBIN
   go env GOPATH
   ```

   Use `GOBIN` if it printed a path; otherwise use `GOPATH` plus `/bin` (example: `/Users/you/go/bin`).

2. **zsh** (default on Mac) — add this line to `~/.zshrc`:

   ```bash
   export PATH="$(go env GOPATH)/bin:$PATH"
   ```

   If `go env GOBIN` printed a path, use that instead of `$(go env GOPATH)/bin`.

3. Save the file. Open a **new** Terminal window (or run `source ~/.zshrc`).

4. Run `command -v` again.

**bash:** put the same `export PATH=…` line in `~/.bashrc` or `~/.bash_profile`.

**Windows:** add the folder from `go env GOPATH` + `\bin` (or `GOBIN`) to your user Path under Environment Variables, open a new Command Prompt, and run `where 4ss2open4ss`.

## Update

1. Upgrade Git if needed (see **Git**).
2. Update Go if needed (see **Go**).
3. If you keep a local copy of the project, run `git pull` (see **This project’s files**).
4. Re-run the same `go install …@latest` commands as in **Install**.
5. Run `command -v` again.

## Usage

```bash
4ss2open4ss <file.4ss>
txt2open4ss <file.txt>
```

Each writes a sibling file named `* (forScore12plus).4ss` (Open Setlist is only readable in forScore 12+), prints the output path, and exits:

- `2` — bad usage
- `1` — read / convert / write error
- `0` — success

### `4ss2open4ss`

Setlist items are classified as:

| Element | Source cues |
|---------|-------------|
| `score` | PDF `FilePath` |
| `bookmark` | `Bookmark` / type fields |
| `placeholder` | `Placeholder` |

### `txt2open4ss`

Blank and whitespace-only lines are skipped. Each non-blank line becomes a `<placeholder title="..." />`.

## Layout

| Path | Role |
|------|------|
| [`cmd/4ss2open4ss`](cmd/4ss2open4ss) | Native `.4ss` → Open Setlist |
| [`cmd/txt2open4ss`](cmd/txt2open4ss) | Text lines → Open Setlist placeholders |
| [`internal/convert`](internal/convert) | Conversion library |

### For developers

From a local checkout of this repository:

```bash
go build -o 4ss2open4ss ./cmd/4ss2open4ss
go build -o txt2open4ss ./cmd/txt2open4ss
```

## Tests

```bash
go test ./...
```

Fixture: [`internal/convert/testdata`](internal/convert/testdata).

## License

[MIT](LICENSE)

## Related docs

- [DEVELOPMENT_HISTORY.md](DEVELOPMENT_HISTORY.md)
- [BACKLOG.md](BACKLOG.md)
