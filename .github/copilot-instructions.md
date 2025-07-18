# GitHub Copilot Instructions for omnistrate-ctl

## Documentation Generation Workflow

**CRITICAL**: After making ANY changes to CLI commands, flags, or help text, you MUST run documentation generation commands to keep docs synchronized.

### Required Commands After CLI Changes:
```bash
make gen-doc    # Regenerates CLI documentation
# OR
make all        # Runs complete build pipeline including doc generation
```

### When to Run Documentation Generation:
- After modifying any file in `cmd/` directory
- After changing command definitions, flags, or help text
- After adding new CLI commands or subcommands
- After updating command usage examples
- After modifying cobra command configurations

### Documentation System:
- Uses `go run doc-gen/main.go` to auto-generate markdown files
- Outputs to `mkdocs/docs/` directory
- Automatically removes old docs and regenerates fresh content
- Documentation must stay synchronized with actual CLI behavior

### Build Pipeline:
The `make all` target runs the complete pipeline:
1. `make tidy` - Clean up Go modules
2. `make build` - Build the binary
3. `make unit-test` - Run tests
4. `make lint` - Code linting
5. `make check-dependencies` - Verify dependencies
6. `make gen-doc` - **Generate documentation**
7. `make pretty` - Format code

## Code Modification Guidelines:
- Always consider if your changes affect CLI behavior
- If modifying anything in `cmd/`, documentation regeneration is required
- Test your changes with `make build` before generating docs
- Use `make all` for comprehensive validation including doc generation

## Quick Reference:
```bash
# After CLI changes - ALWAYS run one of these:
make gen-doc     # Just regenerate docs
make all         # Full pipeline with docs
