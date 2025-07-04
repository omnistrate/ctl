#!/bin/bash
# test-mcp.sh - Simple script to test the MCP server locally

set -e

echo "🚀 Building omnistrate-ctl..."
make build

echo ""
echo "📋 Listing available MCP tools..."
./dist/omnistrate-ctl-darwin-arm64 mcp --list-tools | head -20

echo ""
echo "🔧 Testing MCP server with JSON-RPC requests..."
echo "Starting server in background..."

# Start the server in the background
./dist/omnistrate-ctl-darwin-arm64 mcp &
SERVER_PID=$!

# Give the server a moment to start
sleep 1

# Test 1: List tools
echo ""
echo "📡 Testing tools/list request..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./dist/omnistrate-ctl-darwin-arm64 mcp | head -10

# Test 2: Test an invalid method
echo ""
echo "❌ Testing invalid method (should return error)..."
echo '{"jsonrpc":"2.0","id":2,"method":"invalid/method","params":{}}' | ./dist/omnistrate-ctl-darwin-arm64 mcp

# Clean up
kill $SERVER_PID 2>/dev/null || true

echo ""
echo "✅ MCP server tests completed!"
echo ""
echo "💡 To test interactively:"
echo "   1. Run: make build && ./dist/omnistrate-ctl-darwin-arm64 mcp"
echo "   2. Send JSON-RPC requests via stdin"
echo "   3. Or use VS Code launch configurations"
