# Contributing

Thank you for your interest in contributing to Azure Command Tower!

## Project Status

Azure Command Tower is currently in **alpha stage**. This means:

- **Core features are still being developed** - Some functionality may be incomplete or under active refinement
- **Breaking changes may occur** - APIs and CLI interfaces might change between versions as we improve the design
- **Your contributions are especially valuable** - Early feedback and contributions help shape the project's direction
- **Focus areas for contribution**:
  - Bug fixes and stability improvements
  - Core functionality enhancements
  - Documentation improvements
  - Testing and test coverage
  - Performance optimizations

Please understand that during this alpha phase:
- Some features may be partially implemented
- Code structure and architecture may evolve
- Your patience and constructive feedback are greatly appreciated

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/rafaelherik/azure-control-tower.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes: `go test ./...`
6. Commit: `git commit -m 'Add some feature'`
7. Push: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

1. Install Go 1.23.0 or later
2. Install dependencies: `go mod download`
3. Build: `go build -o azct ./cmd/azct`
4. Run: `./azct`

## Code Style

- Follow Go conventions
- Use `go fmt` to format code
- Write clear, descriptive names
- Add comments for exported functions
- Keep functions focused and small

## Testing

- Write tests for new features
- Ensure all tests pass: `go test ./...`
- Aim for good test coverage

## Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Submit PR with clear description

## Areas for Contribution

- Bug fixes
- New features
- Documentation improvements
- Performance optimizations
- UI/UX improvements
- Resource handler implementations
- Test coverage

## Questions?

Open an issue for questions or discussions about contributions.

