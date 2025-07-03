package server

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllTools(t *testing.T) {
	// Create a simple test command structure
	rootCmd := &cobra.Command{
		Use: "test-ctl",
	}
	
	accountCmd := &cobra.Command{
		Use:   "account",
		Short: "Manage accounts",
	}
	
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create account",
	}
	
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
	}
	
	// Add flags to test flag parsing
	createCmd.Flags().String("name", "", "Account name")
	createCmd.Flags().Bool("enabled", false, "Enable account")
	
	accountCmd.AddCommand(createCmd)
	accountCmd.AddCommand(listCmd)
	rootCmd.AddCommand(accountCmd)
	
	// Test tool generation
	tools, err := GetAllTools(rootCmd)
	require.NoError(t, err)
	
	// Should have 2 tools (create and list)
	assert.Len(t, tools, 2)
	
	// Check the tools
	var createTool, listTool *Tool
	for _, tool := range tools {
		if tool.Name == "omnistrate-ctl-test-ctl-account-create" {
			createTool = &tool
		} else if tool.Name == "omnistrate-ctl-test-ctl-account-list" {
			listTool = &tool
		}
	}
	
	require.NotNil(t, createTool, "create tool should exist")
	require.NotNil(t, listTool, "list tool should exist")
	
	// Check create tool properties
	assert.Equal(t, "Create account", createTool.Description)
	assert.Equal(t, "object", createTool.InputSchema.Type)
	assert.Contains(t, createTool.InputSchema.Properties, "flag_name")
	assert.Contains(t, createTool.InputSchema.Properties, "flag_enabled")
	
	// Check flag types
	assert.Equal(t, "string", createTool.InputSchema.Properties["flag_name"].Type)
	assert.Equal(t, "boolean", createTool.InputSchema.Properties["flag_enabled"].Type)
	
	// Check list tool
	assert.Equal(t, "List accounts", listTool.Description)
}

func TestMCPRequestResponse(t *testing.T) {
	// Create a simple test command structure
	rootCmd := &cobra.Command{
		Use: "test-ctl",
	}
	
	pingCmd := &cobra.Command{
		Use:   "ping",
		Short: "Test ping command",
	}
	
	rootCmd.AddCommand(pingCmd)
	
	// Create MCP server
	server, err := NewMCPServer(rootCmd)
	require.NoError(t, err)
	
	// Test tools/list request
	listRequest := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "tools/list",
		Params:  json.RawMessage(`{}`),
	}
	
	response := server.handleRequest(listRequest)
	
	assert.Equal(t, "2.0", response.Jsonrpc)
	assert.Equal(t, 1, response.ID)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Result)
	
	// Check that result contains tools
	result, ok := response.Result.(map[string]interface{})
	assert.True(t, ok)
	assert.Contains(t, result, "tools")
	
	tools, ok := result["tools"].([]Tool)
	assert.True(t, ok)
	assert.Len(t, tools, 1)
	assert.Equal(t, "omnistrate-ctl-test-ctl-ping", tools[0].Name)
}

func TestMCPErrorHandling(t *testing.T) {
	rootCmd := &cobra.Command{
		Use: "test-ctl",
	}
	
	server, err := NewMCPServer(rootCmd)
	require.NoError(t, err)
	
	// Test unknown method
	request := MCPRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "unknown/method",
		Params:  json.RawMessage(`{}`),
	}
	
	response := server.handleRequest(request)
	
	assert.Equal(t, "2.0", response.Jsonrpc)
	assert.Equal(t, 1, response.ID)
	assert.Nil(t, response.Result)
	assert.NotNil(t, response.Error)
	assert.Equal(t, -32601, response.Error.Code)
	assert.Contains(t, response.Error.Message, "Method not found")
}