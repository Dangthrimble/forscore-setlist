# Installing forscore-setlist from a release package

These packages contain two command-line executables. You do not need Go to run
them.

## Download and extract

Prefer a numbered release from the repository’s
[Releases page](https://github.com/Dangthrimble/forscore-setlist/releases).
Download the archive for your computer, for example
`forscore-setlist-v0.1.2-macos-apple-silicon.zip` or
`forscore-setlist-v0.1.2-windows-x64.zip`.

Extract the archive and keep the extracted folder together. It contains:

- the `4ss2open4ss` executable (`4ss2open4ss.exe` on Windows)
- the `txt2open4ss` executable (`txt2open4ss.exe` on Windows)
- [README.md](README.md)
- this guide
- the `docs/` guides
- `BACKLOG.md` and `DEVELOPMENT_HISTORY.md`
- `LICENSE`, `THIRD_PARTY_NOTICES.txt` and `THIRD_PARTY_LICENSES/`
- `VERSION.txt` and `SHA256SUMS.txt`

## Check the download

Each archive has a `.zip.sha256` file, and the release provides a combined
`SHA256SUMS.txt`. Compare the checksum before extracting:

```sh
shasum -a 256 forscore-setlist-v0.1.2-macos-apple-silicon.zip
```

On Linux, use `sha256sum` with the downloaded filename. In PowerShell, use:

```powershell
Get-FileHash -Algorithm SHA256 -LiteralPath 'forscore-setlist-v0.1.2-windows-x64.zip'
```

## Run the programs

Open a terminal, change to the folder that holds the file you want to convert,
and invoke each executable by its full path (or add the extracted folder to your
PATH). For example:

```sh
/path/to/extracted/4ss2open4ss <file.4ss>
/path/to/extracted/txt2open4ss <file.txt>
```

On Windows, use `4ss2open4ss.exe` and `txt2open4ss.exe` the same way from Command
Prompt or PowerShell.

Usage details are in [docs/usage.md](docs/usage.md). For install with Go instead,
see [docs/install-go.md](docs/install-go.md).
