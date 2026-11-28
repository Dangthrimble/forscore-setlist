# Using the tools

Back to [README](../README.md).

## Usage

```bash
4ss2open4ss <file.4ss>
txt2open4ss <file.txt>
4ss2open4ss --version
txt2open4ss --version
```

Use `--` before a filename that begins with a hyphen.

Each conversion writes a sibling file named `* (forScore12plus).4ss` (Open Setlist is only readable in forScore 12+), prints the output path, and exits. `--version` prints `<command> <version>` and exits `0`. `--help` and `-h` print usage and exit `0`.

Exits:

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

## Input limits

Both commands accept at most **100 entries**. In a text file, blank and
whitespace-only lines are ignored and do not count towards the limit.

The native `.4ss` path also applies these limits:

- compressed input: 2 MiB;
- decompressed plist: 4 MiB;
- object-graph depth: 64;
- individual title or path: 16 KiB.

A setlist with 101 entries is rejected rather than truncated. If you have a
legitimate setlist that needs more than 100 entries, open an issue with the
fixture at <https://github.com/Dangthrimble/forscore-setlist/issues>.
