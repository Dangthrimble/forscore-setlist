# Using the tools

Back to [README](../README.md).

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
