# Installation

Azure Command Tower (azct) can be installed using several methods.

## Prerequisites

- Go 1.23.0 or later (for building from source)
- Azure CLI installed and configured
- Valid Azure subscription(s)

## Homebrew (macOS)

```bash
brew install azct
```

## From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/azure-control-tower.git
   cd azure-control-tower
   ```

2. Build the application:
   ```bash
   go build -o azct ./cmd/azct
   ```

3. Install to your PATH:
   ```bash
   sudo mv azct /usr/local/bin/
   ```

## Go Install

```bash
go install github.com/yourusername/azure-control-tower/cmd/azct@latest
```

## Verify Installation

After installation, verify that azct is working:

```bash
azct --version
```

You should see the version number (0.0.1).

## Next Steps

Once installed, proceed to [Quick Start](quick-start.md) to learn how to use azct.

