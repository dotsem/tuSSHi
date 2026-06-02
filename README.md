<img width="200" alt="tusshi logo" src="assets/tusshi.png" align="right">

# tuSSHi

tuSSHi is a terminal user interface connection manager for your standard SSH configuration files, built in Go using the Bubble Tea framework.

Most SSH managers force you into using custom JSON or YAML databases, locking you into their tool and making it impossible to share host configurations with other terminal programs, Ansible scripts, or IDEs. tuSSHi takes a different approach: your standard ~/.ssh/config file remains the single source of truth.

---

## Why tuSSHi?

Unlike other SSH managers that overwrite your configuration files destructively, tuSSHi is built on top of a lossless Abstract Syntax Tree parser. When you add, edit, or delete hosts through the TUI, tuSSHi surgically modifies only the relevant blocks in your config files. Your custom indentation, comment lines, custom spacing, and unparsed custom variables are preserved exactly as they were.

Furthermore, if you keep your server configurations modular using standard OpenSSH Include directives, tuSSHi automatically detects them and organizes your connections into visual tabs mapped directly to each physical file on disk.

---

## Installation

### Prerequisites

To compile tusshi from source, you need Go 1.26.3 or higher installed on your system.

### Building from Source

Clone the repository and build the binary:

```bash
git clone https://github.com/dotsem/tusshi.git
cd tusshi
go build -o tusshi cmd/tusshi/main.go
```

You can run the compiled binary directly or move it into a directory in your system PATH (such as `/usr/local/bin`):

```bash
mv tusshi /usr/local/bin/
```

### Installing on Linux 

```bash
./linux-install.sh
```

For development it is recommended to create a symlink:

```bash
ln -s $(pwd)/cmd/tusshi/tusshi ~/.local/bin/tusshi
```

### Installing on MacOS

Currently built from source, homebrew support will be added later on.

### Installing on Windows

Currently built from source, package manager support will be added later on.

---

## Getting Started

To launch tusshi with your default primary SSH configuration file (~/.ssh/config):

```bash
tusshi
```

If you want to manage or browse a different SSH configuration file, you can pass its path as an argument:

```bash
tusshi -config ~/path/to/alternative/ssh_config
```

---

## Interface and Navigation

The interface is structured into four main areas:
1. Header: Displays the application name and tabs corresponding to your config files.
2. Main Body: A table listing your connections (Alias, Hostname, User, Port, and Source File).
3. Search / Status Bar: Shows active search queries or error/alert banners.
4. Footer: Displays a quick summary of keybindings and application state.

### Keyboard Shortcuts

Navigation and core commands are designed to be fast and keyboard-friendly:

* Up / Down (or j / k): Navigate through the connections list.
* Left / Right (or h / l): Switch active configuration tabs.
* Enter: Connect to the selected host (launches your system's native ssh client).
* / (Slash): Enter search mode to fuzzy-filter connections by alias, hostname, or user.
* : (Colon): Enter command mode to perform management tasks.
* Esc: Exit search, command mode, or active overlay overlays (like help or forms) and return to normal mode.
* ? (Question mark): Open help overlay showing available commands.

---

## Command Mode Reference

Type a colon (:) in normal mode to open the command prompt. The following commands are supported:

* :new (or :n): Open the interactive connection creator form.
* :edit (or :e): Edit the selected connection.
* :delete-config <filename> (or :rmconf): Delete an empty configuration file.
* :add-config <filename> (or :addconf): Add a new configuration file.
* :rename-config <old-name> <new-name> (or :mvconf): Rename an existing configuration file.
* :help (or :h or :?): Open the interactive help overlay showing available commands.
* :quit (or :q): Exit the application.

---

## Technical Limitations

* Native SSH Execution: tusshi does not implement its own SSH client protocol. Instead, it constructs a subprocess executing your system's native ssh client in-place. While this ensures perfect compatibility with your keys, ssh-agent, and local configuration, it means tusshi requires standard ssh to be accessible in your system environment.
* Deep Inclusions: To prevent circular dependencies or infinite loops, tusshi limits nested configuration file inclusion (via Include directives) to a maximum depth of 5 files.
* Host Patterns: While tusshi handles standard Host blocks perfectly, it excludes wildcard blocks (e.g. Host *) from the main interactive listing to keep the focus solely on functional target connections. It still reads these wildcard blocks to resolve inherited attributes.

---

## Contributing

Contributions are welcome. If you would like to help improve tusshi:

1. Formulate clear, focused pull requests targeting specific features or bugs.
2. Adhere strictly to the Go programming guidelines and coding standards. Keep files focused and limit file length (under 300 lines is preferred).
3. Ensure no redundant comments are added (only comment on why a piece of code works or handles complex scenarios).

View the contributing guidelines in [CONTRIBUTING](CONTRIBUTING).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.