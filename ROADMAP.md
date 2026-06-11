# TUSSHI Project Roadmap

This document outlines the planned design improvements, features, and core focus areas for TUSSHI.

## 1. Live Diagnostics & Organization

### Background Health Checks & Latency Checker (Instant Ping)
- [x] Pool of non-blocking background goroutines performing quick TCP handshakes (e.g., `net.DialTimeout` with a 1-2s limit on the target `HostName:Port`).
- [x] Dynamic Bubble Tea messages (`PingResultMsg`) to update the TUI model reactively.
- [x] Interactive status column displaying **Online (latency in ms)** or **Offline**.

### Native Tagging via Lossless Comments
- [ ] Lossless parser support for custom hashtag metadata inside standard configuration comments (e.g., `# tags: production, database, aws`).
- [ ] Indexing tags on load to enable query filtering in the search input (e.g., `tag:production` or `#aws`).
- [ ] Categorized TUI views or tab structures based on tag groups.

---

## 2. Terminal & Theme Integration

### System Theme & Adaptive Color Palette
- [ ] Integrate dynamic terminal color profile query (using libraries like `termenv` or ANSI queries) to auto-detect dark vs. light terminal backgrounds.
- [ ] Adapt TUI element coloring dynamically to respect user terminal themes (ANSI palette colors) rather than using hardcoded hex color values.
- [ ] Implement system-level accessibility support (respect `NO_COLOR` environment variable).

---

## 3. Remote Execution Shortcuts

### Command Bookmarks & Startup Snippets
- [ ] Support bookmarking remote commands inside the configuration file using lossless comments (e.g., `# bookmark: Tail access logs | tail -f /var/log/nginx/access.log`).
- [ ] Interactive dropdown or list overlay to select a command bookmark for the active host.
- [ ] Direct invocation of the system ssh binary with the selected command bookmark appended (e.g., `ssh host -t "command"`).
