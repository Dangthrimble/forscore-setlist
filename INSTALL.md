# Installing forscore-setlist from a release package

These packages contain two command-line executables. You do not need Go to run
them.

## Download and extract

Prefer a numbered release from the repository’s
[Releases page](https://github.com/Dangthrimble/forscore-setlist/releases).
Download the archive for your computer (for example
`forscore-setlist_darwin_arm64.tar.gz` or
`forscore-setlist_windows_amd64.zip`).

Extract the archive and keep the extracted folder together. It contains:

- the `4ss2open4ss` executable (`4ss2open4ss.exe` on Windows)
- the `txt2open4ss` executable (`txt2open4ss.exe` on Windows)

The repository [README](README.md) and guides are the reference documentation;
they are not included in the current release archives.

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
