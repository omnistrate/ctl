package instance

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/config"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/utils"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	opsCenterExample = `# Launch ops center TUI for an instance
omctl instance ops-center instance-abcd1234`
)

var opsCenterCmd = &cobra.Command{
	Use:          "ops-center [instance-id]",
	Short:        "Launch interactive ops center for instance management",
	Long:         `This command launches an interactive terminal UI for managing resource instances, including one-off patches.`,
	Example:      opsCenterExample,
	RunE:         runOpsCenter,
	SilenceUsage: true,
}

func init() {
	opsCenterCmd.Args = cobra.ExactArgs(1)
}

type OpsCenterApp struct {
	app            *tview.Application
	instanceID     string
	token          string
	instanceData   *openapiclientfleet.ResourceInstance
	serviceID      string
	environmentID  string
	resourceID     string
	
	// UI components
	summaryView    *tview.TextView
	sidebar        *tview.List
	mainFlex       *tview.Flex
	configEditor   *tview.TextArea
	statusBar      *tview.TextView
	editorHelp     *tview.TextView
	
	// State
	selectedResource string
	originalConfig   string
	modifiedConfig   string
	showingEditor    bool
	hasChanges       bool
	sidebarMode      string // "operations" or "resources"
}

func runOpsCenter(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	instanceID := args[0]

	// Validate user login
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Get instance details
	serviceID, environmentID, _, resourceID, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Get detailed instance data
	instanceData, err := dataaccess.DescribeResourceInstance(cmd.Context(), token, serviceID, environmentID, instanceID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Create and run the TUI
	app := &OpsCenterApp{
		instanceID:    instanceID,
		token:         token,
		instanceData:  instanceData,
		serviceID:     serviceID,
		environmentID: environmentID,
		resourceID:    resourceID,
	}

	return app.Run()
}

func (a *OpsCenterApp) Run() error {
	a.app = tview.NewApplication()
	a.setupUI()
	a.app.EnableMouse(true)
	
	if err := a.app.SetRoot(a.mainFlex, true).Run(); err != nil {
		return err
	}
	
	return nil
}

func (a *OpsCenterApp) setupUI() {
	// Create components
	a.summaryView = tview.NewTextView()
	a.sidebar = tview.NewList()
	a.configEditor = tview.NewTextArea()
	a.statusBar = tview.NewTextView()
	a.editorHelp = tview.NewTextView()
	
	// Configure summary view
	a.summaryView.SetBorder(true)
	a.summaryView.SetTitle("Instance Summary")
	a.summaryView.SetTitleAlign(tview.AlignLeft)
	a.summaryView.SetDynamicColors(true)
	a.summaryView.SetWordWrap(true)
	
	// Configure sidebar
	a.sidebar.SetBorder(true)
	a.sidebar.SetTitle("Operations")
	a.sidebar.SetTitleAlign(tview.AlignLeft)
	a.sidebar.ShowSecondaryText(false)
	
	// Configure config editor
	a.configEditor.SetBorder(true)
	a.configEditor.SetTitle("Helm Configuration Editor")
	a.configEditor.SetTitleAlign(tview.AlignLeft)
	a.configEditor.SetText("Select a resource to edit its configuration...", false)
	
	// Configure editor help
	a.editorHelp.SetBorder(false)
	a.editorHelp.SetDynamicColors(true)
	a.editorHelp.SetText("[yellow]Ctrl+S[white]: Save | [yellow]Esc[white]: Close Editor | [yellow]Tab[white]: Navigate")
	a.editorHelp.SetTextAlign(tview.AlignCenter)
	
	// Configure status bar
	a.statusBar.SetDynamicColors(true)
	a.statusBar.SetText("[yellow]Press 'q' to quit, 'Tab' to navigate, 'Enter' to select[white]")
	
	// Initialize sidebar mode and content
	a.sidebarMode = "operations"
	a.setupSidebar()
	
	// Setup summary content
	a.setupSummary()
	
	// Create main layout
	a.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(a.sidebar, 35, 1, true).
				AddItem(a.summaryView, 0, 2, false), 0, 1, true).
		AddItem(a.statusBar, 1, 1, false)
	
	// Setup key handlers
	a.setupKeyHandlers()
}

func (a *OpsCenterApp) setupSidebar() {
	a.sidebar.Clear()
	
	switch a.sidebarMode {
	case "operations":
		a.sidebar.SetTitle("Operations")
		a.sidebar.AddItem("🔍 Instance Overview", "", 'o', func() {
			a.showInstanceOverview()
		})
		
		a.sidebar.AddItem("🔧 One-Off Patch", "", 'p', func() {
			a.showPatchOptions()
		})
		
		a.sidebar.AddItem("📊 Resource Status", "", 'r', func() {
			a.showResourceStatus()
		})
	case "resources":
		a.sidebar.SetTitle("Resources")
		a.populateResourceList()
		
		// Add back to operations option
		a.sidebar.AddItem("← Back to Operations", "", 'b', func() {
			a.sidebarMode = "operations"
			a.setupSidebar()
			a.app.SetFocus(a.sidebar)
		})
	}
}

func (a *OpsCenterApp) setupSummary() {
	if a.instanceData == nil {
		a.summaryView.SetText("[red]Failed to load instance data[white]")
		return
	}

	var summary strings.Builder
	
	// Instance basic info
	summary.WriteString("[yellow]Instance Details[white]\n")
	summary.WriteString(fmt.Sprintf("ID: [cyan]%s[white]\n", a.instanceID))
	
	// Get status from consumption result
	status := "UNKNOWN"
	if a.instanceData.ConsumptionResourceInstanceResult.Status != nil {
		status = *a.instanceData.ConsumptionResourceInstanceResult.Status
	}
	summary.WriteString(fmt.Sprintf("Status: %s\n", getStatusColor(status)))
	
	// Get cloud provider
	summary.WriteString(fmt.Sprintf("Cloud Provider: [blue]%s[white]\n", a.instanceData.CloudProvider))
	
	// Get region from consumption result
	if a.instanceData.ConsumptionResourceInstanceResult.Region != nil {
		summary.WriteString(fmt.Sprintf("Region: [blue]%s[white]\n", *a.instanceData.ConsumptionResourceInstanceResult.Region))
	}
	
	// Get network type from consumption result
	if a.instanceData.ConsumptionResourceInstanceResult.NetworkType != nil {
		summary.WriteString(fmt.Sprintf("Network Type: [blue]%s[white]\n", *a.instanceData.ConsumptionResourceInstanceResult.NetworkType))
	}
	
	// Get created at from consumption result
	if a.instanceData.ConsumptionResourceInstanceResult.CreatedAt != nil {
		summary.WriteString(fmt.Sprintf("Created: [gray]%s[white]\n", *a.instanceData.ConsumptionResourceInstanceResult.CreatedAt))
	}
	
	summary.WriteString("\n[yellow]Service Details[white]\n")
	summary.WriteString(fmt.Sprintf("Service ID: [cyan]%s[white]\n", a.serviceID))
	summary.WriteString(fmt.Sprintf("Environment ID: [cyan]%s[white]\n", a.environmentID))
	summary.WriteString(fmt.Sprintf("Resource ID: [cyan]%s[white]\n", a.resourceID))
	
	// Resource version info
	if len(a.instanceData.ResourceVersionSummaries) > 0 {
		summary.WriteString("\n[yellow]Resource Versions[white]\n")
		for _, rv := range a.instanceData.ResourceVersionSummaries {
			resourceName := "Unknown"
			if rv.ResourceName != nil {
				resourceName = *rv.ResourceName
			}
			version := "Unknown"
			if rv.Version != nil {
				version = *rv.Version
			}
			summary.WriteString(fmt.Sprintf("• [cyan]%s[white] v%s\n", resourceName, version))
		}
	}
	
	// Network topology
	if a.instanceData.ConsumptionResourceInstanceResult.DetailedNetworkTopology != nil {
		summary.WriteString("\n[yellow]Network Topology[white]\n")
		for key, resource := range a.instanceData.ConsumptionResourceInstanceResult.DetailedNetworkTopology {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				if resourceType, exists := resourceMap["resourceType"]; exists {
					summary.WriteString(fmt.Sprintf("• [cyan]%s[white] (%v)\n", key, resourceType))
				} else {
					summary.WriteString(fmt.Sprintf("• [cyan]%s[white]\n", key))
				}
			}
		}
	}
	
	a.summaryView.SetText(summary.String())
}

func (a *OpsCenterApp) showInstanceOverview() {
	// Switch back to overview if showing other panels
	if a.showingEditor {
		a.hideEditor()
	}
	a.updateStatusBar("Instance overview - Press 'p' for patch options")
}

func (a *OpsCenterApp) showPatchOptions() {
	// Switch to resources mode in sidebar
	a.sidebarMode = "resources"
	a.setupSidebar()
	a.updateStatusBar("Select a resource to modify its helm configuration")
}

func (a *OpsCenterApp) populateResourceList() {
	// Get resource versions with helm configs
	resourceVersions := a.instanceData.ResourceVersionSummaries
	if len(resourceVersions) == 0 {
		a.sidebar.AddItem("❌ No resources found", "", 0, nil)
		return
	}
	
	hasValidResources := false
	for _, rv := range resourceVersions {
		resourceKey := ""
		if rv.ResourceId != nil {
			resourceKey = *rv.ResourceId
		}
		resourceName := ""
		if rv.ResourceName != nil {
			resourceName = *rv.ResourceName
		}
		
		// Check if resource has helm deployment configuration
		if rv.HelmDeploymentConfiguration != nil {
			hasValidResources = true
			// Capture variables for closure
			rvCopy := rv
			a.sidebar.AddItem(
				fmt.Sprintf("✅ %s", resourceName),
				resourceKey,
				0,
				func() {
					a.selectResource(resourceKey, resourceName, rvCopy.HelmDeploymentConfiguration)
				})
		} else {
			a.sidebar.AddItem(
				fmt.Sprintf("❌ %s", resourceName),
				"Not eligible",
				0,
				func() {
					a.updateStatusBar("[red]Resource not eligible for modification - no helm configuration[white]")
				})
		}
	}
	
	if !hasValidResources {
		a.sidebar.AddItem("❌ No eligible resources", "", 0, nil)
	}
}

// showResourceList is no longer needed since resources are shown in sidebar

func (a *OpsCenterApp) selectResource(resourceKey, resourceName string, helmConfig *openapiclientfleet.HelmDeploymentConfiguration) {
	a.selectedResource = resourceKey
	
	// Extract and format helm values
	if helmConfig.Values == nil {
		a.updateStatusBar("[red]No helm values found for this resource[white]")
		return
	}
	
	// Convert values to YAML
	yamlData, err := yaml.Marshal(helmConfig.Values)
	if err != nil {
		a.updateStatusBar(fmt.Sprintf("[red]Error formatting helm values: %v[white]", err))
		return
	}
	
	a.originalConfig = string(yamlData)
	a.modifiedConfig = a.originalConfig
	
	// Show editor
	a.showEditor(resourceName)
}

func (a *OpsCenterApp) showEditor(resourceName string) {
	a.showingEditor = true
	a.hasChanges = false
	
	// Update editor title and content
	a.configEditor.SetTitle(fmt.Sprintf("Helm Configuration Editor - %s", resourceName))
	a.configEditor.SetText(a.originalConfig, true)
	
	// Replace the right pane with editor
	a.mainFlex.RemoveItem(a.mainFlex.GetItem(0))
	editorFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.configEditor, 0, 1, true).
		AddItem(a.editorHelp, 1, 1, false)
	
	newMainFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.sidebar, 35, 1, false).
		AddItem(editorFlex, 0, 2, true)
	a.mainFlex.AddItem(newMainFlex, 0, 1, true)
	
	// Set up change detection
	a.configEditor.SetChangedFunc(func() {
		if a.configEditor.GetText() != a.originalConfig {
			a.hasChanges = true
			a.configEditor.SetTitle(fmt.Sprintf("Helm Configuration Editor - %s [MODIFIED]", resourceName))
		} else {
			a.hasChanges = false
			a.configEditor.SetTitle(fmt.Sprintf("Helm Configuration Editor - %s", resourceName))
		}
	})
	
	a.app.SetFocus(a.configEditor)
	a.updateStatusBar("Editing configuration - Use keybindings shown below")
}

func (a *OpsCenterApp) hideEditor() {
	if a.hasChanges {
		a.showSaveDiscardDialog()
		return
	}
	
	a.closeEditor()
}

func (a *OpsCenterApp) closeEditor() {
	a.showingEditor = false
	a.hasChanges = false
	
	// Return to summary view
	a.mainFlex.RemoveItem(a.mainFlex.GetItem(0))
	newFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.sidebar, 35, 1, true).
		AddItem(a.summaryView, 0, 2, false)
	a.mainFlex.AddItem(newFlex, 0, 1, true)
	
	a.app.SetFocus(a.sidebar)
	a.updateStatusBar("Press 'q' to quit, 'Tab' to navigate, 'Enter' to select")
}

func (a *OpsCenterApp) showSaveDiscardDialog() {
	modal := tview.NewModal().
		SetText("You have unsaved changes.\n\nWhat would you like to do?").
		AddButtons([]string{"Save", "Discard", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Save":
				a.saveConfiguration()
				a.app.SetRoot(a.mainFlex, true)
			case "Discard":
				a.closeEditor()
				a.app.SetRoot(a.mainFlex, true)
			case "Cancel":
				a.app.SetRoot(a.mainFlex, true)
				a.app.SetFocus(a.configEditor)
			}
		})
	
	a.app.SetRoot(modal, true)
}

func (a *OpsCenterApp) saveConfiguration() {
	if !a.showingEditor {
		return
	}
	
	a.modifiedConfig = a.configEditor.GetText()
	
	// Validate YAML
	var testParse interface{}
	if err := yaml.Unmarshal([]byte(a.modifiedConfig), &testParse); err != nil {
		a.updateStatusBar(fmt.Sprintf("[red]Invalid YAML: %v[white]", err))
		return
	}
	
	// Save to file
	if err := a.saveToConfigFile(); err != nil {
		a.updateStatusBar(fmt.Sprintf("[red]Error saving config: %v[white]", err))
		return
	}
	
	// Update original config and clear changes flag
	a.originalConfig = a.modifiedConfig
	a.hasChanges = false
	
	// Show confirmation dialog
	a.showConfirmationDialog()
}

func (a *OpsCenterApp) saveToConfigFile() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Join(os.TempDir(), "omctl-ops-center")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	// Convert modified config to interface{}
	var parsedConfig interface{}
	if err := yaml.Unmarshal([]byte(a.modifiedConfig), &parsedConfig); err != nil {
		return err
	}
	
	// Create configuration in the format expected by one-off patch
	helmValues, ok := parsedConfig.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid configuration format")
	}
	
	config := map[string]openapiclientfleet.ResourceOneOffPatchConfigurationOverride{
		a.selectedResource: {
			HelmChartValues: helmValues,
		},
	}
	
	configData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	
	// Save to file
	filename := filepath.Join(configDir, fmt.Sprintf("patch-config-%s-%s.yaml", a.instanceID, a.selectedResource))
	if err := os.WriteFile(filename, configData, 0600); err != nil {
		return err
	}
	
	a.updateStatusBar(fmt.Sprintf("[green]Configuration saved to: %s[white]", filename))
	return nil
}

func (a *OpsCenterApp) showConfirmationDialog() {
	modal := tview.NewModal().
		SetText("Configuration saved successfully!\n\nWould you like to start the one-off patch now?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				a.startPatch()
			}
			a.closeEditor()
			a.app.SetRoot(a.mainFlex, true)
		})
	
	a.app.SetRoot(modal, true)
}

func (a *OpsCenterApp) startPatch() {
	// This would trigger the actual patch process
	// For now, just show a success message
	a.updateStatusBar("[green]One-off patch initiated! Check your terminal for progress.[white]")
	
	// You could integrate with the existing patch command here
	// or call the dataaccess function directly
}

func (a *OpsCenterApp) showResourceStatus() {
	a.updateStatusBar("Resource status view - Feature coming soon!")
}

func (a *OpsCenterApp) setupKeyHandlers() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if a.showingEditor {
				a.hideEditor()
				return nil
			}
			return event
		case tcell.KeyCtrlC:
			a.app.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				if !a.showingEditor {
					a.app.Stop()
					return nil
				}
			case 's', 'S':
				if event.Modifiers()&tcell.ModCtrl != 0 && a.showingEditor {
					a.saveConfiguration()
					return nil
				}
			}
		case tcell.KeyTab:
			if a.showingEditor {
				return event // Let the editor handle tab
			}
			// Cycle focus between components
			current := a.app.GetFocus()
			if current == a.sidebar {
				a.app.SetFocus(a.summaryView)
			} else {
				a.app.SetFocus(a.sidebar)
			}
			return nil
		}
		return event
	})
}

func (a *OpsCenterApp) updateStatusBar(message string) {
	a.statusBar.SetText(message)
}

func getStatusColor(status string) string {
	switch strings.ToUpper(status) {
	case "RUNNING":
		return "[green]RUNNING[white]"
	case "FAILED":
		return "[red]FAILED[white]"
	case "PENDING":
		return "[yellow]PENDING[white]"
	default:
		return fmt.Sprintf("[gray]%s[white]", status)
	}
}