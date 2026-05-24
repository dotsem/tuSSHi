# TUSSHI Project Roadmap & Checklist

## Live Diagnostics & Organization
### Background Health Checks & Latency Checker (Instant Ping)

**Goal:** Avoid standard 30-second SSH timeouts by knowing if a server is online before connecting.

- [ ] Pool of non-blocking background goroutines performing quick TCP handshakes (e.g., `net.DialTimeout` with a 1-2s limit on the target `HostName:Port`).
- [ ] Dynamic Bubble Tea messages (`PingResultMsg`) to update the TUI model reactively.
- [ ] Interactive status column displaying **Online (Green with latency in ms)** or **Offline (Red)**.

### Native Tagging via Lossless Comments

**Goal:** Enable dynamic filtering and custom groupings in larger configurations without modifying the standard OpenSSH config schema.

- [ ] Lossless parser support for custom hashtag metadata inside standard configuration comments (e.g., `# tags: production, database, aws`).
- [ ] Indexing tags on load to enable query filtering in the search input (e.g., `tag:production` or `#aws`).
- [ ] Categorized TUI split views based on tag groups.

---

## Password & Secret Manager Integrations

### Bitwarden (`bw`) & 1Password (`op`) CLI Integration

**Goal:** Securely fetch passwords and SSH key passphrases without manual entries or insecure hardcoding.

- [ ] Allow declaration of custom properties inside host blocks that TUSSHI's lossless parser preserves (e.g., `_PassSource bitwarden:item-uuid` or `_PassSource op:vault/item`).
- [ ] Direct queries to local password manager CLIs to retrieve passwords or key passphrases prior to launching a connection.
- [ ] Secure TUI prompt flow for master password entry or biometric unlock if the user is unauthenticated with their manager CLI.

### Dynamic `SSH_ASKPASS` Integration

**Goal:** Provide secure, hands-off secret injection into the standard SSH client without putting sensitive data in process arguments or environment variables.

- [ ] Act as a local, short-lived `SSH_ASKPASS` provider.
- [ ] Point the `SSH_ASKPASS` environment variable to a secure internal hook or custom short-lived local script.
- [ ] Dynamically pipe decrypted credentials from secret manager integrations directly to the native `ssh` prompt.

---

## Advanced Connectivity & Tunnelling

### Interactive Port Forwarding & SSH Tunnel Visualizer

**Goal:** Simplify complex database and SOCKS5 proxy port forwarding setups (`LocalForward`, `RemoteForward`, `DynamicForward`).

- [ ] Dedicated "Tunnels" pane (toggleable via `t`) for the highlighted host.
- [ ] Visual interface to list, add, or toggle active port-forwarding tunnels.
- [ ] Background connection keep-alives with a real-time status light (Green = Active, Red = Stopped).

### Startup Snippets & Command Bookmarks

**Goal:** Speed up repetitive remote administrative tasks.

- [ ] Association of custom command bookmarks with individual hosts (e.g., `tail -f /var/log/nginx/access.log`).
- [ ] Command bookmark dropdown menu accessible via connection hotkeys (e.g., `Shift+Enter`).
- [ ] Auto-executing shell launcher to immediately connect and run the selected bookmark.

### Multi-Host Command Broadcast (CSSH Mode)

**Goal:** Run diagnostic or setup commands across multiple servers concurrently.

- [ ] Multi-select hosts in the TUI using `spacebar`.
- [ ] Unified command execution overlay to type a single broadcast command.
- [ ] Parallel command execution on all selected servers with a beautifully formatted, collapsable report of `stdout` and `stderr`.
