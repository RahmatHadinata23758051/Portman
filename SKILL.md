# SKILL.md

## Role

You are an autonomous software engineering agent working inside this repository.

Your job is to build **portman**, a developer-first terminal port manager.

portman helps developers solve this common problem:

```txt
Error: listen EADDRINUSE: address already in use
```

Instead of forcing users to remember commands like `lsof`, `netstat`, or `ss`, portman provides a clean CLI and TUI to list active local ports, identify the owning process, and safely kill processes that block development ports.

The project must be useful, clean, maintainable, open-source friendly, and suitable for an MIT licensed public GitHub repository.

---

## Project Identity

Project name:

```txt
portman
```

Tagline:

```txt
Stop fighting EADDRINUSE.
```

Description:

```txt
Live port manager for developers. See what's running, inspect the process, and safely kill what's in the way.
```

Primary target users:

* Web developers
* Backend developers
* Full-stack developers
* Students learning development
* DevOps beginners
* Developers who frequently run local servers

This is not an enterprise network monitor.
This is a local developer utility.

---

## Core Principles

Always follow these principles:

1. Keep the project simple and focused.
2. Prioritize developer experience.
3. Avoid unnecessary features.
4. Make every command predictable.
5. Do not kill processes without confirmation unless explicitly requested.
6. Prefer safe termination before force killing.
7. Keep output clean, readable, and script-friendly.
8. Maintain cross-platform design, but implement incrementally.
9. Write code that is easy to test.
10. Do not turn this into a web dashboard.

---

## Tech Stack

Use Go.

Required libraries:

```txt
github.com/spf13/cobra
github.com/charmbracelet/bubbletea
github.com/charmbracelet/lipgloss
github.com/shirou/gopsutil/v4
gopkg.in/yaml.v3
```

Reasons:

* Go produces single binaries.
* Go is fast and cross-platform.
* Cobra is standard for CLI apps.
* Bubble Tea is mature for TUI apps.
* Lip Gloss gives clean terminal styling.
* gopsutil helps with process and connection data.

Go version: 1.21 or later.

---

## Platform Priority

Implement and test in this order:

1. macOS (primary)
2. Linux
3. Windows

Do not write Windows-specific code until macOS and Linux are working correctly.

---

## License

Use MIT License.

The repository must include:

```txt
LICENSE
```

with the standard MIT License text.

---

## Initial Project Structure

Create this structure:

```txt
portman/
├── cmd/
│   ├── root.go
│   ├── list.go
│   ├── kill.go
│   └── version.go
├── internal/
│   ├── ports/
│   │   ├── scanner.go
│   │   └── types.go
│   ├── process/
│   │   └── kill.go
│   ├── output/
│   │   ├── table.go
│   │   └── json.go
│   └── tui/
│       ├── model.go
│       ├── update.go
│       └── view.go
├── .github/
│   └── workflows/
│       └── release.yml
├── README.md
├── LICENSE
├── go.mod
├── go.sum
└── main.go
```

Keep files small and organized.

---

## CLI Commands

Implement these commands.

### 1. Open TUI

```bash
portman
```

Opens the live terminal dashboard.

In early versions, if the TUI is not ready, this command may print a helpful message and suggest:

```bash
portman list
```

But the final goal is to make `portman` open the TUI.

---

### 2. List Active Ports

```bash
portman list
```

Print all active local ports in a clean table.

Expected fields:

```txt
PORT
PROCESS
PID
STATE
ADDRESS
PROTOCOL
```

Example:

```txt
PORT    PROCESS       PID     STATE        ADDRESS        PROTOCOL
3000    node          12345   LISTEN       127.0.0.1      tcp
5432    postgres      1203    LISTEN       127.0.0.1      tcp
8080    go            52104   ESTABLISHED  0.0.0.0        tcp
```

---

### 3. JSON Output

```bash
portman list --json
```

Output machine-readable JSON.

Example:

```json
[
  {
    "port": 3000,
    "process": "node",
    "pid": 12345,
    "state": "LISTEN",
    "address": "127.0.0.1",
    "protocol": "tcp"
  }
]
```

---

### 4. Kill by Port

```bash
portman kill 3000
```

Find the process using port `3000` and ask for confirmation before terminating it.

Expected interaction:

```txt
Found node PID 12345 using port 3000.
Kill this process? [y/N]
```

Default answer must be no.

---

### 5. Kill by Process Name

```bash
portman kill node
```

Find processes matching the name `node` and ask for confirmation.

If multiple processes match, show all matches and ask before killing.

---

### 6. Force Kill

```bash
portman kill 3000 --force
```

Force kill the process.

On Unix-like systems, safe termination should use SIGTERM first.
Force kill should use SIGKILL only when requested.

On Windows, use the appropriate process kill mechanism.

---

### 7. Yes Flag

```bash
portman kill 3000 --yes
```

Skip confirmation.

This is useful for scripting.

Do not make `--yes` imply `--force`.

---

### 8. Filter TUI

```bash
portman filter node
portman filter 3000
```

Open the TUI with a pre-applied filter.

Filter matches against both process name and port number.
Filter is case-insensitive.

---

### 9. Version

```bash
portman version
```

Print version.

Example:

```txt
portman v0.1.0
```

---

## TUI Requirements

The TUI should look clean, minimal, and developer-friendly.

Default refresh interval: 2 seconds.
Refresh interval is not configurable in v1.

Main screen should show:

* App name and version
* Refresh countdown
* Active ports count
* Listening ports count
* Established connections count
* Table of active ports
* Footer with keybindings
* Current OS info if easy to obtain

Example layout:

```txt
portman v0.1.0                                      refreshing in 2s

┌ active ports ┐   ┌ listening ┐   ┌ established ┐
│ 6            │   │ 5         │   │ 1           │
└──────────────┘   └───────────┘   └─────────────┘

PORT     PROCESS       PID      STATE        ACTION
:5432    postgres      1203     LISTEN       kill
:6379    redis-server  1391     LISTEN       kill
:8080    go            52104    ESTAB        kill
:8888    python        33901    LISTEN       kill

↑↓ navigate    k kill    r refresh    / filter    q quit
```

---

## TUI Keybindings

Implement:

```txt
↑ / ↓       Navigate rows
k           Kill selected process
r           Refresh manually
/           Filter
q           Quit
esc         Quit or leave filter mode
enter       Confirm action when confirmation is active
```

When killing from the TUI, always show confirmation first:

```txt
Kill process node PID 12345 using port 3000? [y/N]
```

---

## Safety Rules

Process killing must be safe by default.

Rules:

1. Never kill without confirmation unless `--yes` is provided.
2. Never force kill unless `--force` is provided.
3. Warn when killing common long-running or system-related processes.
4. Default answer must always be no.
5. If the process cannot be killed due to permissions, show a clear message.

Warn for process names such as:

```txt
systemd
launchd
svchost
postgres
mysql
redis-server
docker
dockerd
containerd
sshd
nginx
apache2
```

Example warning:

```txt
Warning: postgres may be a long-running database process.
Use --force or confirm carefully if you really want to stop it.
```

Do not block the user completely. Just warn clearly.

---

## Port Scanner Behavior

The scanner should return normalized data.

Create a type similar to:

```go
type PortEntry struct {
    Port     uint32 `json:"port"`
    Address  string `json:"address"`
    Protocol string `json:"protocol"`
    PID      int32  `json:"pid"`
    Process  string `json:"process"`
    State    string `json:"state"`
}
```

The scanner should:

* list TCP connections
* prioritize listening ports
* include established connections when available
* map port to PID
* map PID to process name
* handle permission errors gracefully
* avoid crashing on missing process info

If process name cannot be resolved, use:

```txt
unknown
```

---

## Output Requirements

### Human Table Output

Keep it clean and aligned.

Do not print noisy debug information.

Good:

```txt
PORT    PROCESS       PID     STATE
3000    node          12345   LISTEN
```

Bad:

```txt
map[pid:12345 port:3000 state:LISTEN]
```

---

### JSON Output

Must be valid JSON.

Do not mix logs with JSON output.

If `--json` is used, stdout must contain JSON only.

Errors should go to stderr.

---

## Error Handling

Be helpful.

Examples:

If no active ports:

```txt
No active ports found.
```

If port is not used:

```txt
No process found using port 3000.
```

If permission denied:

```txt
Permission denied while reading process information.
Try running portman with elevated privileges.
```

If invalid port:

```txt
Invalid port: 99999
Port must be between 1 and 65535.
```

---

## Versioning

Start with:

```txt
v0.1.0
```

Use semantic versioning.

Suggested roadmap:

```txt
v0.1.0  CLI list and kill
v0.2.0  TUI dashboard
v0.3.0  JSON output and filtering
v0.4.0  Windows support improvements
v1.0.0  Stable cross-platform release
```

---

## Development Phases

### Phase 1 — Foundation

Tasks:

* Initialize Go module
* Add Cobra
* Add command structure
* Add README
* Add MIT License
* Add version command

Acceptance criteria:

```bash
go run . version
```

prints:

```txt
portman v0.1.0
```

---

### Phase 2 — Port Listing

Tasks:

* Implement port scanner
* Implement `portman list`
* Print table output
* Handle empty result
* Handle permission errors

Acceptance criteria:

```bash
go run . list
```

prints active ports in a readable table.

---

### Phase 3 — JSON Output

Tasks:

* Add `--json`
* Ensure valid JSON only on stdout
* Add stable struct tags

Acceptance criteria:

```bash
go run . list --json
```

prints valid JSON.

---

### Phase 4 — Kill by Port

Tasks:

* Implement `portman kill <port>`
* Find process by port
* Ask confirmation
* Add `--yes`
* Add `--force`
* Show useful success/error messages

Acceptance criteria:

```bash
go run . kill 3000
```

finds the process and asks before killing.

---

### Phase 5 — Kill by Name

Tasks:

* Support process name argument
* Show all matches
* Ask confirmation
* Avoid accidental mass kill

Acceptance criteria:

```bash
go run . kill node
```

lists matching processes and asks confirmation.

---

### Phase 6 — TUI

Tasks:

* Implement Bubble Tea model
* Show active ports table
* Add navigation
* Add refresh every 2 seconds
* Add filter (case-insensitive, matches name and port)
* Add kill confirmation
* Style with Lip Gloss

Acceptance criteria:

```bash
go run .
```

opens a usable TUI.

---

### Phase 7 — Release Workflow

Tasks:

* Add GitHub Actions build
* Build for Linux, macOS, Windows
* Upload release artifacts
* Document install methods

Acceptance criteria:

A GitHub release can publish binaries for:

```txt
linux-amd64
linux-arm64
darwin-amd64
darwin-arm64
windows-amd64.exe
```

---

## What "Done" Looks Like

The project is considered complete for v1 when:

* `go build ./...` succeeds with no errors
* `go test ./...` passes
* `portman list` works correctly on macOS and Linux
* `portman kill` works with confirmation on macOS and Linux
* `portman` opens a functional TUI
* README is complete and accurate
* LICENSE file is present
* GitHub Actions release workflow is present

Do not declare the project done until all of the above are true.

---

## README Requirements

README must include:

1. Project name
2. Tagline
3. Problem statement
4. Screenshot or mockup
5. Features
6. Installation
7. Usage examples
8. Commands
9. Safety note
10. Roadmap
11. Contributing
12. License

README tone:

* Clear
* Confident
* Not overhyped
* Beginner-friendly

---

## Example README Opening

````md
# portman

> Stop fighting EADDRINUSE.

portman is a live port manager for developers. It shows which process is using each local port and lets you safely stop the process without remembering `lsof`, `netstat`, or `ss`.

## Why?

Every developer eventually sees this:

```txt
Error: listen EADDRINUSE: address already in use :::3000
```

portman gives you one clean command to find and fix it.
````

---

## Testing Requirements

Add tests where practical.

Prioritize testing:

* argument parsing
* port validation
* output formatting
* process matching
* JSON output
* safety confirmation logic

Do not require real process killing in unit tests.

Use interfaces so killing can be mocked.

---

## Code Quality Rules

Follow these rules:

1. Use clear package names.
2. Avoid global mutable state.
3. Keep command logic thin.
4. Put business logic in internal packages.
5. Avoid large files.
6. Avoid unnecessary abstractions.
7. Always run `gofmt`.
8. Always run `go test ./...`.
9. Do not commit generated binaries.
10. Do not print debug logs in normal output.

---

## Commit Style

Use simple conventional commits:

```txt
feat: add list command
feat: add port scanner
feat: add kill by port
fix: handle missing process name
docs: update readme usage
chore: add release workflow
```

---

## Out of Scope for v1

Do not implement these in v1:

* Web UI
* Remote machine monitoring
* SSH monitoring
* Port forwarding
* Tunneling
* Alerting
* Cloud sync
* Daemon mode
* User accounts
* Telemetry

portman must remain local-first and simple.

---

## Final Goal

The final project should allow a developer to run:

```bash
portman
```

and immediately understand:

* what ports are active
* which process owns each port
* which PID is responsible
* which connections are listening or established
* what can be safely killed

The project should also support:

```bash
portman list
portman list --json
portman kill 3000
portman kill node
```

This project must feel like a real open-source developer tool, not a demo.