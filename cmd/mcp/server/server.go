package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// MCPServer represents the Model Context Protocol server
type MCPServer struct {
	tools   []Tool
	rootCmd *cobra.Command
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema represents the JSON schema for tool input
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

// Property represents a property in the input schema
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// MCPRequest represents an MCP JSON-RPC request
type MCPRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// MCPResponse represents an MCP JSON-RPC response
type MCPResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewMCPServer creates a new MCP server
func NewMCPServer(rootCmd *cobra.Command) (*MCPServer, error) {
	tools, err := GetAllTools(rootCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools: %w", err)
	}

	return &MCPServer{
		tools:   tools,
		rootCmd: rootCmd,
	}, nil
}

// Start starts the MCP server
func (s *MCPServer) Start() error {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			log.Printf("Failed to parse request: %v", err)
			continue
		}

		response := s.handleRequest(request)

		// Skip sending response for notifications
		if response.Jsonrpc == "" {
			continue
		}

		responseData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			continue
		}

		fmt.Println(string(responseData))
	}

	return scanner.Err()
}

// handleRequest handles an MCP request
func (s *MCPServer) handleRequest(request MCPRequest) MCPResponse {
	response := MCPResponse{
		Jsonrpc: "2.0",
		ID:      request.ID,
	}

	switch request.Method {
	case "initialize":
		// Handle MCP initialization
		capabilities := s.generateCapabilities()
		response.Result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    capabilities,
			"serverInfo": map[string]interface{}{
				"name":    "omnistrate-ctl",
				"version": "1.0.0",
			},
		}
	case "notifications/initialized":
		// Handle initialized notification - no response needed
		return MCPResponse{}
	case "tools/list":
		response.Result = map[string]interface{}{
			"tools": s.tools,
		}
	case "tools/call":
		result, err := s.callTool(request.Params)
		if err != nil {
			response.Error = &MCPError{
				Code:    -32602,
				Message: err.Error(),
			}
		} else {
			response.Result = result
		}
	default:
		response.Error = &MCPError{
			Code:    -32601,
			Message: fmt.Sprintf("Method not found: %s", request.Method),
		}
	}

	return response
}

// callTool executes a tool
func (s *MCPServer) callTool(params json.RawMessage) (interface{}, error) {
	var toolCall struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(params, &toolCall); err != nil {
		return nil, fmt.Errorf("failed to parse tool call: %w", err)
	}

	// Find the tool
	var tool *Tool
	for _, t := range s.tools {
		if t.Name == toolCall.Name {
			tool = &t
			break
		}
	}

	if tool == nil {
		return nil, fmt.Errorf("tool not found: %s", toolCall.Name)
	}

	// Execute the command
	return s.executeCommand(toolCall.Name, toolCall.Arguments)
}

// executeCommand executes a CLI command
func (s *MCPServer) executeCommand(toolName string, arguments map[string]interface{}) (interface{}, error) {
	// Build command arguments
	args := []string{}

	// Parse tool name to get command parts
	parts := strings.Split(toolName, "-")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid tool name: %s", toolName)
	}

	// Skip "omnistrate" prefix and "ctl" suffix
	if parts[0] == "omnistrate" {
		parts = parts[1:]
	}
	if len(parts) > 0 && parts[len(parts)-1] == "ctl" {
		parts = parts[:len(parts)-1]
	}

	// Add command parts
	args = append(args, parts...)

	// Add arguments
	for key, value := range arguments {
		if key == "output" {
			args = append(args, "--output", fmt.Sprintf("%v", value))
		} else if strings.HasPrefix(key, "flag_") {
			// Handle flags
			flagName := strings.TrimPrefix(key, "flag_")
			if flagName != "" {
				args = append(args, "--"+flagName, fmt.Sprintf("%v", value))
			}
		} else {
			// Handle positional arguments
			args = append(args, fmt.Sprintf("%v", value))
		}
	}

	// Execute the command
	cmd := exec.Command("omnistrate-ctl", args...)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w, output: %s", err, string(output))
	}

	return map[string]interface{}{
		"content": []*TextContent{{
			Type: "text",
			Text: string(output),
		}},
	}, nil
}

// TextContent represents text content in MCP response
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// GetAllTools gets all available tools from the CLI
func GetAllTools(rootCmd *cobra.Command) ([]Tool, error) {
	tools := []Tool{}

	// Add a tool for each command
	err := walkCommands(rootCmd, "", &tools)
	if err != nil {
		return nil, err
	}

	return tools, nil
}

// walkCommands recursively walks through all commands
func walkCommands(cmd *cobra.Command, prefix string, tools *[]Tool) error {
	// Skip root command and parent commands that have subcommands
	if cmd.Use != "omnistrate-ctl" && !cmd.HasSubCommands() {
		toolName := prefix + cmd.Use

		// Clean up tool name (remove everything after space)
		if strings.Contains(toolName, " ") {
			toolName = strings.Split(toolName, " ")[0]
		}

		// Create tool
		tool := Tool{
			Name:        "omnistrate-ctl-" + strings.ReplaceAll(toolName, " ", "-"),
			Description: cmd.Short,
			InputSchema: InputSchema{
				Type:       "object",
				Properties: make(map[string]Property),
				Required:   []string{},
			},
		}

		// Add long description if available
		if cmd.Long != "" {
			tool.Description = cmd.Long
		}

		// Add flags as properties
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if flag.Name == "help" || flag.Name == "version" {
				return
			}

			property := Property{
				Type:        "string",
				Description: flag.Usage,
			}

			// Handle boolean flags
			if flag.Value.Type() == "bool" {
				property.Type = "boolean"
			}

			// Handle enum flags (like output format)
			if flag.Name == "output" {
				property.Enum = []string{"text", "table", "json"}
			}

			tool.InputSchema.Properties["flag_"+flag.Name] = property
		})

		*tools = append(*tools, tool)
	}

	// Process subcommands
	for _, subCmd := range cmd.Commands() {
		if subCmd.Hidden {
			continue
		}

		newPrefix := prefix
		if cmd.Use != "omnistrate-ctl" {
			// Add the command name to the prefix
			cmdName := cmd.Use
			if strings.Contains(cmdName, " ") {
				cmdName = strings.Split(cmdName, " ")[0]
			}
			newPrefix = prefix + cmdName + "-"
		}

		err := walkCommands(subCmd, newPrefix, tools)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateCapabilities creates the capabilities object based on available tools
func (s *MCPServer) generateCapabilities() map[string]interface{} {
	capabilities := map[string]interface{}{
		"tools": map[string]interface{}{
			"listChanged": false,
		},
	}

	// Add logging capability information
	log.Printf("MCP Server initialized with %d tools", len(s.tools))

	// Log some example tools for debugging
	if len(s.tools) > 0 {
		log.Printf("Example tools: %s, %s", s.tools[0].Name,
			func() string {
				if len(s.tools) > 1 {
					return s.tools[1].Name
				}
				return "(only one tool)"
			}())
	}

	return capabilities
}
