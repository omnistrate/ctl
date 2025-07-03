# MCP Server for Omnistrate CTL

This document describes the Model Context Protocol (MCP) server implementation for the omnistrate-ctl tool.

## Overview

The MCP server exposes all omnistrate-ctl commands as tools that can be used by AI models following the Model Context Protocol specification. This allows AI assistants to interact with the Omnistrate platform through a standardized protocol.

## Usage

### Start the MCP Server

```bash
omnistrate-ctl mcp
```

The server will start and listen for MCP requests on stdin, sending responses to stdout. This follows the standard MCP JSON-RPC protocol.

### List Available Tools

```bash
omnistrate-ctl mcp --list-tools
```

This command lists all available tools in JSON format, showing the complete tool definitions including parameters and descriptions.

## Tool Structure

Each CLI command is exposed as an MCP tool with the following structure:

- **Name**: `omnistrate-ctl-<command-hierarchy>` (e.g., `omnistrate-ctl-instance-create`)
- **Description**: The command's short or long description
- **Input Schema**: JSON schema defining the tool's parameters

### Tool Naming Convention

Tools are named using the full command hierarchy:

- `omnistrate-ctl-account-create` for `omnistrate-ctl account create`
- `omnistrate-ctl-instance-list` for `omnistrate-ctl instance list`
- `omnistrate-ctl-service-plan-release` for `omnistrate-ctl service-plan release`

### Parameter Mapping

CLI flags are mapped to tool parameters with the `flag_` prefix:

- `--output json` becomes `"flag_output": "json"`
- `--environment prod` becomes `"flag_environment": "prod"`
- `--enabled` becomes `"flag_enabled": true`

## Available Tools

The MCP server exposes 74+ tools covering all CLI functionality:

### Account Management

- `omnistrate-ctl-account-create` - Create a Cloud Provider Account
- `omnistrate-ctl-account-delete` - Delete a Cloud Provider Account
- `omnistrate-ctl-account-describe` - Describe a Cloud Provider Account
- `omnistrate-ctl-account-list` - List Cloud Provider Accounts

### Service Building

- `omnistrate-ctl-build` - Build Services from image, compose spec or service plan spec
- `omnistrate-ctl-build-from-repo` - Build Service from Git Repository

### Instance Management

- `omnistrate-ctl-instance-create` - Create an instance deployment
- `omnistrate-ctl-instance-delete` - Delete an instance deployment
- `omnistrate-ctl-instance-describe` - Describe an instance deployment
- `omnistrate-ctl-instance-list` - List instance deployments
- `omnistrate-ctl-instance-start` - Start an instance deployment
- `omnistrate-ctl-instance-stop` - Stop an instance deployment
- `omnistrate-ctl-instance-restart` - Restart an instance deployment
- `omnistrate-ctl-instance-modify` - Modify an instance deployment
- And many more...

### Service Plan Management

- `omnistrate-ctl-service-plan-create` - Create a service plan
- `omnistrate-ctl-service-plan-delete` - Delete a service plan
- `omnistrate-ctl-service-plan-list` - List service plans
- `omnistrate-ctl-service-plan-release` - Release a service plan
- And more...

[Complete list available via `omnistrate-ctl mcp --list-tools`]

## MCP Protocol Support

The server implements the following MCP methods:

### `tools/list`

Returns all available tools with their schemas.

**Request:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/list",
  "params": {}
}
```

**Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "omnistrate-ctl-instance-create",
        "description": "Create an instance deployment for your service",
        "inputSchema": {
          "type": "object",
          "properties": {
            "flag_environment": {
              "type": "string",
              "description": "Environment name"
            },
            "flag_cloud-provider": {
              "type": "string",
              "description": "Cloud provider (aws|gcp)"
            }
          },
          "required": []
        }
      }
    ]
  }
}
```

### `tools/call`

Executes a tool (CLI command).

**Request:**

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "omnistrate-ctl-instance-list",
    "arguments": {
      "flag_output": "json"
    }
  }
}
```

**Response:**

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "[CLI command output here]"
      }
    ]
  }
}
```

## Error Handling

The server returns standard JSON-RPC errors:

- `-32601`: Method not found
- `-32602`: Invalid params (e.g., tool not found, invalid arguments)

## Integration Examples

### Claude Desktop

Add to your Claude Desktop configuration:

```json
{
  "mcpServers": {
    "omnistrate-ctl": {
      "command": "omnistrate-ctl",
      "args": ["mcp"]
    }
  }
}
```

### Other MCP Clients

The server works with any MCP-compatible client. Start the server with:

```bash
omnistrate-ctl mcp
```

Then communicate using the JSON-RPC protocol over stdin/stdout.

## Development

The MCP server is implemented in Go and automatically discovers all CLI commands using reflection on the Cobra command structure. When new commands are added to the CLI, they automatically become available as MCP tools.

### Testing

Run the MCP server tests:

```bash
go test ./cmd/mcp/server/...
```

### Building

The MCP server is built as part of the main omnistrate-ctl binary:

```bash
make build
```
