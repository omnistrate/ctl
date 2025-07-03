package mcp

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/omnistrate-oss/omnistrate-ctl/cmd/mcp/server"
	"github.com/spf13/cobra"
)

// Cmd represents the MCP command
var Cmd = &cobra.Command{
	Use:   "mcp",
	Short: "Model Context Protocol server for omnistrate-ctl",
	Long: `This command starts an MCP server that exposes all omnistrate-ctl commands as tools
that can be used by AI models following the Model Context Protocol specification.`,
	Run: runMCP,
}

var listTools bool

func init() {
	Cmd.Flags().BoolVar(&listTools, "list-tools", false, "List all available tools and exit")
}

func runMCP(cmd *cobra.Command, args []string) {
	if listTools {
		// List all available tools
		tools, err := server.GetAllTools(cmd.Root())
		if err != nil {
			log.Fatalf("Failed to get tools: %v", err)
		}
		
		data, err := json.MarshalIndent(tools, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal tools: %v", err)
		}
		
		fmt.Println(string(data))
		return
	}

	// Start MCP server
	mcpServer, err := server.NewMCPServer(cmd.Root())
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	if err := mcpServer.Start(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}