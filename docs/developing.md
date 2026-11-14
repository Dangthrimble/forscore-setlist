# Developing

Back to [README](../README.md).

## Layout

| Path | Role |
|------|------|
| [`cmd/4ss2open4ss`](../cmd/4ss2open4ss) | Native `.4ss` → Open Setlist |
| [`cmd/txt2open4ss`](../cmd/txt2open4ss) | Text lines → Open Setlist placeholders |
| [`internal/convert`](../internal/convert) | Conversion library |

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

## CI and releases

[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) runs `go test ./...`, `go vet ./...`, and builds both commands on Ubuntu, macOS, and Windows for every push and pull request.

[`.github/workflows/release.yml`](../.github/workflows/release.yml) runs when a `v*` tag is pushed. It publishes `4ss2open4ss` and `txt2open4ss` for darwin, linux, and windows (amd64 and arm64). macOS and Linux archives are `.tar.gz`; Windows archives are `.zip` with `.exe` binaries. `go install …@latest` remains the documented install path.
