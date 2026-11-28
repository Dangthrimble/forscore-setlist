# Development history

| When | What |
|------|------|
| 2026-10-17 | Initial checkpoint: converter CLI (`main.go`, `convert.go`), tests, Mount Pleasant Baptist fixture; module originally named `4ss2xml`. |
| 2026-10-17 | Module renamed to `4ss24ssxml`. |
| 2026-10-17 | Scrubbed `.DS_Store` from git history on `main`. |
| 2026-10-17 | Restructure to `cmd/4ss24ssxml` + `internal/convert`; module path `github.com/Dangthrimble/4ss24ssxml`; README, backlog, MIT licence; fixture hygiene. |
| 2026-10-17 | Renamed CLI to `4ss2open4ss`; added `txt2open4ss` (non-blank lines → placeholders); shared `(forScore12plus)` output marker. |
| 2026-11-09 | Added `.agents/skills/go-cli-install-docs` and rewrote README install/update docs (Git, repo, Go, PATH) for copy-paste users. |
| 2026-11-12 | GitHub repo and module path moved from `4ss24ssxml` to `forscore-setlist`. |
| 2026-11-13 | Split the README into a hub with linked guides; added the `maintain-before-commit` skill. |
| 2026-11-14 | Added CI and multi-platform release workflows; both commands report the Git tag from `--version`; released v0.1.1. |
| 2026-11-20 | Adopted the shared cross-agent guidance, UK English rule, and release and packaging conventions. |
| 2026-11-20 | Standardised CI with race tests, gofmt, vet, staticcheck, govulncheck, and a required `CI gate`; raised the minimum Go version to 1.26.11. |
| 2026-11-21 | Added the Go release packager, six native CI-built packages, versioned ZIP assets, per-asset and combined checksums, and a tag-validated release workflow. |
| 2026-11-21 | Standardised CLI behaviour: `<command> <version>` output, `--help`/`-h`, `--` termination, command-prefixed errors, and exit `2` for usage errors. |
| 2026-11-21 | Added the shared repository-settings script, Dependabot, issue and pull-request templates, and the security policy. |

## Design notes

- Native `.4ss` input is treated as gzip → binary plist → NSKeyedArchiver (`$objects` / `$top.setlist`).
- Already-XML / Open Setlist inputs are rejected.
- Output file name: input stem + ` (forScore12plus)` + `.4ss`.
- Dependency: [`howett.net/plist`](https://pkg.go.dev/howett.net/plist).
