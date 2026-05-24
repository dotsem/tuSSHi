# TUSSHI Project Roadmap

### Live Diagnostics & Organization
#### 1. Background Health Checks & Latency Checker (Instant Ping)
* **Goal:** Avoid standard 30-second SSH timeouts by knowing if a server is online before connecting.
* **Mechanism:** 
  * Periodically launch a pool of non-blocking background goroutines performing quick TCP handshakes (e.g. `net.DialTimeout` on the target `HostName:Port` with a 1-2s limit).
  * Send dynamic Bubble Tea messages (`PingResultMsg`) to update the TUI model.
  * Render an interactive status column displaying **Online (Green / Latency in ms)** or **Offline (Red)**.

#### 2. Native Tagging via Lossless Comments
* **Goal:** Enable dynamic filtering and custom groupings in larger configurations without modifying the standard OpenSSH config schema.
* **Mechanism:**
  * Parse special hashtag metadata inside standard configuration file comments (e.g., `# tags: production, database, aws`).
  * Index these tags on load to allow query filtering in the search input (e.g., `tag:production` or `#aws`) or split views based on tag categories.

### Password & Secret Manager Integrations (e.g., Bitwarden, 1Password)
Managing passwords and SSH key passphrases manually is tedious and insecure. TUSSHI will provide elegant integration with command-line password managers.

#### 1. Bitwarden (`bw`) & 1Password (`op`) CLI Integration
* **Mechanism:**
  * Allow users to declare custom properties inside their host blocks that TUSSHI's lossless parser preserves (e.g., `_PassSource bitwarden:item-uuid` or `_PassSource op:vault/item`).
  * Before launching the connection, TUSSHI queries the local password manager CLI to retrieve the password or key passphrase.
  * *Security First:* The TUI will check if the user is currently authenticated with their CLI manager; if not, it will prompt securely inside the TUI for master password/biometric unlock.

#### 2. Dynamic `SSH_ASKPASS` Integration
* **Mechanism:**
  * Directly injecting passwords into interactive shell commands is highly insecure and prone to leaking secrets in shell history or process trees.
  * **Solution:** TUSSHI can act as a local, temporary `SSH_ASKPASS` provider. When launching the `ssh` process, TUSSHI points the `SSH_ASKPASS` environment variable to a secure internal hook or custom short-lived local script. When the native `ssh` client prompts for a passphrase or password, it automatically retrieves the decrypted secret from the password manager and pipes it in safely.

### Advanced Connectivity & Tunnelling
#### 1. Interactive Port Forwarding & SSH Tunnel Visualizer
* **Goal:** Simplify complex database and SOCKS5 proxy port forwarding setups (`LocalForward`, `RemoteForward`, `DynamicForward`).
* **Mechanism:**
  * A dedicated "Tunnels" pane (toggleable via `t`) for the highlighted host.
  * Visually list, add, or toggle active port-forwarding setups.
  * Keep connections alive in the background with a visual status light (Green = Active Tunnel, Red = Stopped).

#### 2. Startup Snippets & Command Bookmarks
* **Goal:** Speed up repetitive remote tasks.
* **Mechanism:**
  * Associate custom command bookmarks with individual hosts (e.g. `tail -f /var/log/nginx/access.log`).
  * Connect using a hotkey (e.g., `Shift+Enter`) to open a command bookmark dropdown.
  * Selecting a bookmark immediately launches the native `ssh` shell and runs the command.

#### 3. Multi-Host Command Broadcast (CSSH Mode)
* **Goal:** Run diagnostic or setup commands across multiple servers concurrently.
* **Mechanism:**
  * Multi-select hosts using `spacebar`.
  * Open an execution input prompt to type a single bash command.
  * Run the command in parallel on all selected servers and display a beautifully formatted, unified report of `stdout` and `stderr`.
