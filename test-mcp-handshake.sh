#!/bin/bash

# Test MCP Server Handshake
echo "Testing MCP Server Handshake..."

# Build the project first
echo "Building project..."
make build

# Test the complete MCP handshake
echo -e "\n1. Testing initialize method..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | \
./dist/omnistrate-ctl-darwin-arm64 mcp | jq .

echo -e "\n2. Testing initialized notification..."
echo '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}' | \
./dist/omnistrate-ctl-darwin-arm64 mcp

echo -e "\n3. Testing tools/list..."
echo '{"jsonrpc":"2.0","id":3,"method":"tools/list","params":{}}' | \
./dist/omnistrate-ctl-darwin-arm64 mcp 2>/dev/null | jq '.result.tools | length'

echo -e "\n4. Testing invalid method..."
echo '{"jsonrpc":"2.0","id":4,"method":"invalid/method","params":{}}' | \
./dist/omnistrate-ctl-darwin-arm64 mcp 2>/dev/null | jq .

echo -e "\nMCP Server handshake test complete!"
