# Changelog - English

All notable changes to this project will be documented in this file.

---

### [v1.0.0] Worker Concurrency Handling (Minimum Version) (2026-01-03)

- **Refactor:** Enhance worker concurrency handling.

### [v0.1.7] CI/Release Configuration Fix (2026-01-02)

- **Fix:** Fixed indentation in `.goreleaser.yaml`.

### [v0.1.6] / [v0.1.5] Build Pipeline Optimization (2026-01-02)

- **Fix:** Simplified GoReleaser step by removing unnecessary skip arguments.

### [v0.1.4] Multi-platform Distribution Support (2026-01-02)

- **Feat:** Added initial `.goreleaser.yaml` configuration.

### [v0.1.3] Package Namespace Standardization (2026-01-02)

- **Fix:** Updated package name from `main` to `oncamq` in `worker.go`.

### [v0.1.2] Scope Visibility Adjustment (2026-01-02)

- **Fix:** Changed package name from `oncamq` to `main`.

### [v0.1.1] Import Path Conflict and Workflow Resolution (2026-01-02)

- **CI:** Fixed package/program conflict at import path.
- **Fix:** Updated release version message in GitHub Actions workflow.

### [v0.1.0] Core Implementation and Idiomatic Patterns (Gopher Way) (2026-01-02)

- **Feat:** Implementation of idiomatic Go worker patterns (removed global state, explicit context).
- **Feat:** Success actions: completed queue management, return values, and attempts.
- **Feat:** GitHub Actions integration for automatic publishing.
- **Docs:** Expanded documentation with contribution guide and realistic examples.
- **Chore:** Renamed module from `go-bullmq-consumer` to `oncamq`.
