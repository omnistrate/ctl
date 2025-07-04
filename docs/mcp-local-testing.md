# MCP Server Local Testing Guide

This guide explains how to locally test the MCP (Model Context Protocol) server for omnistrate-ctl.

## Prerequisites

1. **Go 1.21+** - Make sure you have Go installed
2. **VS Code** - For the integrated development experience
3. **Make** - For using the project's build commands

## Quick Start

### 1. List Available Tools

First, see what tools are available:

```bash
# Using make
make build
./dist/omnistrate-ctl-darwin-arm64 mcp --list-tools

# Or using go run
go run . mcp --list-tools
```

### 2. Start the MCP Server

```bash
# Using make
make build
./dist/omnistrate-ctl-darwin-arm64 mcp

# Or using go run
go run . mcp
```

The server will start and listen on stdin/stdout for JSON-RPC messages.

## VS Code Integration

### Launch Configurations

The project includes several VS Code launch configurations:

1. **MCP Server** - Runs the MCP server directly
2. **MCP Server - List Tools** - Lists all available tools
3. **Build and Run MCP Server** - Builds first, then runs the server

To use:

1. Open VS Code in the project root
2. Go to Run and Debug (Ctrl+Shift+D)
3. Select one of the MCP configurations
4. Press F5 to run

### Tasks

Available VS Code tasks (Ctrl+Shift+P → "Tasks: Run Task"):

- **build** - Build the project using make
- **test-mcp** - Run MCP server tests
- **run-mcp-server** - Run the MCP server
- **list-mcp-tools** - List all available MCP tools

## Testing the MCP Server

### 1. Unit Tests

Run the MCP server tests:

```bash
# Using make
make unit-test

# Or directly
go test ./cmd/mcp/server/... -v
```

### 2. Manual Testing

You can manually test the MCP server by sending JSON-RPC requests:

#### List Tools Request

```json
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
```

#### Call Tool Request

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "omnistrate-ctl-account-list",
    "arguments": {
      "flag_output": "json"
    }
  }
}
```

### 3. Testing Script

Create a test script to verify the server:

```bash
#!/bin/bash
# test-mcp.sh

# Build the project
make build

# Test 1: List tools
echo "Testing tools/list..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./dist/omnistrate-ctl-darwin-arm64 mcp

# Test 2: Call a tool (this will fail without authentication, but shows the call works)
echo "Testing tools/call..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"omnistrate-ctl-account-list","arguments":{"flag_output":"json"}}}' | ./dist/omnistrate-ctl-darwin-arm64 mcp
```

## MCP Client Integration

### Claude Desktop

To use with Claude Desktop, add this to your Claude config:

```json
{
  "mcp.servers": {
    "omnistrate-ctl": {
      "command": "/path/to/omnistrate-ctl",
      "args": ["mcp"]
    }
  }
}
```

### VS Code MCP Extension

If you have an MCP extension for VS Code, you can configure it to use the local server:

```json
{
  "mcp.servers": {
    "omnistrate-ctl": {
      "command": "go",
      "args": ["run", ".", "mcp"],
      "env": {
        "PATH": "${env:PATH}:${workspaceFolder}/dist"
      }
    }
  }
}
```

## Development Workflow

1. **Make changes** to the MCP server code
2. **Run tests** to ensure everything works:

   ```bash
   make unit-test
   ```

3. **Test locally** using the VS Code configurations or manual testing
4. **Build and deploy** when ready:

   ```bash
   make build
   ```

## Troubleshooting

### Common Issues

1. **Command not found**: Make sure the omnistrate-ctl binary is built and in your PATH
2. **Authentication errors**: The MCP server inherits the same authentication as the CLI
3. **Tool not found**: Check that the tool name matches exactly what's returned by `--list-tools`

### Debug Mode

Add debug logging to the server by modifying the log level:

```go
// In cmd/mcp/server/server.go
log.SetLevel(log.DebugLevel)
```

### Logs

The MCP server logs to stderr, so you can redirect logs while keeping the JSON-RPC communication on stdout:

```bash
go run . mcp 2>debug.log
```

## Adding New Tools

When you add new commands to omnistrate-ctl, they automatically become available as MCP tools. The server uses reflection to discover all available commands and their flags.

To test new tools:

1. Add your new command
2. Run `go run . mcp --list-tools` to verify it appears
3. Test the tool using the JSON-RPC interface
