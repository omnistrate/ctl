package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MCPTestClient provides functionality to test the MCP server as a subprocess
type MCPTestClient struct {
	process   *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	stderr    io.ReadCloser
	scanner   *bufio.Scanner
	requestID int
}

// NewMCPTestClient creates a new MCP test client
func NewMCPTestClient(t *testing.T) *MCPTestClient {
	// Build the binary path based on the current platform
	binaryPath := getBinaryPath(t)
	
	// Create the command with --no-auth flag for testing
	cmd := exec.Command(binaryPath, "mcp", "--no-auth")
	
	// Set up pipes
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	
	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)
	
	client := &MCPTestClient{
		process: cmd,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		scanner: bufio.NewScanner(stdout),
	}
	
	return client
}

// Start starts the MCP server process
func (c *MCPTestClient) Start(t *testing.T) {
	err := c.process.Start()
	require.NoError(t, err)
	
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
}

// Stop stops the MCP server process
func (c *MCPTestClient) Stop(t *testing.T) {
	if c.process != nil && c.process.Process != nil {
		c.stdin.Close()
		c.process.Process.Kill()
		c.process.Wait()
	}
}

// SendRequest sends a JSON-RPC request and returns the response
func (c *MCPTestClient) SendRequest(t *testing.T, method string, params interface{}) map[string]interface{} {
	c.requestID++
	
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      c.requestID,
		"method":  method,
		"params":  params,
	}
	
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)
	
	// Send request
	_, err = c.stdin.Write(append(requestBytes, '\n'))
	require.NoError(t, err)
	
	// Read response
	require.True(t, c.scanner.Scan(), "Failed to read response")
	responseText := c.scanner.Text()
	
	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseText), &response)
	require.NoError(t, err)
	
	return response
}

// SendNotification sends a JSON-RPC notification (no response expected)
func (c *MCPTestClient) SendNotification(t *testing.T, method string, params interface{}) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}
	
	requestBytes, err := json.Marshal(request)
	require.NoError(t, err)
	
	// Send notification
	_, err = c.stdin.Write(append(requestBytes, '\n'))
	require.NoError(t, err)
}

// getBinaryPath returns the path to the built binary for the current platform
func getBinaryPath(t *testing.T) string {
	// Get the project root directory
	projectRoot := getProjectRoot(t)
	
	// Determine the binary name based on platform
	system := runtime.GOOS
	arch := runtime.GOARCH
	
	binaryName := fmt.Sprintf("omnistrate-ctl-%s-%s", system, arch)
	if system == "windows" {
		binaryName += ".exe"
	}
	
	binaryPath := filepath.Join(projectRoot, "dist", binaryName)
	
	// Check if binary exists, if not try to build it
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Logf("Binary not found at %s, attempting to build...", binaryPath)
		buildCmd := exec.Command("make", "build")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Failed to build binary")
	}
	
	return binaryPath
}

// getProjectRoot finds the project root directory
func getProjectRoot(t *testing.T) string {
	// Start from current directory and walk up to find go.mod
	dir, err := os.Getwd()
	require.NoError(t, err)
	
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	t.Fatal("Could not find project root (go.mod)")
	return ""
}

// TestMCPServerIntegration tests the MCP server as a subprocess
func TestMCPServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := NewMCPTestClient(t)
	client.Start(t)
	defer client.Stop(t)
	
	t.Run("Initialize", func(t *testing.T) {
		// Test initialize method
		initParams := map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		}
		
		response := client.SendRequest(t, "initialize", initParams)
		
		assert.Equal(t, "2.0", response["jsonrpc"])
		assert.Equal(t, 1, int(response["id"].(float64)))
		assert.Nil(t, response["error"])
		assert.NotNil(t, response["result"])
		
		// Check result structure
		result := response["result"].(map[string]interface{})
		assert.Contains(t, result, "protocolVersion")
		assert.Contains(t, result, "capabilities")
		assert.Contains(t, result, "serverInfo")
	})
	
	t.Run("Initialized", func(t *testing.T) {
		// Send initialized notification
		client.SendNotification(t, "notifications/initialized", map[string]interface{}{})
		// No response expected for notifications
	})
	
	t.Run("ToolsList", func(t *testing.T) {
		// Test tools/list method
		response := client.SendRequest(t, "tools/list", map[string]interface{}{})
		
		assert.Equal(t, "2.0", response["jsonrpc"])
		assert.Nil(t, response["error"])
		assert.NotNil(t, response["result"])
		
		// Check result structure
		result := response["result"].(map[string]interface{})
		assert.Contains(t, result, "tools")
		
		tools := result["tools"].([]interface{})
		assert.Greater(t, len(tools), 0, "Should have at least one tool")
		
		// Check first tool structure
		firstTool := tools[0].(map[string]interface{})
		assert.Contains(t, firstTool, "name")
		assert.Contains(t, firstTool, "description")
		assert.Contains(t, firstTool, "inputSchema")
		
		// Tool name should start with "omnistrate-ctl-"
		toolName := firstTool["name"].(string)
		assert.True(t, strings.HasPrefix(toolName, "omnistrate-ctl-"), "Tool name should start with 'omnistrate-ctl-'")
		
		t.Logf("Found %d tools, first tool: %s", len(tools), toolName)
	})
	
	t.Run("InvalidMethod", func(t *testing.T) {
		// Test invalid method
		response := client.SendRequest(t, "invalid/method", map[string]interface{}{})
		
		assert.Equal(t, "2.0", response["jsonrpc"])
		assert.Nil(t, response["result"])
		assert.NotNil(t, response["error"])
		
		// Check error structure
		errorObj := response["error"].(map[string]interface{})
		assert.Equal(t, float64(-32601), errorObj["code"])
		assert.Contains(t, errorObj["message"], "Method not found")
	})
	
	t.Run("ToolsCall", func(t *testing.T) {
		// First get the list of tools to find one we can call
		toolsResponse := client.SendRequest(t, "tools/list", map[string]interface{}{})
		result := toolsResponse["result"].(map[string]interface{})
		tools := result["tools"].([]interface{})
		
		require.Greater(t, len(tools), 0, "Should have at least one tool")
		
		// Find a simple tool to call (prefer help or version commands)
		var toolToCall string
		for _, tool := range tools {
			toolMap := tool.(map[string]interface{})
			name := toolMap["name"].(string)
			if strings.Contains(name, "help") || strings.Contains(name, "version") {
				toolToCall = name
				break
			}
		}
		
		// If no help/version tool found, use the first tool
		if toolToCall == "" {
			toolToCall = tools[0].(map[string]interface{})["name"].(string)
		}
		
		// Test tools/call method
		callParams := map[string]interface{}{
			"name":      toolToCall,
			"arguments": map[string]interface{}{},
		}
		
		response := client.SendRequest(t, "tools/call", callParams)
		
		assert.Equal(t, "2.0", response["jsonrpc"])
		
		// Check if there's an error (tool might fail due to missing arguments)
		if response["error"] != nil {
			errorObj := response["error"].(map[string]interface{})
			t.Logf("Tool call failed (expected): %v", errorObj["message"])
			return
		}
		
		// If no error, check result structure
		if response["result"] != nil {
			result := response["result"].(map[string]interface{})
			assert.Contains(t, result, "content")
			
			content := result["content"].([]interface{})
			assert.Greater(t, len(content), 0, "Should have at least one content item")
			
			// Check first content item
			firstContent := content[0].(map[string]interface{})
			assert.Contains(t, firstContent, "type")
			assert.Contains(t, firstContent, "text")
			assert.Equal(t, "text", firstContent["type"])
			
			t.Logf("Called tool %s successfully", toolToCall)
		}
	})
	
	t.Run("ToolsCallInvalid", func(t *testing.T) {
		// Test tools/call with invalid tool name
		callParams := map[string]interface{}{
			"name":      "non-existent-tool",
			"arguments": map[string]interface{}{},
		}
		
		response := client.SendRequest(t, "tools/call", callParams)
		
		assert.Equal(t, "2.0", response["jsonrpc"])
		assert.Nil(t, response["result"])
		assert.NotNil(t, response["error"])
		
		// Check error structure
		errorObj := response["error"].(map[string]interface{})
		assert.Equal(t, float64(-32602), errorObj["code"])
		assert.Contains(t, errorObj["message"], "tool not found")
	})
}

// TestMCPServerHandshake tests the complete MCP handshake process
func TestMCPServerHandshake(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := NewMCPTestClient(t)
	client.Start(t)
	defer client.Stop(t)
	
	// Step 1: Initialize
	initParams := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "test-client",
			"version": "1.0.0",
		},
	}
	
	initResponse := client.SendRequest(t, "initialize", initParams)
	assert.Nil(t, initResponse["error"])
	assert.NotNil(t, initResponse["result"])
	
	// Step 2: Send initialized notification
	client.SendNotification(t, "notifications/initialized", map[string]interface{}{})
	
	// Step 3: Use the server (list tools)
	toolsResponse := client.SendRequest(t, "tools/list", map[string]interface{}{})
	assert.Nil(t, toolsResponse["error"])
	assert.NotNil(t, toolsResponse["result"])
	
	result := toolsResponse["result"].(map[string]interface{})
	tools := result["tools"].([]interface{})
	
	t.Logf("Handshake completed successfully, server reported %d tools", len(tools))
	assert.Greater(t, len(tools), 0, "Should have at least one tool after handshake")
}

// TestMCPServerToolCount tests that the server exposes the expected number of tools
func TestMCPServerToolCount(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := NewMCPTestClient(t)
	client.Start(t)
	defer client.Stop(t)
	
	// Get tools list
	response := client.SendRequest(t, "tools/list", map[string]interface{}{})
	require.Nil(t, response["error"])
	require.NotNil(t, response["result"])
	
	result := response["result"].(map[string]interface{})
	tools := result["tools"].([]interface{})
	
	// We expect a significant number of tools (74+ as mentioned in the PR description)
	assert.Greater(t, len(tools), 50, "Should have more than 50 tools")
	
	// Check that all tools follow the naming convention
	for _, tool := range tools {
		toolMap := tool.(map[string]interface{})
		name := toolMap["name"].(string)
		assert.True(t, strings.HasPrefix(name, "omnistrate-ctl-"), 
			"Tool name '%s' should start with 'omnistrate-ctl-'", name)
	}
	
	t.Logf("Server exposes %d tools", len(tools))
}