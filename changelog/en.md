# Changelog - English

All notable changes to this project will be documented in this file.

---

### [v0.1.7] - 2026-01-02
**Friendly Name:** Precision Tweaks
**Synopsis:** Minor technical fix in the distribution configuration (GoReleaser) to ensure builds are generated correctly.

- **Fix:** Fixed indentation in `.goreleaser.yaml`.

### [v0.1.6] / [v0.1.5] - 2026-01-02
**Friendly Name:** Delivery Refinement
**Synopsis:** Simplification of the release process, removing redundant steps to make the publishing cycle more agile.

- **Fix:** Simplified GoReleaser step by removing unnecessary skip arguments.

### [v0.1.4] - 2026-01-02
**Friendly Name:** Binary Automation
**Synopsis:** Introduction of GoReleaser to automate the creation of binaries and multi-platform releases.

- **Feat:** Added initial `.goreleaser.yaml` configuration.

### [v0.1.3] - 2026-01-02
**Friendly Name:** Identity Consolidation
**Synopsis:** Final standardization of the package name to `oncamq`, ensuring consistency between the code and the repository.

- **Fix:** Updated package name from `main` to `oncamq` in `worker.go`.

### [v0.1.2] - 2026-01-02
**Friendly Name:** Scope Adjust
**Synopsis:** Temporary fix in package scope for internal import testing.

- **Fix:** Changed package name from `oncamq` to `main`.

### [v0.1.1] - 2026-01-02
**Friendly Name:** Dependency Fix
**Synopsis:** Resolution of package and program conflicts in the Go import path.

- **CI:** Fixed package/program conflict at import path.
- **Fix:** Updated release version message in GitHub Actions workflow.

### [v0.1.0] - 2026-01-02
**Friendly Name:** The Birth of OncaMQ
**Synopsis:** Consolidated release of the consumer core.

- **Feat:** Implementation of idiomatic Go worker patterns (removed global state, explicit context).
- **Feat:** Success actions: completed queue management, return values, and attempts.
- **Feat:** GitHub Actions integration for automatic publishing.
- **Docs:** Expanded documentation with contribution guide and realistic examples.
- **Chore:** Renamed module from `go-bullmq-consumer` to `oncamq`.
