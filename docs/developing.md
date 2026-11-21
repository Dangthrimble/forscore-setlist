# Developing

Back to [README](../README.md).

## Layout

| Path | Role |
|------|------|
| [`cmd/4ss2open4ss`](../cmd/4ss2open4ss) | Native `.4ss` → Open Setlist |
| [`cmd/txt2open4ss`](../cmd/txt2open4ss) | Text lines → Open Setlist placeholders |
| [`cmd/buildpackage`](../cmd/buildpackage) | Release archive builder |
| [`internal/cli`](../internal/cli) | Shared version, usage and single-file argument handling |
| [`internal/convert`](../internal/convert) | Conversion library |
| [`release.json`](../release.json) | Product, module, command and extra-file configuration for packaging |

## Local builds

From a local checkout of this repository:

```bash
go build -o 4ss2open4ss ./cmd/4ss2open4ss
go build -o txt2open4ss ./cmd/txt2open4ss
```

## Tests

```bash
go test ./...
```

Fixture: [`internal/convert/testdata`](../internal/convert/testdata).

## Packaging

`cmd/buildpackage` builds the release archive. It verifies that each executable
was built with `CGO_ENABLED=0`, checks its `--version` output, adds the
documentation, dependency notices and an internal checksum list, and writes the
archive plus a `.zip.sha256` file.

```bash
go run ./cmd/buildpackage \
  -binary dist/bin/4ss2open4ss \
  -binary dist/bin/txt2open4ss \
  -version v0.1.2 -out dist/packages
```

Archives are named `forscore-setlist-v<version>-<platform>.zip` for all six
targets.

## CI and releases

[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) runs tests natively on
all six release targets (macOS Intel/Apple Silicon, Windows x64/ARM64 and Linux
x64/ARM64). The race detector runs on every target except Windows ARM64, where
Go does not support it. The workflow also runs `go vet`, `gofmt`, staticcheck and
govulncheck, and builds both commands. A final `CI gate` job aggregates the
required checks.

[`.github/workflows/release.yml`](../.github/workflows/release.yml) runs when a
`v*` tag is pushed. It calls the CI workflow with the tag as the package version,
then publishes the six `forscore-setlist-v<version>-<platform>.zip` archives and
a combined `SHA256SUMS.txt` only after every package has verified. The release
tags must match `vMAJOR.MINOR.PATCH`.

The release binaries stamp `--version` from the tag. `go install …@latest`
remains the documented install path; those binaries report the module version
from the install. A local `go build` reports `dev`.
