# portman

> Stop fighting EADDRINUSE.

portman is a local developer utility to monitor and manage ports. It shows which process is using each local port and lets you safely stop the process without remembering commands like `lsof`, `netstat`, or `ss`.

## Why?

Every developer running local servers eventually encounters this error:

```txt
Error: listen EADDRINUSE: address already in use :::3000
```

Instead of manually searching for the process holding the port and extracting the PID to kill it, `portman` provides a clean CLI and interactive terminal dashboard to list, filter, and safely stop processes that block development ports.

---

## Features

* **Interactive Live Dashboard**: Visual TUI to monitor active ports, connection states, and processes in real-time.
* **Process Termination**: Safe and forced process killing directly by port number or process name.
* **Smart Safety**: Warnings when attempting to terminate database engines or critical system processes.
* **Script Friendly**: JSON output support for listing ports programmatically.
* **Zero Configuration**: Single statically compiled binary with no external runtimes required.

---

## Installation

### Using Go

If you have Go installed, install `portman` directly:

```bash
go install github.com/nata/portman@latest
```

### Homebrew (macOS / Linux)

```bash
# Placeholder for future homebrew formula
brew install nata/tap/portman
```

### Manual Binary Download

Download pre-compiled binaries directly from the Releases page for your operating system:
* Linux (amd64 / arm64)
* macOS (amd64 / arm64)
* Windows (amd64)

---

## Usage

### Interactive Terminal Dashboard (TUI)

To launch the live interactive dashboard:

```bash
portman
```

This displays a real-time table of all active ports, process statistics, and keybindings.

To launch the dashboard pre-filtered for a specific process or port:

```bash
portman filter node
```

### Command Line Interface (CLI)

#### List Active Ports

List active local ports in a human-readable table (shows LISTEN ports only by default):

```bash
portman list
```

Expected output:
```txt
PORT  PROCESS   PID    STATE   ADDRESS  PROTOCOL
3000  node      12345  LISTEN  0.0.0.0  tcp
5432  postgres  1203   LISTEN  0.0.0.0  tcp
```

To also include established outbound connections:

```bash
portman list --all
```

#### List as JSON

Get machine-readable output:

```bash
portman list --json
```

Expected output:
```json
[
  {
    "port": 3000,
    "address": "0.0.0.0",
    "protocol": "tcp",
    "pid": 12345,
    "process": "node",
    "state": "LISTEN"
  }
]
```

#### Kill a Process by Port

To stop a process using a specific port (e.g., port 3000):

```bash
portman kill 3000
```

Expected interaction:
```txt
Found node PID 12345 using port 3000.
Kill this process? [y/N] y
Killed node (PID 12345).
```

#### Kill a Process by Name

To stop processes by process name (e.g., node):

```bash
portman kill node
```

Expected interaction:
```txt
Found the following processes:
- node (PID: 12345, Port: 3000)
- node (PID: 12346, Port: 3001)
Kill these processes? [y/N] y
Killed node (PID 12345).
Killed node (PID 12346).
```

#### Skip Confirmation Prompt

Use `--yes` to skip the confirmation prompt:

```bash
portman kill 3000 --yes
```

#### Force Kill

To force terminate a process using SIGKILL instead of SIGTERM:

```bash
portman kill 3000 --force
```

---

## CLI Command Reference

* `portman` - Opens the interactive TUI.
* `portman list` - Prints active LISTEN ports in an aligned table.
* `portman list --all` - Prints active LISTEN and ESTABLISHED ports.
* `portman list --json` - Output active ports as JSON.
* `portman kill <port|process>` - Find processes and prompt to stop them.
* `portman filter <query>` - Opens the TUI with a pre-applied search filter.
* `portman version` - Prints current version.

---

## TUI Keybindings

When running inside the interactive terminal dashboard:

* `↑ / ↓`: Navigate active process list.
* `k`: Prompt to kill the selected process.
* `r`: Manually refresh the port list.
* `/`: Enter filter mode to search by port or process name.
* `esc`: Clear active filter or quit.
* `q`: Quit the application.
* `enter`: Confirm actions (e.g., confirming process termination).

---

## Safety Note

* **Graceful Termination**: `portman` attempts graceful termination using SIGTERM first. It only issues SIGKILL when the `--force` flag is specified.
* **Confirmation Prompts**: By default, no process is terminated without explicit confirmation from the user.
* **System Process Protection**: If a process matches a known system or persistent service (e.g., `postgres`, `mysql`, `redis-server`, `nginx`, `docker`, `systemd`, `launchd`), `portman` displays a clear warning prior to prompting for action:
  ```txt
  Warning: postgres (PID 1203) may be a long-running system process.
  ```

---

## Roadmap

* **v0.1.0**: CLI listing (`list`), JSON outputs, and safe termination (`kill`).
* **v0.2.0**: Live terminal interactive dashboard (TUI) powered by Bubble Tea.
* **v0.3.0**: Command filtering and platform-specific improvement stubs.
* **v0.4.0**: Enhanced Windows and cross-platform process resolving support.
* **v1.0.0**: Stable, optimized cross-platform release.

---

## Contributing

Contributions are welcome! Please follow these guidelines:
1. Fork the repository and create your branch.
2. Ensure all changes pass unit tests (`go test ./...`).
3. Format all code using `gofmt` before committing.
4. Keep dependencies minimal and standard library-focused where possible.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
