package instance

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"

	"github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
)

var debugCmd = &cobra.Command{
	Use:     "debug [instance-id]",
	Short:   "Debug instance resources",
	Long:    `Debug instance resources with an interactive TUI showing helm charts, terraform files, and logs.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDebug,
	Example: `  omnistrate-ctl instance debug <instance-id>`,
}

type DebugData struct {
	InstanceID string         `json:"instanceId"`
	Resources  []ResourceInfo `json:"resources"`
}

type ResourceInfo struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Type          string         `json:"type"` // "helm" or "terraform"
	DebugData     interface{}    `json:"debugData"`
	HelmData      *HelmData      `json:"helmData,omitempty"`
	TerraformData *TerraformData `json:"terraformData,omitempty"`
}

type HelmData struct {
	ChartRepoName string                 `json:"chartRepoName"`
	ChartRepoURL  string                 `json:"chartRepoURL"`
	ChartVersion  string                 `json:"chartVersion"`
	ChartValues   map[string]interface{} `json:"chartValues"`
	InstallLog    string                 `json:"installLog"`
	Namespace     string                 `json:"namespace"`
	ReleaseName   string                 `json:"releaseName"`
}

type TerraformData struct {
	Files map[string]string `json:"files"`
	Logs  map[string]string `json:"logs"`
}

func runDebug(_ *cobra.Command, args []string) error {
	instanceID := args[0]

	token, err := common.GetTokenWithLogin()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	ctx := context.Background()

	// Get instance details
	serviceID, environmentID, _, _, err := getInstance(ctx, token, instanceID)
	if err != nil {
		return fmt.Errorf("failed to get instance: %w", err)
	}

	// Get debug information
	debugResult, err := dataaccess.DebugResourceInstance(ctx, token, serviceID, environmentID, instanceID)
	if err != nil {
		return fmt.Errorf("failed to get debug information: %w", err)
	}

	// Process debug result
	data := DebugData{
		InstanceID: instanceID,
		Resources:  []ResourceInfo{},
	}

	if debugResult.ResourcesDebug != nil {
		for resourceKey, resourceDebugInfo := range debugResult.ResourcesDebug {
			resourceInfo := ResourceInfo{
				ID:        resourceKey,
				Name:      resourceKey,
				Type:      "unknown",
				DebugData: resourceDebugInfo,
			}

			if debugData, ok := resourceDebugInfo.(map[string]interface{}); ok {
				if actualDebugData, ok := debugData["debugData"].(map[string]interface{}); ok {
					// Check if it's a helm resource
					if _, hasChart := actualDebugData["chartRepoName"]; hasChart {
						resourceInfo.Type = "helm"
						resourceInfo.HelmData = parseHelmData(actualDebugData)
					} else {
						// Check if it's a terraform resource by looking for terraform files or logs
						hasTerraformFiles := false
						hasTerraformLogs := false

						for key := range actualDebugData {
							if strings.HasPrefix(key, "rendered/") && strings.HasSuffix(key, ".tf") {
								hasTerraformFiles = true
							} else if strings.HasPrefix(key, "log/") && strings.Contains(key, "terraform") {
								hasTerraformLogs = true
							}
						}

						if hasTerraformFiles || hasTerraformLogs {
							resourceInfo.Type = "terraform"
							resourceInfo.TerraformData = parseTerraformData(actualDebugData)
						}
					}
				}
			}

			data.Resources = append(data.Resources, resourceInfo)
		}
	}

	// Launch TUI
	return launchDebugTUI(data)
}

func parseHelmData(debugData map[string]interface{}) *HelmData {
	helmData := &HelmData{
		ChartValues: make(map[string]interface{}),
	}

	if chartRepoName, ok := debugData["chartRepoName"].(string); ok {
		helmData.ChartRepoName = chartRepoName
	}
	if chartRepoURL, ok := debugData["chartRepoURL"].(string); ok {
		helmData.ChartRepoURL = chartRepoURL
	}
	if chartVersion, ok := debugData["chartVersion"].(string); ok {
		helmData.ChartVersion = chartVersion
	}
	if namespace, ok := debugData["namespace"].(string); ok {
		helmData.Namespace = namespace
	}
	if releaseName, ok := debugData["releaseName"].(string); ok {
		helmData.ReleaseName = releaseName
	}

	// Parse chart values
	if chartValuesStr, ok := debugData["chartValues"].(string); ok {
		var chartValues map[string]interface{}
		if err := json.Unmarshal([]byte(chartValuesStr), &chartValues); err == nil {
			helmData.ChartValues = chartValues
		}
	}

	// Parse install log
	if installLog, ok := debugData["log/install.log"].(string); ok {
		helmData.InstallLog = installLog
	}

	return helmData
}

func parseTerraformData(debugData map[string]interface{}) *TerraformData {
	terraformData := &TerraformData{
		Files: make(map[string]string),
		Logs:  make(map[string]string),
	}

	// Parse all files and logs
	for key, value := range debugData {
		if strValue, ok := value.(string); ok {
			if strings.HasPrefix(key, "rendered/") && strings.HasSuffix(key, ".tf") {
				terraformData.Files[key] = strValue
			} else if strings.HasPrefix(key, "log/") {
				terraformData.Logs[key] = strValue
			}
		}
	}

	return terraformData
}

func launchDebugTUI(data DebugData) error {
	app := tview.NewApplication()

	// Global state to track current selection and terraform data for file browser
	var currentTerraformData *TerraformData
	var currentSelectionIsTerraformFiles bool

	// Create main layout
	flex := tview.NewFlex()

	// Left panel - Resources (accordion style)
	leftPanel := tview.NewTreeView()
	leftPanel.SetBorder(true).SetTitle("Resources")

	// Create root node
	root := tview.NewTreeNode(fmt.Sprintf("Instance: %s", data.InstanceID))
	root.SetColor(tcell.ColorYellow)
	leftPanel.SetRoot(root)

	// Add resources (only helm and terraform, skip unknown types)
	for _, resource := range data.Resources {
		// Skip unknown resource types
		if resource.Type != "helm" && resource.Type != "terraform" {
			continue
		}

		resourceNode := tview.NewTreeNode(fmt.Sprintf("%s (%s)", resource.Name, resource.Type))
		resourceNode.SetReference(resource)
		resourceNode.SetColor(tcell.ColorBlue)

		// Add options based on resource type
		if resource.Type == "helm" && resource.HelmData != nil {
			// Add Chart Values option
			chartValuesNode := tview.NewTreeNode("Chart Values")
			chartValuesNode.SetReference(map[string]interface{}{
				"type":     "helm-chart-values",
				"resource": resource,
			})
			chartValuesNode.SetColor(tcell.ColorGreen)
			resourceNode.AddChild(chartValuesNode)

			// Add Install Log option
			if resource.HelmData.InstallLog != "" {
				installLogNode := tview.NewTreeNode("Install Log")
				installLogNode.SetReference(map[string]interface{}{
					"type":     "helm-install-log",
					"resource": resource,
				})
				installLogNode.SetColor(tcell.ColorGreen)
				resourceNode.AddChild(installLogNode)
			}
		} else if resource.Type == "terraform" && resource.TerraformData != nil {
			// Add Terraform Files option
			if len(resource.TerraformData.Files) > 0 {
				filesNode := tview.NewTreeNode("Terraform Files")
				filesNode.SetReference(map[string]interface{}{
					"type":     "terraform-files",
					"resource": resource,
				})
				filesNode.SetColor(tcell.ColorGreen)
				resourceNode.AddChild(filesNode)
			}

			// Add Install Log option
			if len(resource.TerraformData.Logs) > 0 {
				installLogNode := tview.NewTreeNode("Install Logs")
				installLogNode.SetReference(map[string]interface{}{
					"type":     "terraform-install-logs",
					"resource": resource,
				})
				installLogNode.SetColor(tcell.ColorGreen)
				resourceNode.AddChild(installLogNode)
			}
		}

		root.AddChild(resourceNode)
	}

	root.SetExpanded(true)

	// Right panel - Content
	rightPanel := tview.NewTextView()
	rightPanel.SetBorder(true).SetTitle("Content")
	rightPanel.SetDynamicColors(true)
	rightPanel.SetWrap(true)
	rightPanel.SetScrollable(true)
	rightPanel.SetText("Select a resource option to view details")

	// Add focus handlers to show which panel is active
	leftPanel.SetFocusFunc(func() {
		leftPanel.SetBorderColor(tcell.ColorGreen)
		rightPanel.SetBorderColor(tcell.ColorDefault)
	})
	rightPanel.SetFocusFunc(func() {
		rightPanel.SetBorderColor(tcell.ColorGreen)
		leftPanel.SetBorderColor(tcell.ColorDefault)
	})

	// Handle tree selection
	leftPanel.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			rightPanel.SetTitle("Content")
			rightPanel.SetText("Select a resource option to view details")
			// Clear terraform file selection state when no valid selection
			currentSelectionIsTerraformFiles = false
			return
		}

		switch ref := reference.(type) {
		case ResourceInfo:
			// Show resource information
			content := formatResourceInfo(ref)
			rightPanel.SetTitle(fmt.Sprintf("Resource: %s", ref.Name))
			rightPanel.SetText(content)
			// Clear terraform file selection state when selecting resource node
			currentSelectionIsTerraformFiles = false
		case map[string]interface{}:
			handleOptionSelection(ref, rightPanel)
			// Update current terraform data and selection state for file browser
			if optionType, ok := ref["type"].(string); ok && optionType == "terraform-files" {
				if resource, ok := ref["resource"].(ResourceInfo); ok {
					currentTerraformData = resource.TerraformData
					currentSelectionIsTerraformFiles = true
				}
			} else {
				currentSelectionIsTerraformFiles = false
			}
		}
	})

	// Also handle direct selection (Enter key)
	leftPanel.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference != nil {
			// If it's an option, show its content
			switch ref := reference.(type) {
			case ResourceInfo:
				content := formatResourceInfo(ref)
				rightPanel.SetTitle(fmt.Sprintf("Resource: %s", ref.Name))
				rightPanel.SetText(content)
				// Clear terraform file selection state when selecting resource node
				currentSelectionIsTerraformFiles = false
			case map[string]interface{}:
				handleOptionSelection(ref, rightPanel)
				// Update current terraform data and selection state for file browser
				if optionType, ok := ref["type"].(string); ok && optionType == "terraform-files" {
					if resource, ok := ref["resource"].(ResourceInfo); ok {
						currentTerraformData = resource.TerraformData
						currentSelectionIsTerraformFiles = true
					}
				} else {
					currentSelectionIsTerraformFiles = false
				}
				return // Don't toggle expansion for options
			}
		}
		// Toggle expansion for resource nodes
		node.SetExpanded(!node.IsExpanded())
	})

	// Set up layout
	flex.AddItem(leftPanel, 0, 1, true)
	flex.AddItem(rightPanel, 0, 2, false)

	// Create main layout with help text
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	mainFlex.AddItem(flex, 0, 1, true)
	mainFlex.AddItem(createHelpText(), 1, 0, false)

	// Create main input handler
	var mainInputHandler func(event *tcell.EventKey) *tcell.EventKey
	mainInputHandler = func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyEnter:
			// Switch from left panel to right panel to view content
			if app.GetFocus() == leftPanel {
				app.SetFocus(rightPanel)
				return nil
			}
			// If on right panel, let default behavior handle scrolling
		case tcell.KeyEscape:
			// Go back to left panel from right panel
			if app.GetFocus() == rightPanel {
				app.SetFocus(leftPanel)
				return nil
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				app.Stop()
				return nil
			case 'f', 'F':
				if currentSelectionIsTerraformFiles && currentTerraformData != nil && len(currentTerraformData.Files) > 0 {
					showFileBrowser(app, currentTerraformData, mainFlex, mainInputHandler)
				}
				return nil
			}
		}
		return event
	}

	// Set the main input handler
	app.SetInputCapture(mainInputHandler)

	// Set initial focus and selection
	app.SetFocus(leftPanel)

	// Set initial selection to first resource if available
	if len(data.Resources) > 0 {
		// Find the first resource node
		if len(root.GetChildren()) > 0 {
			firstResource := root.GetChildren()[0]
			leftPanel.SetCurrentNode(firstResource)
			// Expand the first resource to show its options
			firstResource.SetExpanded(true)
		}
	}

	// Start the application (disable mouse to allow terminal text selection)
	if err := app.SetRoot(mainFlex, true).EnableMouse(false).Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}

func createHelpText() *tview.TextView {
	helpText := tview.NewTextView()
	helpText.SetText("Navigate: ↑/↓ to move | Enter: view content/expand | Esc: go back | f: file browser | q: quit")
	helpText.SetTextAlign(tview.AlignCenter)
	helpText.SetDynamicColors(true)
	return helpText
}

func handleOptionSelection(ref map[string]interface{}, rightPanel *tview.TextView) {
	optionType, _ := ref["type"].(string)
	resource, _ := ref["resource"].(ResourceInfo)

	switch optionType {
	case "helm-chart-values":
		if resource.HelmData != nil {
			content := formatHelmChartValues(resource.HelmData)
			rightPanel.SetTitle("Chart Values")
			rightPanel.SetText(content)
		}
	case "helm-install-log":
		if resource.HelmData != nil {
			content := formatHelmInstallLog(resource.HelmData.InstallLog)
			rightPanel.SetTitle("Install Log")
			rightPanel.SetText(content)
		}
	case "terraform-files":
		if resource.TerraformData != nil {
			content := formatTerraformFileList(resource.TerraformData.Files)
			rightPanel.SetTitle("Terraform Files")
			rightPanel.SetText(content)
		}
	case "terraform-install-logs":
		if resource.TerraformData != nil {
			content := formatTerraformLogs(resource.TerraformData.Logs)
			rightPanel.SetTitle("Install Logs")
			rightPanel.SetText(content)
		}
	}
}

func formatResourceInfo(resource ResourceInfo) string {
	debugInfo := ""
	if resource.Type == "terraform" && resource.TerraformData != nil {
		debugInfo = fmt.Sprintf("\n\nTerraform Files: %d\nTerraform Logs: %d", len(resource.TerraformData.Files), len(resource.TerraformData.Logs))
	} else if resource.Type == "helm" && resource.HelmData != nil {
		debugInfo = fmt.Sprintf("\n\nChart: %s\nInstall Log: %t", resource.HelmData.ChartRepoName, resource.HelmData.InstallLog != "")
	}

	return fmt.Sprintf(`[yellow]Resource Information[white]

Name: %s
Type: %s
ID: %s%s

Select an option from the tree to view specific details.`, resource.Name, resource.Type, resource.ID, debugInfo)
}

func formatHelmChartValues(helmData *HelmData) string {
	content := fmt.Sprintf(`[yellow]Helm Chart Values[white]

Chart: %s
Version: %s
Repo: %s
Namespace: %s
Release: %s

[yellow]Values:[white]
`, helmData.ChartRepoName, helmData.ChartVersion, helmData.ChartRepoURL, helmData.Namespace, helmData.ReleaseName)

	if len(helmData.ChartValues) > 0 {
		jsonBytes, err := json.MarshalIndent(helmData.ChartValues, "", "  ")
		if err == nil {
			// Apply YAML syntax highlighting to the JSON output (similar structure)
			highlightedContent := addYAMLSyntaxHighlighting(string(jsonBytes))
			content += highlightedContent
		} else {
			content += fmt.Sprintf("Error formatting values: %v", err)
		}
	} else {
		content += "No chart values available"
	}

	return content
}

func formatHelmInstallLog(installLog string) string {
	if installLog == "" {
		return "[yellow]Install Log[white]\n\nNo install log available"
	}
	// Apply log syntax highlighting
	highlightedLog := addLogSyntaxHighlighting(installLog)
	return fmt.Sprintf(`[yellow]Install Log[white]

%s`, highlightedLog)
}

func formatTerraformFileList(files map[string]string) string {
	if len(files) == 0 {
		return "[yellow]Terraform Files[white]\n\nNo terraform files available"
	}

	content := "[yellow]Terraform Files[white]\n\nFiles available (press 'f' to open file browser):\n\n"

	// Build a hierarchical tree structure
	type TreeNode struct {
		Name     string
		IsDir    bool
		Children map[string]*TreeNode
		Files    []string
	}

	root := &TreeNode{
		Name:     "root",
		IsDir:    true,
		Children: make(map[string]*TreeNode),
		Files:    []string{},
	}

	// Get sorted file paths for deterministic ordering
	filePaths := make([]string, 0, len(files))
	for filePath := range files {
		filePaths = append(filePaths, filePath)
	}
	sort.Strings(filePaths)

	// Build the tree structure
	for _, filePath := range filePaths {
		parts := strings.Split(filePath, "/")
		currentNode := root

		// Navigate through directory parts
		for i, part := range parts {
			if i == len(parts)-1 {
				// This is a file
				currentNode.Files = append(currentNode.Files, part)
			} else {
				// This is a directory
				if currentNode.Children[part] == nil {
					currentNode.Children[part] = &TreeNode{
						Name:     part,
						IsDir:    true,
						Children: make(map[string]*TreeNode),
						Files:    []string{},
					}
				}
				currentNode = currentNode.Children[part]
			}
		}
	}

	// Function to render the tree
	var renderTree func(node *TreeNode, prefix string, isLast bool) string
	renderTree = func(node *TreeNode, prefix string, isLast bool) string {
		result := ""

		// Sort children directories and files
		var childNames []string
		for name := range node.Children {
			childNames = append(childNames, name)
		}
		sort.Strings(childNames)
		sort.Strings(node.Files)

		// Render child directories
		for i, childName := range childNames {
			child := node.Children[childName]
			isLastChild := (i == len(childNames)-1) && len(node.Files) == 0

			// Choose the right tree symbol
			var symbol, nextPrefix string
			if isLastChild {
				symbol = "└── "
				nextPrefix = prefix + "    "
			} else {
				symbol = "├── "
				nextPrefix = prefix + "│   "
			}

			result += fmt.Sprintf("%s[blue]%s%s/[-]\n", prefix, symbol, childName)
			result += renderTree(child, nextPrefix, true)
		}

		// Render files
		for i, fileName := range node.Files {
			isLastFile := i == len(node.Files)-1
			var symbol string
			if isLastFile {
				symbol = "└── "
			} else {
				symbol = "├── "
			}
			result += fmt.Sprintf("%s%s%s\n", prefix, symbol, fileName)
		}

		return result
	}

	// Render the tree starting from root
	content += renderTree(root, "", true)
	content += "\n[green]Press 'f' to open file browser and view individual files[-]"

	return content
}

func formatTerraformLogs(logs map[string]string) string {
	if len(logs) == 0 {
		return "[yellow]Terraform Logs[white]\n\nNo terraform logs available"
	}

	content := "[yellow]Terraform Logs[white]\n\n"

	// Sort log names for consistent output
	var sortedLogNames []string
	for logName := range logs {
		sortedLogNames = append(sortedLogNames, logName)
	}
	sort.Strings(sortedLogNames)

	for _, logName := range sortedLogNames {
		logContent := logs[logName]
		content += fmt.Sprintf("[blue]%s:[white]\n", logName)
		if logContent != "" {
			// Apply log syntax highlighting
			highlightedContent := addLogSyntaxHighlighting(logContent)
			content += fmt.Sprintf("%s\n\n", highlightedContent)
		} else {
			content += "(empty log)\n\n"
		}
	}

	return content
}

// addYAMLSyntaxHighlighting adds basic syntax highlighting for YAML content
func addYAMLSyntaxHighlighting(content string) string {
	lines := strings.Split(content, "\n")
	var highlighted []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			highlighted = append(highlighted, line)
			continue
		}

		// Comments
		if strings.HasPrefix(trimmed, "#") {
			highlighted = append(highlighted, fmt.Sprintf("[green]%s[-]", line))
			continue
		}

		// Keys (lines containing ':')
		if strings.Contains(line, ":") && !strings.HasPrefix(trimmed, "-") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				highlighted = append(highlighted, fmt.Sprintf("[cyan]%s[-]:[yellow]%s[-]", key, value))
				continue
			}
		}

		// List items
		if strings.HasPrefix(trimmed, "-") {
			highlighted = append(highlighted, fmt.Sprintf("[blue]%s[-]", line))
			continue
		}

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

// addTerraformSyntaxHighlighting adds basic syntax highlighting for Terraform files
func addTerraformSyntaxHighlighting(content string) string {
	lines := strings.Split(content, "\n")
	var highlighted []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			highlighted = append(highlighted, line)
			continue
		}

		// Comments
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			highlighted = append(highlighted, fmt.Sprintf("[green]%s[-]", line))
			continue
		}

		// Resource/data/variable/output blocks
		if strings.HasPrefix(trimmed, "resource ") || strings.HasPrefix(trimmed, "data ") ||
			strings.HasPrefix(trimmed, "variable ") || strings.HasPrefix(trimmed, "output ") ||
			strings.HasPrefix(trimmed, "provider ") || strings.HasPrefix(trimmed, "module ") {
			highlighted = append(highlighted, fmt.Sprintf("[fuchsia]%s[-]", line))
			continue
		}

		// String assignments (key = "value")
		if strings.Contains(line, "=") && strings.Contains(line, "\"") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := strings.TrimSpace(parts[1])
				// Highlight strings in quotes
				if strings.Contains(value, "\"") {
					value = strings.ReplaceAll(value, "\"", "[yellow]\"[-]")
				}
				highlighted = append(highlighted, fmt.Sprintf("[cyan]%s[-]= %s", key, value))
				continue
			}
		}

		// Simple assignments
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				highlighted = append(highlighted, fmt.Sprintf("[cyan]%s[-]=[blue]%s[-]", key, value))
				continue
			}
		}

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

// addLogSyntaxHighlighting adds basic syntax highlighting for log content
func addLogSyntaxHighlighting(content string) string {
	lines := strings.Split(content, "\n")
	var highlighted []string

	for _, line := range lines {
		lower := strings.ToLower(line)

		// Error messages
		if strings.Contains(lower, "error") || strings.Contains(lower, "failed") ||
			strings.Contains(lower, "panic") || strings.Contains(lower, "fatal") {
			highlighted = append(highlighted, fmt.Sprintf("[red]%s[-]", line))
			continue
		}

		// Warning messages
		if strings.Contains(lower, "warn") || strings.Contains(lower, "warning") {
			highlighted = append(highlighted, fmt.Sprintf("[yellow]%s[-]", line))
			continue
		}

		// Success messages
		if strings.Contains(lower, "success") || strings.Contains(lower, "complete") ||
			strings.Contains(lower, "applied") || strings.Contains(lower, "created") {
			highlighted = append(highlighted, fmt.Sprintf("[green]%s[-]", line))
			continue
		}

		// Info messages
		if strings.Contains(lower, "info") || strings.Contains(lower, "applying") ||
			strings.Contains(lower, "planning") || strings.Contains(lower, "refreshing") {
			highlighted = append(highlighted, fmt.Sprintf("[blue]%s[-]", line))
			continue
		}

		// Timestamps (basic detection)
		if strings.Contains(line, ":") && (strings.Contains(line, "T") ||
			strings.Contains(line, "[") && strings.Contains(line, "]")) {
			highlighted = append(highlighted, fmt.Sprintf("[gray]%s[-]", line))
			continue
		}

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

func showFileBrowser(app *tview.Application, terraformData *TerraformData, mainFlex *tview.Flex, originalInputHandler func(event *tcell.EventKey) *tcell.EventKey) {
	// Create file tree view (hierarchical)
	fileTree := tview.NewTreeView()
	fileTree.SetBorder(true).SetTitle("Terraform Files")

	// Create root node
	root := tview.NewTreeNode("Files")
	root.SetColor(tcell.ColorYellow)
	fileTree.SetRoot(root)

	// Build hierarchical file structure
	dirNodes := make(map[string]*tview.TreeNode)

	// Get sorted file paths for deterministic ordering
	filePaths := make([]string, 0, len(terraformData.Files))
	for filePath := range terraformData.Files {
		filePaths = append(filePaths, filePath)
	}
	sort.Strings(filePaths)

	// Helper function to get or create directory node
	var getOrCreateDirNode func(path string) *tview.TreeNode
	getOrCreateDirNode = func(path string) *tview.TreeNode {
		if path == "." || path == "" {
			return root
		}

		// Check if we already have this directory
		if node, exists := dirNodes[path]; exists {
			return node
		}

		// Create the directory node
		dirName := filepath.Base(path)
		dirNode := tview.NewTreeNode(dirName + "/")
		dirNode.SetColor(tcell.ColorBlue)
		dirNode.SetExpanded(false) // Allow user to expand/collapse
		dirNodes[path] = dirNode

		// Get parent directory and add this node to it
		parentPath := filepath.Dir(path)
		parentNode := getOrCreateDirNode(parentPath)
		parentNode.AddChild(dirNode)

		return dirNode
	}

	// Build the tree structure
	for _, filePath := range filePaths {
		dir := filepath.Dir(filePath)
		fileName := filepath.Base(filePath)

		// Get the parent directory node (creates all intermediate directories)
		parentNode := getOrCreateDirNode(dir)

		// Add file to parent directory
		fileNode := tview.NewTreeNode(fileName)
		fileNode.SetReference(filePath)
		fileNode.SetColor(tcell.ColorWhite)
		parentNode.AddChild(fileNode)
	}

	root.SetExpanded(true)

	// Create file content viewer
	fileViewer := tview.NewTextView()
	fileViewer.SetBorder(true).SetTitle("File Content")
	fileViewer.SetScrollable(true)
	fileViewer.SetWrap(false)
	fileViewer.SetDynamicColors(true) // Enable color rendering
	fileViewer.SetText("Select a file from the tree to view its content")

	// Handle tree selection
	fileTree.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference != nil {
			if filePath, ok := reference.(string); ok {
				if content, exists := terraformData.Files[filePath]; exists {
					fileViewer.SetTitle(fmt.Sprintf("File: %s", filePath))
					// Apply syntax highlighting based on file extension
					if strings.HasSuffix(filePath, ".tf") || strings.HasSuffix(filePath, ".tfvars") {
						highlightedContent := addTerraformSyntaxHighlighting(content)
						fileViewer.SetText(highlightedContent)
					} else {
						fileViewer.SetText(content)
					}
				}
			}
		}
	})

	// Handle tree node selection (Enter key)
	fileTree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference != nil {
			// If it's a file, show content and don't toggle expansion
			if filePath, ok := reference.(string); ok {
				if content, exists := terraformData.Files[filePath]; exists {
					fileViewer.SetTitle(fmt.Sprintf("File: %s", filePath))
					// Apply syntax highlighting based on file extension
					if strings.HasSuffix(filePath, ".tf") || strings.HasSuffix(filePath, ".tfvars") {
						highlightedContent := addTerraformSyntaxHighlighting(content)
						fileViewer.SetText(highlightedContent)
					} else {
						fileViewer.SetText(content)
					}
				}
				return // Don't toggle expansion for files
			}
		}
		// Toggle expansion for directory nodes (including root and subdirectories)
		node.SetExpanded(!node.IsExpanded())
	})

	// Add focus handlers to show which panel is active
	fileTree.SetFocusFunc(func() {
		fileTree.SetBorderColor(tcell.ColorGreen)
		fileViewer.SetBorderColor(tcell.ColorDefault)
	})
	fileViewer.SetFocusFunc(func() {
		fileViewer.SetBorderColor(tcell.ColorGreen)
		fileTree.SetBorderColor(tcell.ColorDefault)
	})

	// Create layout for file browser
	fileBrowserFlex := tview.NewFlex()
	fileBrowserFlex.AddItem(fileTree, 0, 1, true)
	fileBrowserFlex.AddItem(fileViewer, 0, 2, false)

	// Create modal frame
	modal := tview.NewFlex().SetDirection(tview.FlexRow)
	modal.AddItem(nil, 0, 1, false)
	modal.AddItem(tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(fileBrowserFlex, 0, 8, true).
		AddItem(nil, 0, 1, false), 0, 8, true)
	modal.AddItem(nil, 0, 1, false)

	// Help text for file browser
	helpText := tview.NewTextView()
	helpText.SetText("Navigate: ↑/↓ to select file | Enter: view content/expand | Esc: back/close | Content scrollable when focused")
	helpText.SetTextAlign(tview.AlignCenter)
	helpText.SetDynamicColors(true)

	// Final modal layout
	modalLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	modalLayout.AddItem(modal, 0, 1, true)
	modalLayout.AddItem(helpText, 1, 0, false)

	// Handle key events in file browser
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if app.GetFocus() == fileViewer {
				// If viewing content, go back to file tree
				app.SetFocus(fileTree)
				return nil
			} else {
				// If on file tree, close file browser and return to main view
				app.SetInputCapture(originalInputHandler) // Restore original input handler
				app.SetRoot(mainFlex, true)
				return nil
			}
		case tcell.KeyEnter:
			if app.GetFocus() == fileTree {
				// Let the tree view handle Enter first (for expand/collapse)
				// Only switch to content viewer if a file is selected
				currentNode := fileTree.GetCurrentNode()
				if currentNode != nil {
					reference := currentNode.GetReference()
					// If it's a file (has reference), switch to content viewer
					if _, isFile := reference.(string); isFile {
						app.SetFocus(fileViewer)
						return nil
					}
					// If it's a directory (no reference), let tree handle expansion
					// Don't consume the event, let it pass through to the tree
					return event
				}
			}
			// If already viewing content, let default behavior handle scrolling
		}
		return event
	})

	// Set initial focus and selection
	app.SetFocus(fileTree)

	// Set initial selection to first file if available
	if len(filePaths) > 0 {
		// Find the first file node in the tree
		var findFirstFileNode func(node *tview.TreeNode) *tview.TreeNode
		findFirstFileNode = func(node *tview.TreeNode) *tview.TreeNode {
			if node.GetReference() != nil {
				// This is a file node
				return node
			}
			// Check children for file nodes
			for _, child := range node.GetChildren() {
				if result := findFirstFileNode(child); result != nil {
					return result
				}
			}
			return nil
		}

		if firstFileNode := findFirstFileNode(root); firstFileNode != nil {
			fileTree.SetCurrentNode(firstFileNode)
		}
	}

	app.SetRoot(modalLayout, true).EnableMouse(false)
}

func init() {
	// Command will be added by the parent instance command
}
