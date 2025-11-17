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

## Reporting Issues

Before you start coding, please check if there's already an issue for what you want to work on. When creating a new issue, please use the appropriate issue template:

- **Bug Report**: Report bugs or unexpected behavior
- **Feature Request**: Suggest new features or enhancements
- **Question**: Ask questions or get help using azct

Using these templates ensures we have all the necessary information to understand and address your concern quickly.

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

## Reporting Bugs

If you find a bug, please report it using the **Bug Report** issue template. Include:

- Clear description of the bug
- Steps to reproduce
- Expected vs. actual behavior
- Your environment (OS, azct version, Azure CLI version)
- Screenshots or logs if applicable

## Requesting Features

Have an idea for a new feature? Use the **Feature Request** issue template to:

- Describe the problem you're trying to solve
- Propose your solution
- Explain the use case and benefits
- Provide any additional context

## Questions?

If you have questions or need help:

- Check the [documentation](https://rafaelherik.github.io/azure-control-tower) first
- Use the **Question** issue template for support
- Join discussions in [GitHub Discussions](https://github.com/rafaelherik/azure-control-tower/discussions)

