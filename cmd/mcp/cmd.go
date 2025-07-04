package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/omnistrate-oss/omnistrate-ctl/cmd/mcp/server"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/config"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
	"github.com/spf13/cobra"
)

// Cmd represents the MCP command
var Cmd = &cobra.Command{
	Use:   "mcp",
	Short: "Model Context Protocol server for omnistrate-ctl",
	Long: `This command starts an MCP server that exposes all omnistrate-ctl commands as tools
that can be used by AI models following the Model Context Protocol specification.

The server automatically handles authentication by checking for existing login credentials.
If no valid credentials are found, you can provide email and password to login automatically.`,
	Run: runMCP,
}

var (
	listTools bool
	email     string
	password  string
	noAuth    bool
)

func init() {
	Cmd.Flags().BoolVar(&listTools, "list-tools", false, "List all available tools and exit")
	Cmd.Flags().StringVar(&email, "email", "", "Email for authentication (optional)")
	Cmd.Flags().StringVar(&password, "password", "", "Password for authentication (optional)")
	Cmd.Flags().BoolVar(&noAuth, "no-auth", false, "Skip authentication (for testing only)")
	
	// Hide the no-auth flag as it's for testing only
	Cmd.Flags().MarkHidden("no-auth")
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

	// Handle authentication (skip if --no-auth flag is set)
	if !noAuth {
		if err := ensureAuthentication(); err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}
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

// ensureAuthentication ensures the user is authenticated before starting the MCP server
func ensureAuthentication() error {
	// If email and password are provided, attempt to login
	if email != "" && password != "" {
		return performLogin()
	}
	
	// For MCP server, we should not prompt for interactive login
	// Just check if there's an existing valid token
	_, err := config.GetToken()
	if err != nil {
		return fmt.Errorf("no valid authentication found. Please provide --email and --password flags or log in first using 'omnistrate-ctl login'")
	}
	
	return nil
}

// performLogin performs authentication using the provided email and password
func performLogin() error {
	// Validate input
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	
	// Use the same approach as the login command
	token, err := dataaccess.LoginWithPassword(nil, email, password)
	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}
	
	// Save the token
	authConfig := config.AuthConfig{
		Token: token,
	}
	if err = config.CreateOrUpdateAuthConfig(authConfig); err != nil {
		return fmt.Errorf("failed to save authentication config: %v", err)
	}
	
	return nil
}