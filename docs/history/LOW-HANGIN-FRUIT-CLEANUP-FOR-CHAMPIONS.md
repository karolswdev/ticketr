# Low-Hanging Fruit Cleanup for Champions

## Why This Exists
Ticketr’s repo needs to feel as polished as the TUI experience we’re selling. This brief captures obvious transitory artifacts that erode trust, leak sensitive data, or bloat the tree. Treat it as the intake for a rapid cleanup sprint before we hand the keys to new contributors.

## Highest-Risk Offenders (Nuke First)
- `tmp.delete.me.scratchpad` *(root)* — contains a live Jira API token and workspace command transcript. **P0: delete immediately and rotate that token.**
- `pulled_tickets.md`, `tickets.md` *(root)* — real Jira ticket dumps. Remove and purge from history if ever committed.

## Binary & Build Artifacts
- `ticketr`, `ticketr-test` *(root, ~9–10 MB each)* — compiled binaries; replace with `make clean` or `go clean` target. Ensure `.gitignore` (already present) stays enforced and add a pre-commit hook to block them.
- `.ticketr/`, `.ticketr.state` *(root)* — runtime state/logs. These are ignored but sitting in-tree; script a cleanup (`rm -rf .ticketr .ticketr.state`) and add guidance so contributors keep state under XDG paths only.

## Test & Experiment Debris
- `demo-search.go` — ad-hoc demo program; either move into `examples/` with docs or delete once the feature is tested.
- `coverage.out` — leftover from `go test -cover`; add to CI artifacts or include in `.gitignore` (pattern already there) and sweep existing file.
- `tests/`, `evidence/`, legacy weekly reports — confirm whether they’re still serving release artifacts; if not, archive to an `archive/` branch or delete.

## Logs & Scratch Files
- `.ticketr/logs/*.log` — local run logs. Route logs to the global PathResolver cache directory so they never touch the repo.
- Any `tmp.*`, `*.log`, `*.bak` (none tracked now, but enforce with git clean or pre-commit).

## Recommended Cleanup Checklist
1. **Immediate purge:** delete sensitive and binary artifacts (`tmp.delete.me.scratchpad`, ticket dumps, build outputs) and rotate the exposed API token.
2. **Standardise “clean” tooling:** add `make clean` / `scripts/clean.sh` that removes binaries, state files, coverage reports, and logs.
3. **Pre-commit guardrails:** integrate `pre-commit` hooks (or Git hooks via `scripts/git-hooks/install.sh`) to block binaries, ticket dumps, and `.scratch` files.
4. **Doc refresh:** update `CONTRIBUTING.md` with “keep the repo tidy” guidance plus instructions for local state paths.
5. **CI enforcement:** extend lint/test jobs with a `git status --porcelain` check after builds to ensure nothing new is generated.

## Ownership & Next Steps
- **Director:** break the checklist into Builder-ready tasks (binary cleanup, credential sanitisation, tooling updates).
- **Builder:** implement cleanup scripts, relocate demos/examples, adjust `.gitignore` if gaps remain.
- **Verifier:** add regression tests (e.g., shell script) ensuring `git status` is clean after CI steps.
- **Scribe:** document cleanup expectations in README/CONTRIBUTING.
- **Steward:** confirm no security regressions and that workspace/TUI flows still reference PathResolver paths.

This pass alone will make the repo look intentional and production-minded—perfect ammunition when we pitch Ticketr’s craftsmanship to technical teams.
