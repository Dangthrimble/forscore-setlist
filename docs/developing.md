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
