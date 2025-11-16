# Building

Instructions for building Azure Command Tower from source.

## Prerequisites

- Go 1.23.0 or later
- Git
- Azure CLI (for testing)

## Building

### Basic Build

```bash
go build -o azct ./cmd/azct
```

This creates an `azct` binary in the current directory.

### Cross-Platform Builds

Build for different platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o azct-linux-amd64 ./cmd/azct

# macOS
GOOS=darwin GOARCH=amd64 go build -o azct-darwin-amd64 ./cmd/azct
GOOS=darwin GOARCH=arm64 go build -o azct-darwin-arm64 ./cmd/azct

# Windows
GOOS=windows GOARCH=amd64 go build -o azct-windows-amd64.exe ./cmd/azct
```

### Release Build

For production releases, use build flags:

```bash
go build -ldflags="-s -w -X main.version=0.0.1" -o azct ./cmd/azct
```

This:
- Strips debug symbols (`-s`)
- Disables DWARF generation (`-w`)
- Sets version information

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Development Workflow

1. Make your changes
2. Run tests: `go test ./...`
3. Build: `go build -o azct ./cmd/azct`
4. Test the binary: `./azct`
5. Commit changes

## Dependencies

Update dependencies:

```bash
go get -u ./...
go mod tidy
```

## Code Formatting

Format code:

```bash
go fmt ./...
```

## Linting

Use `golangci-lint` or similar:

```bash
golangci-lint run
```

