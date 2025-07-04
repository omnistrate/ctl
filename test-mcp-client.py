#!/usr/bin/env python3
"""
Interactive MCP client for testing the omnistrate-ctl MCP server.
This script provides a simple way to test JSON-RPC calls to the MCP server.
"""

import json
import subprocess
import sys
import threading
import time
from typing import Optional

class MCPClient:
    def __init__(self, server_command: list):
        self.server_command = server_command
        self.process: Optional[subprocess.Popen] = None
        self.request_id = 0
        
    def start_server(self):
        """Start the MCP server process."""
        try:
            self.process = subprocess.Popen(
                self.server_command,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                bufsize=1
            )
            print("✅ MCP server started")
            return True
        except Exception as e:
            print(f"❌ Failed to start server: {e}")
            return False
    
    def send_request(self, method: str, params: dict = None) -> dict:
        """Send a JSON-RPC request to the server."""
        if not self.process:
            raise RuntimeError("Server not started")
        
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        
        request_json = json.dumps(request)
        print(f"📤 Sending: {request_json}")
        
        # Send request
        self.process.stdin.write(request_json + "\n")
        self.process.stdin.flush()
        
        # Read response
        response_line = self.process.stdout.readline()
        if response_line:
            try:
                response = json.loads(response_line.strip())
                print(f"📥 Response: {json.dumps(response, indent=2)}")
                return response
            except json.JSONDecodeError as e:
                print(f"❌ Invalid JSON response: {e}")
                print(f"Raw response: {response_line}")
                return {}
        else:
            print("❌ No response received")
            return {}
    
    def stop_server(self):
        """Stop the MCP server process."""
        if self.process:
            self.process.terminate()
            self.process.wait()
            print("🛑 Server stopped")

def main():
    print("🚀 MCP Client for omnistrate-ctl")
    print("=" * 50)
    
    # Build the server first
    print("Building omnistrate-ctl...")
    try:
        subprocess.run(["make", "build"], check=True, capture_output=True)
        print("✅ Build successful")
    except subprocess.CalledProcessError as e:
        print(f"❌ Build failed: {e}")
        sys.exit(1)
    
    # Determine the binary name based on platform
    import platform
    system = platform.system().lower()
    machine = platform.machine().lower()
    
    if machine == "x86_64":
        machine = "amd64"
    elif machine in ["aarch64", "arm64"]:
        machine = "arm64"
    
    binary_name = f"omnistrate-ctl-{system}-{machine}"
    if system == "windows":
        binary_name += ".exe"
    
    server_command = [f"./dist/{binary_name}", "mcp"]
    
    # Create client
    client = MCPClient(server_command)
    
    if not client.start_server():
        sys.exit(1)
    
    try:
        # Wait a moment for server to start
        time.sleep(0.5)
        
        print("\n🔧 Testing MCP server...")
        
        # Test 1: List tools
        print("\n1. Testing tools/list")
        response = client.send_request("tools/list")
        
        if response.get("result"):
            tools = response["result"].get("tools", [])
            print(f"✅ Found {len(tools)} tools")
            if tools:
                print("📋 Available tools:")
                for i, tool in enumerate(tools[:5]):  # Show first 5 tools
                    print(f"   {i+1}. {tool['name']}: {tool['description']}")
                if len(tools) > 5:
                    print(f"   ... and {len(tools) - 5} more")
        
        # Test 2: Test invalid method
        print("\n2. Testing invalid method")
        response = client.send_request("invalid/method")
        
        if response.get("error"):
            print("✅ Error handling works correctly")
        
        # Test 3: Interactive mode
        print("\n🎮 Interactive mode (type 'quit' to exit)")
        print("Available commands:")
        print("  tools/list - List all available tools")
        print("  tools/call - Call a tool (you'll need to provide tool name and arguments)")
        print("  quit - Exit the client")
        
        while True:
            try:
                user_input = input("\n> ").strip()
                
                if user_input.lower() in ['quit', 'exit']:
                    break
                
                if user_input == "tools/list":
                    client.send_request("tools/list")
                elif user_input.startswith("tools/call"):
                    # Simple tool call example
                    print("Example tool call:")
                    print('{"name": "omnistrate-ctl-account-list", "arguments": {"flag_output": "json"}}')
                    tool_params = input("Enter tool parameters (JSON): ").strip()
                    try:
                        params = json.loads(tool_params)
                        client.send_request("tools/call", params)
                    except json.JSONDecodeError:
                        print("❌ Invalid JSON parameters")
                elif user_input:
                    # Try to parse as JSON-RPC method
                    try:
                        if user_input.startswith("{"):
                            # Full JSON-RPC request
                            request = json.loads(user_input)
                            method = request.get("method")
                            params = request.get("params")
                            client.send_request(method, params)
                        else:
                            # Simple method name
                            client.send_request(user_input)
                    except json.JSONDecodeError:
                        print("❌ Invalid JSON or method name")
                    except Exception as e:
                        print(f"❌ Error: {e}")
                
            except KeyboardInterrupt:
                print("\n\n👋 Goodbye!")
                break
            except EOFError:
                print("\n\n👋 Goodbye!")
                break
    
    finally:
        client.stop_server()

if __name__ == "__main__":
    main()
