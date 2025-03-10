# CLAUDE.md - CLI Development Guidelines

## Build Commands

- Build: `make build`
- Run all tests: `make unit-test`
- Run single test: `go test ./path/to/package -run TestName`
- Run smoke tests: `make smoke-test` (requires TEST_EMAIL, TEST_PASSWORD)
- Lint code: `make lint` (install with `make lint-install`)
- Format code: `make pretty` (uses prettier and go fmt)

## Code Style Guidelines

- Follow Go standard naming conventions: camelCase for variables, PascalCase for exports
- Use error handling with appropriate checks - wrap errors with context when needed
- Format imports with standard Go style (stdlib first, then external, then internal)
- Use descriptive variable/function names and maintain consistent indentation
- Add tests for all new functionality and maintain high coverage
- Follow project structure with cmd/ for commands and internal/ for implementation
- Commit messages should be clear and descriptive (feature/fix/chore: message)

## Project-Specific Patterns

- Prefer functional options for configuration
- Use cobra for CLI commands with consistent flags/args pattern
- Use tabwriter for formatted terminal output
