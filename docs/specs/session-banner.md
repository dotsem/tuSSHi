# Spec: Session Identification Banner

| | |
|---|---|
| **Status** | Research / parked |
| **Last researched** | 2026-11-09 |
| **Related code** | [`internal/ssh/launcher.go`](../../internal/ssh/launcher.go), [`internal/tui/keybinds.go`](../../internal/tui/keybinds.go) (`connectSelectedHost`) |

## 1. Problem

With several terminals open, it is not obvious which server a given session is connected to. On headless systems (TTY console, no window manager) there is no window title to fall back on either.

## 2. Goal

Show persistent session context while connected via tuSSHi:

- host alias / hostname
- connection status (connected, reconnecting, dropped)
- local time
- optionally: tags, latency, session duration

### Non-goals

- Modifying the remote host (no `PS1` / shell rc injection).
- Becoming a general-purpose terminal multiplexer.

## 3. Why a "real" banner is hard

Rendering a line is trivial. Reserving it is not:

- Remote programs (`vim`, `htop`, `less`, `tmux`, the shell) assume they own the full screen. They use absolute cursor positioning (`ESC[1;1H`), clear the screen (`ESC[2J`), reset scroll regions (`ESC[r`) and switch to the alternate screen (`ESC[?1049h`). Each of these overwrites or wipes a banner.
- `ssh` forwards the real terminal size to the remote, so the remote never knows a row is reserved.

A reliable banner therefore needs a component between `ssh` and the real terminal that owns the screen: a pty plus a VT emulator. That is the core of what `tmux` does.

Today tuSSHi hands the TTY to `ssh` entirely via `tea.ExecProcess` and has no control during the session.

## 4. Options

| # | Approach | Effort | Robustness | Notes |
|---|---|---|---|---|
| A | Window title (OSC 2) + per-host background tint (OSC 11), set before exec and reset after | ~30 LOC | High | Solves "which terminal is which" via tab bar and colour. No live info. No effect on bare TTY consoles. |
| B | Run `ssh` inside `tmux` with a generated config whose status bar shows host/time/state | ~100 LOC | High | Works on headless consoles. Requires `tmux` on `PATH`. Nested-tmux and prefix-key friction (prefix can be `None`). |
| C | Pty proxy that parses output and rewrites cursor/scroll sequences (offset by banner height) | Weeks | Low–medium | Hack. Alt screen, clears and scroll-region resets need special handling. Endless edge cases. |
| D | Embedded terminal emulator (pty + VT emulator + own renderer); banner is a normal Bubble Tea view | Weeks–months | Medium → high | The "proper" approach. Enables multi-session tabs/splits. Loses native scrollback, selection and image protocols. |

### Relevant Go packages

| Purpose | Package |
|---|---|
| pty | `github.com/creack/pty`, `github.com/charmbracelet/x/xpty` |
| VT emulator | `github.com/charmbracelet/x/vt`, `github.com/hinshun/vt10x`, `github.com/liamg/darktile` |
| ANSI parser only (option C) | `github.com/charmbracelet/x/ansi` (already indirect dep), `github.com/danielgatis/go-vte` |
| Bubble Tea terminal component (option D) | `github.com/taigrr/bubbleterm` |

## 5. Evaluation: `taigrr/bubbleterm`

Snapshot as of 2026-11-09.

### Facts

- Created 2025-07, 31 stars, 0BSD licence, actively maintained (last push 2026-10-31), has external contributors.
- Thin layer: a Bubble Tea `Model` (`bubbleterm.go`), key/mouse translation (`keys.go`) and an `emulator` package. Actual VT emulation is `github.com/charmbracelet/x/vt` (untagged pseudo-version). Maturity is effectively that of Charm's `x/vt`.
- Recent history is the expected bug tail, being actively fixed: idle CPU burn, poll goroutine leak, key translation rewrite, DSR query responses not reaching the child.
- Includes a `cmd/multiwindow` example (multiple terminals in one TUI).

### Blocker: Bubble Tea v2

| | tuSSHi | bubbleterm |
|---|---|---|
| bubbletea | `github.com/charmbracelet/bubbletea v1.3.10` | `charm.land/bubbletea/v2 v2.0.10` |
| lipgloss | `v1.1.0` | `charm.land/lipgloss/v2` |
| `View()` return type | `string` | `tea.View` |
| `go` directive | 1.26.3 | 1.27.1 |

A v2 model cannot be embedded in a v1 program. Adopting bubbleterm requires migrating tuSSHi to the v2 stack first:

- `tea.KeyMsg` → `tea.KeyPressMsg`
- `View()` returns `tea.View`; alt screen / cursor become view properties
- lipgloss v2 colour handling changes
- bubbles v2; confirm `huh` has a v2-compatible release
- bump `go` directive to 1.27 and the CI setup-go version

The v2 migration is likely worthwhile regardless, but it is a separate piece of work with its own risk.

### Architectural consequences of option D

```
keyboard ──► bubbletea key parse ──► bubbleterm re-encode ──► pty ──► ssh
screen   ◄── bubbletea render ◄── banner + x/vt grid ◄────── pty ◄── ssh
```

- **Lossy input round-trip**: keystrokes are decoded into `KeyPressMsg` and re-encoded. Alt combos, F-keys, kitty keyboard protocol and bracketed paste are the likely failure points. Biggest fidelity risk.
- **Lifetime coupling**: a tuSSHi panic or hang kills the SSH session. Today a crash cannot affect a running session.
- **Escape key**: a key is needed to leave the embedded terminal; it will conflict with something remote.
- **Lost native features**: terminal scrollback (now inside the emulator), OSC 52 clipboard, kitty/sixel images; native selection includes the banner.
- **Throughput**: bulk output (`cat` of a large file, build logs) goes through emulate → render → diff.
- **Upside**: multi-session tabs/splits and live connection status become possible.

## 6. Spike plan (no tuSSHi changes required)

Clone bubbleterm, change `cmd/bubbleterm/main.go` to run `exec.Command("ssh", "<host>")`, wrap it with a one-line banner, and verify:

- [ ] `vim` / `nvim`: insert, visual, `:term`, mouse
- [ ] `htop`, `less`, `man`: alt screen enter and restore
- [ ] `tmux` on the remote (nested)
- [ ] window resize mid-session
- [ ] large multi-line paste (bracketed paste)
- [ ] Alt+key, Ctrl+arrows, F-keys, Home/End in the shell
- [ ] `cat` of a ~50 MB file: CPU usage and responsiveness
- [ ] Unicode: emoji, CJK, Nerd Font glyphs
- [ ] 10+ minutes idle: CPU and connection stability
- [ ] `exit` and dropped connection: clean teardown

**Decision rule**: if input and paste fidelity pass, proceed with the v2 migration and add bubbleterm as an opt-in "banner mode" with `ExecProcess` kept as default. If they fail, use option B instead.

## 7. Recommended phasing

1. **Option A** — ship independently, no dependencies. Highest value for lowest effort.
2. **Spike** (section 6) for option D.
3. Depending on the spike: **option D** behind a flag (after the v2 migration), or **option B** when `tmux` is available.
4. Option C: do not pursue.

## 8. Open questions

- Per-host colour source: a lossless config comment consistent with existing tags (e.g. `# color: red`), or derived automatically from tags (e.g. `production` → red)?
- Should banner mode be global, per host, or per tag?
- Option B: generate a temporary tmux config per session, or ship a static one?
- Option D: which key leaves/detaches the embedded session?

## 9. Prerequisite fix: shell injection in the launcher

`NewSSHCommand` in `internal/ssh/launcher.go` interpolates the host alias into a `bash -c` string via `fmt.Sprintf`. An alias containing `;`, `$()` or quotes is executed as shell code. Any of the options above should start by replacing this with:

- printing the "connecting…" line and OSC sequences from Go, then `exec.Command("ssh", "--", alias)`, or
- passing the alias as a positional argument: `bash -c '… "$1"' _ alias`.
