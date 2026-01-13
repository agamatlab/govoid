# govoid

`govoid` is a lightweight, macOS-focused Terminal User Interface (TUI) application for managing active processes. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, it provides a clean and interactive way to view and terminate running applications.

## Features

- **Process Listing**: Automatically detects and lists visible macOS applications.
- **Interactive UI**: A split-view interface showing active processes and a status/spinner area.
- **Process Management**: Quickly terminate selected applications directly from the TUI.
- **Responsive Design**: Adapts to terminal window resizing.

## Prerequisites

- **OS**: macOS (specifically relies on `lsappinfo` and `pkill`).
- **Go**: version 1.25.5 or later.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/agamatlab/govoid.git
   cd govoid
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

Run the application:
```bash
go run main.go
```

### Controls

| Key | Action |
|-----|--------|
| `tab` | Switch focus between Process List and Status View |
| `p` | Kill the selected application (when Process List is focused) |
| `n` | Next spinner / Reset timer |
| `q` / `ctrl+c` | Quit application |
| `↑`/`↓` | Navigate the process list |

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework.
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions.
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components (list, spinner, timer).

## License

[MIT](LICENSE)
