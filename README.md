# Azure Command Tower

Azure Command Tower (azct) is a terminal-based UI for exploring and managing Azure resources. Navigate your Azure subscriptions, resource groups, and resources with an intuitive TUI interface.

> **Note**: This project is inspired by [k9s](https://k9scli.io/), the popular Kubernetes terminal UI.

![Status](https://img.shields.io/badge/status-alpha-orange)
![Version](https://img.shields.io/badge/version-0.0.1-blue)
![Go](https://img.shields.io/badge/go-1.23.0-blue)
![License](https://img.shields.io/badge/license-MIT-green)

> **⚠️ ALPHA SOFTWARE - ONGOING DEVELOPMENT**
>
> This project is in **early alpha stage** and under active development. Please be aware:
> - Features may be incomplete or subject to change
> - APIs and command-line interfaces may change between versions
> - Bugs and unexpected behavior should be expected
> - **Not recommended for production use**
> - Your feedback and contributions are highly valued as we shape the project
>
> We're actively working on improving stability and adding features. Check the [changelog](docs/changelog.md) for updates.

## Features

- 🔍 **Browse Azure Resources**: Navigate through subscriptions, resource groups, and resources
- 📦 **Storage Explorer**: Explore Azure Storage accounts, containers, and blobs
- 🔎 **Filter & Search**: Quickly find resources using built-in filtering
- 📊 **Resource Details**: View detailed information about any Azure resource
- 🎨 **Modern TUI**: Beautiful terminal interface built with [tview](https://github.com/rivo/tview)
- ⚡ **Fast & Lightweight**: Native Go application with minimal dependencies

## Prerequisites

- Go 1.23.0 or later
- Azure CLI installed and configured (`az login`)
- Valid Azure subscription(s)

## Installation

### Homebrew (macOS)

Install from tap:
```bash
brew tap rafaelherik/azct
brew install azct
```

Or in one command:
```bash
brew install rafaelherik/azct/azct
```

### From Source

```bash
git clone https://github.com/rafaelherik/azure-control-tower.git
cd azure-control-tower
go build -o azct ./cmd/azct
sudo mv azct /usr/local/bin/
```

### Go Install

```bash
go install github.com/rafaelherik/azure-control-tower/cmd/azct@latest
```

## Quick Start

1. **Authenticate with Azure**:
   ```bash
   az login
   ```

2. **Run azct**:
   ```bash
   azct
   ```

3. **Navigate**:
   - Use arrow keys to navigate
   - Press `Enter` to select/view details
   - Press `/` to filter
   - Press `m` to open resource type menu
   - Press `q` to quit

## Usage

### Navigation

- **Subscriptions View**: Lists all available Azure subscriptions
- **Resource Groups View**: Browse resource groups within a subscription
- **Resource Types View**: See resource type summaries for a resource group
- **Resources View**: View all resources filtered by type
- **Storage Explorer**: Explore storage accounts and containers
- **Blobs View**: Browse blob storage with folder navigation

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate up/down |
| `Enter` | Select/view details |
| `/` | Open filter/search |
| `d` | Show details |
| `e` | Explore storage (for storage accounts) |
| `m` | Open resource type menu |
| `ESC` | Go back |
| `q` | Quit |

### Filtering

Press `/` in any table view to filter resources. The filter is case-insensitive and matches against all visible columns.

## Architecture

```
azct/
├── cmd/azct/           # Main application entry point
├── internal/
│   ├── auth/           # Azure authentication
│   ├── azure/          # Azure SDK client wrappers
│   ├── models/         # Data models
│   ├── navigation/     # Navigation state management
│   └── ui/             # Terminal UI components
└── pkg/
    └── resource/       # Resource handlers and registry
```

## Development

### Building

```bash
go build -o azct ./cmd/azct
```

### Running Tests

```bash
go test ./...
```

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Documentation

Full documentation is available at: [https://rafaelherik.github.io/azure-control-tower](https://rafaelherik.github.io/azure-control-tower)

Or build locally:

```bash
pip install mkdocs mkdocs-material mkdocs-git-revision-date-localized-plugin
mkdocs serve
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [k9s](https://k9scli.io/) - the Kubernetes terminal UI
- Built with [tview](https://github.com/rivo/tview) and [tcell](https://github.com/gdamore/tcell)
- Uses [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go)

## Version

Current version: **0.0.1**

