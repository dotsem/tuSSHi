

# tuSSHi AI Coding Guidelines

This document defines the coding standards, architectural decisions, and practices for the tuSSHi codebase. Every AI assistant and developer working on this project **must** adhere strictly to these rules.

---

## 1. Universal Principles

### 1.1 Comments & Documentation
- **No Inline Comments**: Do not write comments within function bodies to explain syntax, control flow, or standard actions. If logic is complex, rewrite it to be self-documenting or extract it into a descriptive helper function.
- **Redundant Comments Ban**: Comments that echo the code syntax (e.g., `counter++ // increment counter`) are strictly forbidden.
- **Allowed Comments**:
  - Package/Library-level docstrings and standard Go package documentation.
  - Standard API documentation on public interfaces/functions/types.
  - High-level `// Why:` comments explaining non-obvious architecture or constraints.
  - Actionable `// TODO:` comments explaining shortcuts taken or rough edges, including what should be done in the future.
- **Style**: Keep short comments lowercase. Use capital letters only for multi-line contextual explanations.
- **Self-Debates**: Never start self-debates in comments or code.

### 1.2 Function & File Design
- **Single Responsibility (SOLID)**: Every function must do exactly one thing and do it right. If a function performs multiple operations, break it down.
- **DRY (Don't Repeat Yourself)**: Generalized utilities must be extracted. Never duplicate logic or functions.
- **Function Names**: The name must clearly reflect the function's single action. Avoid generic names or names with "And" (which signals multiple responsibilities).
- **File Length Limit**: A strict ceiling of **300 lines per file**. If a file exceeds this:
  - Decouple, group, and split into sub-modules or logical extensions unless structurally impossible.
  - **Important for AI**: Ask the developer for permission before performing a major file split/refactoring. Do not do it autonomously.
- **Clean Code**: No dead, unused, or commented-out code. Use named constants instead of magic numbers.

### 1.3 Development & Workflows
- **Test Guard**: **NEVER** execute test suites natively (`go test` or `just test`, etc.) unless explicitly instructed by the developer.
- **Linter & Formatter Integrity**: Never bypass formatters or linters. Before presenting or pushing code, formatting and linting checks must pass cleanly.
- **AI Tooling Constraints**:
  - AI assistants must **never** run git commits or push changes under any circumstance. Commits are reserved exclusively for human execution.
  - AI assistants must **never** begin execution of an implementation plan without explicit authorization/approval from the developer. A plan is a plan until the human developer explicitly says "go ahead".
- **Commit Format**: All commit messages must follow the `action(part): description` format.

---

## 2. Go (TUI & SSH Connection Manager)

### 2.1 Tooling & Verification
- **Formatting**: Format using standard Go tools via:
  ```bash
  just fmt
  ```
- **Linting**: Lint using the standard linter suite via:
  ```bash
  just lint
  ```
  All code must strictly pass `gofmt` and `golangci-lint` default rules (including `errcheck`, `govet`, `staticcheck`, and `revive`).

### 2.2 Structural Rules & Architecture
- **Folder Structure**: Follow the Standard Go Project Layout:
  - `/cmd` -> Application entry points (main.go).
  - `/internal` -> Private application/business logic (e.g., config handling, TUI layout/views/components, SSH file parsing); absolute ban on external package exposure.
  - `/pkg` -> Explicitly exportable, reusable utility packages.
- **TUI Architecture**:
  - The application is built using the Bubble Tea framework (Model-View-Update pattern).
  - UI components (under `/internal/tui/components`) should be clean, modular, and self-contained.
  - Theme configurations, colors, and layout borders are defined centrally using `lipgloss` under `/internal/tui/style` and `/internal/tui/theme`. Use these existing styles and tokens instead of hardcoding styles inline.
  - Forms and text inputs should utilize standard components or `github.com/charmbracelet/huh` where appropriate.
- **SSH Configuration Management**:
  - Manage the primary source of truth (`~/.ssh/config` or alternative config paths) losslessly using the AST-based parser (`github.com/kevinburke/ssh_config`).
  - Do not overwrite configuration files destructively; preserve comments, custom spacing, and formatting.
- **Implicit Interfaces**: Leverage Go's implicit interface implementation. Do not define interfaces before they are actually needed by consumers.
- **Error & Logging**:
  - Return early with errors. Avoid deeply nested conditional blocks (e.g., write `if err != nil { return err }` instead of deep nesting).
  - Avoid writing manual logging workarounds. If structured logs or diagnostics are needed, use standard CLI reporting patterns or log outputs.

### 2.3 Testing
- **Unit Testing**:
  - Use `t.Run` to group subtests.
  - Name test files explicitly matching target: `xxx_test.go`.
  - Use `stretchr/testify/assert` consistently.
  - Keep test components mocked where necessary to avoid real filesystem mutations unless specifically testing lossless filesystem writes.

---

## 3. Verification Commands Reference

Run the following commands within their respective targets to ensure compliance (note: tests should only be run when explicitly instructed):

| Action | Command | Description |
|---|---|---|
| **Format** | `just fmt` | Formats all Go files using `gofmt` |
| **Lint / Analyze** | `just lint` | Runs `golangci-lint` on the project |
| **Run Application**| `just run` | Runs the TUI application directly |
| **Test Run** | `go test ./...` | Runs the test suite (**Only if explicitly instructed**) |
