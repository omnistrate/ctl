package instance

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/config"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/utils"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ResourceConfigChange tracks configuration changes for a resource
type ResourceConfigChange struct {
	ResourceID     string
	ResourceName   string
	ResourceType   string // "new", "existing", "deprecated"
	OriginalConfig string
	ModifiedConfig string
	IsConfigured   bool
}

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
	app           *tview.Application
	instanceID    string
	token         string
	instanceData  *openapiclientfleet.ResourceInstance
	serviceID     string
	environmentID string
	resourceID    string

	// UI components
	summaryView   *tview.TextView
	sidebar       *tview.List
	mainFlex      *tview.Flex
	contentArea   *tview.Flex // The right side content area
	configEditor  *tview.TextArea
	statusBar     *tview.TextView
	editorHelp    *tview.TextView
	versionsView  *tview.Table    // For target service plan versions
	resourcesView *tview.List     // For new/existing resources
	summaryView2  *tview.TextView // For patch summary
	summaryList   *tview.List     // For interactive summary selection
	summaryTable  *tview.Table    // For resource summary table
	diffViewer    *tview.TextView // For git-style diff
	commandViewer *tview.TextView // For command execution output

	// State
	selectedResource string
	originalConfig   string
	modifiedConfig   string
	showingEditor    bool
	hasChanges       bool
	sidebarMode      string // "operations", "resources", "versions", "resourceSelect", "summary"

	// One-off patch state
	targetVersion       string
	currentVersionSet   *openapiclient.TierVersionSet
	targetVersionSet    *openapiclient.TierVersionSet
	newResources        []openapiclient.ResourceSummary
	existingResources   []openapiclient.ResourceSummary
	deprecatedResources []openapiclient.ResourceSummary
	resourceMode        string                          // "new", "existing"
	configChanges       map[string]ResourceConfigChange // resource ID -> changes
	configFilename      string                          // Path to saved config file
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

	// Ensure terminal is properly reset on exit
	defer func() {
		if app.app != nil {
			app.app.Stop()
		}
	}()

	return app.Run()
}

func (a *OpsCenterApp) Run() error {
	a.app = tview.NewApplication()
	a.setupUI()
	a.app.EnableMouse(true)

	// Set up proper screen handling
	defer func() {
		if r := recover(); r != nil {
			a.app.Stop()
		}
	}()

	// Use SetRoot with resizeToFit=true for better screen handling
	if err := a.app.SetRoot(a.mainFlex, true).SetFocus(a.sidebar).Run(); err != nil {
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
	a.versionsView = tview.NewTable()
	a.resourcesView = tview.NewList()
	a.summaryView2 = tview.NewTextView()
	a.summaryList = tview.NewList()
	a.summaryTable = tview.NewTable()
	a.diffViewer = tview.NewTextView()
	a.commandViewer = tview.NewTextView()

	// Initialize config changes tracking
	a.configChanges = make(map[string]ResourceConfigChange)

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

	// Configure editor help
	a.editorHelp.SetBorder(false)
	a.editorHelp.SetDynamicColors(true)
	a.editorHelp.SetText("[yellow]Ctrl+S[white]: Save | [yellow]Esc[white]: Close Editor | [yellow]Tab[white]: Navigate")
	a.editorHelp.SetTextAlign(tview.AlignCenter)

	// Configure versions view
	a.versionsView.SetBorder(true)
	a.versionsView.SetTitle("Target Service Plan Versions")
	a.versionsView.SetTitleAlign(tview.AlignLeft)
	a.versionsView.SetSelectable(true, false)
	a.versionsView.SetSeparator('|')

	// Configure resources view
	a.resourcesView.SetBorder(true)
	a.resourcesView.SetTitle("Resources")
	a.resourcesView.SetTitleAlign(tview.AlignLeft)
	a.resourcesView.ShowSecondaryText(true)

	// Configure summary view 2 (for patch summary)
	a.summaryView2.SetBorder(true)
	a.summaryView2.SetTitle("Patch Summary")
	a.summaryView2.SetTitleAlign(tview.AlignLeft)
	a.summaryView2.SetDynamicColors(true)
	a.summaryView2.SetWordWrap(true)

	// Configure summary list (for interactive selection)
	a.summaryList.SetBorder(true)
	a.summaryList.SetTitle("Modified Resources")
	a.summaryList.SetTitleAlign(tview.AlignLeft)
	a.summaryList.ShowSecondaryText(true)

	// Configure summary table (for patch summary)
	a.summaryTable.SetBorder(true)
	a.summaryTable.SetTitle("Modified Resources")
	a.summaryTable.SetTitleAlign(tview.AlignLeft)
	a.summaryTable.SetSelectable(true, false) // Rows selectable, columns not
	a.summaryTable.SetSeparator('|')
	a.summaryTable.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite))

	// Configure diff viewer
	a.diffViewer.SetBorder(true)
	a.diffViewer.SetTitle("Configuration Diff")
	a.diffViewer.SetTitleAlign(tview.AlignLeft)
	a.diffViewer.SetDynamicColors(true)
	a.diffViewer.SetWordWrap(false)

	// Set up input capture for diff viewer
	a.diffViewer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Return focus to summary table when in summary mode
			if a.sidebarMode == "summary" {
				a.app.SetFocus(a.summaryTable)
				return nil
			}
		}
		return event
	})

	// Configure command viewer
	a.commandViewer.SetBorder(true)
	a.commandViewer.SetTitle("Command Execution")
	a.commandViewer.SetTitleAlign(tview.AlignLeft)
	a.commandViewer.SetDynamicColors(true)
	a.commandViewer.SetWordWrap(false)
	a.commandViewer.SetScrollable(true)

	// Configure status bar
	a.statusBar.SetDynamicColors(true)
	a.statusBar.SetText("[yellow]Press 'q' to quit, 'Tab' to navigate, 'Enter' to select[white]")

	// Initialize sidebar mode and content
	a.sidebarMode = "operations"
	a.setupSidebar()

	// Setup summary content
	a.setupSummary()

	// Create FIXED layout structure that never changes
	// Only the content area will be swapped
	a.contentArea = tview.NewFlex()
	a.showSummaryView() // Set initial content

	contentFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(a.sidebar, 35, 1, true).
		AddItem(a.contentArea, 0, 2, false)

	a.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentFlex, 0, 1, true).
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
	case "versions":
		a.sidebar.SetTitle("Target Versions")
		// Add back to operations option
		a.sidebar.AddItem("← Back to Operations", "", 'b', func() {
			a.sidebarMode = "operations"
			a.setupSidebar()
			a.showSummaryView()
			a.app.SetFocus(a.sidebar)
		})
	case "resourceSelect":
		a.sidebar.SetTitle("Resource Navigation")
		a.sidebar.AddItem("📋 Summary & Finish", "", 's', func() {
			a.showPatchSummary()
		})
		a.sidebar.AddItem("← Back to Versions", "", 'b', func() {
			a.sidebarMode = "versions"
			a.setupSidebar()
			a.showVersionsView()
			a.app.SetFocus(a.sidebar)
		})
	case "summary":
		a.sidebar.SetTitle("Patch Actions")
		a.sidebar.AddItem("✅ Finish & Apply Patch", "", 'f', func() {
			a.finalizePatch()
		})
		a.sidebar.AddItem("← Back to Resources", "", 'b', func() {
			a.sidebarMode = "resourceSelect"
			a.setupSidebar()
			a.showResourceSelectionView()
			a.app.SetFocus(a.sidebar)
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

// showSummaryView displays the summary view in the content area
func (a *OpsCenterApp) showSummaryView() {
	a.contentArea.Clear()
	a.contentArea.AddItem(a.summaryView, 0, 1, false)
}

// showEditorView displays the editor in the content area
func (a *OpsCenterApp) showEditorView() {
	a.contentArea.Clear()
	editorFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.configEditor, 0, 1, true).
		AddItem(a.editorHelp, 1, 1, false)
	a.contentArea.AddItem(editorFlex, 0, 1, true)
}

func (a *OpsCenterApp) showInstanceOverview() {
	// Switch back to overview if showing other panels
	if a.showingEditor {
		a.hideEditor()
	}

	// Ensure we're in operations mode
	if a.sidebarMode != "operations" {
		a.sidebarMode = "operations"
		a.setupSidebar()
	}

	a.updateStatusBar("Instance overview - Press 'p' for patch options")
}

func (a *OpsCenterApp) showPatchOptions() {
	// Switch to versions mode to select target version
	a.sidebarMode = "versions"
	a.setupSidebar()
	a.loadServicePlanVersions()
	a.updateStatusBar("Select a target service plan version for the one-off patch")
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

	// Simply swap the content area to show editor
	a.showEditorView()

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

	// Clear the editor state
	a.configEditor.SetText("", false)
	a.configEditor.SetTitle("Helm Configuration Editor")
	a.configEditor.SetChangedFunc(nil)

	// Return to appropriate view based on mode
	if a.sidebarMode == "resourceSelect" {
		// Return to resource selection view to allow editing more resources
		a.showResourceSelectionView()
		a.app.SetFocus(a.resourcesView)
		a.updateStatusBar("Select another resource to edit or press ESC to go to sidebar")
	} else {
		// Legacy mode - return to summary
		a.showSummaryView()
		// Return to operations mode in sidebar
		a.sidebarMode = "operations"
		a.setupSidebar()
		a.app.SetFocus(a.sidebar)
		a.updateStatusBar("Press 'q' to quit, 'Tab' to navigate, 'Enter' to select")
	}
}

func (a *OpsCenterApp) showSaveDiscardDialog() {
	modal := tview.NewModal().
		SetText("You have unsaved changes.\n\nWhat would you like to do?").
		AddButtons([]string{"Save", "Discard", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// Return to main UI first
			a.app.SetRoot(a.mainFlex, true)

			switch buttonLabel {
			case "Save":
				a.saveConfiguration()
			case "Discard":
				a.closeEditor()
			case "Cancel":
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

	// Update config changes tracking
	if change, exists := a.configChanges[a.selectedResource]; exists {
		change.ModifiedConfig = a.modifiedConfig
		change.IsConfigured = a.modifiedConfig != change.OriginalConfig
		a.configChanges[a.selectedResource] = change
	}

	// Update original config and clear changes flag
	a.originalConfig = a.modifiedConfig
	a.hasChanges = false

	// Save to temp file (for one-off patch mode) or show legacy confirmation
	if a.sidebarMode == "resourceSelect" {
		a.updateStatusBar("[green]Configuration saved! Continue editing or go to Summary[white]")
		a.closeEditor()
	} else {
		// Legacy single-resource mode
		if err := a.saveToConfigFile(); err != nil {
			a.updateStatusBar(fmt.Sprintf("[red]Error saving config: %v[white]", err))
			return
		}
		a.showConfirmationDialog()
	}
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

	patchConfig := map[string]openapiclientfleet.ResourceOneOffPatchConfigurationOverride{
		a.selectedResource: {
			HelmChartValues: helmValues,
		},
	}

	configData, err := yaml.Marshal(patchConfig)
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
			// Return to main UI first
			a.app.SetRoot(a.mainFlex, true)

			if buttonLabel == "Yes" {
				a.startPatch()
			}
			a.closeEditor()
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
			// Handle ESC in diff viewer to return to summary
			if a.sidebarMode == "summary" && a.app.GetFocus() == a.diffViewer {
				a.app.SetFocus(a.summaryTable)
				return nil
			}
			return event
		case tcell.KeyCtrlC:
			a.app.Stop()
			return nil
		case tcell.KeyCtrlS:
			// Handle Ctrl+S for saving in editor
			if a.showingEditor {
				a.saveConfiguration()
				return nil
			}
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
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight, tcell.KeyEnter:
			// Let the focused component handle navigation keys
			if a.sidebarMode == "summary" && a.app.GetFocus() == a.summaryTable {
				return event
			}
			return event
		case tcell.KeyTab:
			if a.showingEditor {
				return event // Let the editor handle tab
			}
			// Handle tab in resource selection mode to switch between new/existing
			if a.sidebarMode == "resourceSelect" {
				current := a.app.GetFocus()
				if current == a.resourcesView {
					// Switch between new and existing resources
					if a.resourceMode == "existing" {
						a.resourceMode = "new"
					} else {
						a.resourceMode = "existing"
					}
					a.populateResourcesList()
					return nil
				}
			}
			// Cycle focus between components
			current := a.app.GetFocus()
			if current == a.sidebar {
				switch a.sidebarMode {
				case "versions":
					a.app.SetFocus(a.versionsView)
				case "resourceSelect":
					a.app.SetFocus(a.resourcesView)
				case "summary":
					a.app.SetFocus(a.summaryTable)
				default:
					a.app.SetFocus(a.summaryView)
				}
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

// compareVersions compares two version strings and returns:
// >0 if version1 > version2
// <0 if version1 < version2
// 0 if version1 == version2
func compareVersions(version1, version2 string) int {
	// Parse version strings into components
	v1Parts := parseVersionParts(version1)
	v2Parts := parseVersionParts(version2)

	// Compare each component
	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		v1Val := 0
		v2Val := 0

		if i < len(v1Parts) {
			v1Val = v1Parts[i]
		}
		if i < len(v2Parts) {
			v2Val = v2Parts[i]
		}

		if v1Val > v2Val {
			return 1
		} else if v1Val < v2Val {
			return -1
		}
	}

	return 0
}

// parseVersionParts parses a version string into numeric components
// Examples: "1.2.3" -> [1, 2, 3], "2.1.0-beta" -> [2, 1, 0]
func parseVersionParts(version string) []int {
	// Remove common prefixes like 'v'
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")

	// Split on dots and handle pre-release suffixes
	parts := strings.Split(version, ".")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		// Remove any non-numeric suffixes (like "-beta", "-alpha", etc.)
		numericPart := strings.Split(part, "-")[0]
		numericPart = strings.Split(numericPart, "+")[0]

		if num, err := strconv.Atoi(numericPart); err == nil {
			result = append(result, num)
		} else {
			// If we can't parse as number, treat as 0
			result = append(result, 0)
		}
	}

	return result
}

// loadServicePlanVersions loads eligible service plan versions for one-off patch
func (a *OpsCenterApp) loadServicePlanVersions() {
	go func() {
		// Get current service plan ID from instance
		currentServicePlanID := a.instanceData.ProductTierId

		// Search for all versions of this service plan
		searchRes, err := dataaccess.SearchInventory(context.Background(), a.token, fmt.Sprintf("serviceplan:%s", currentServicePlanID))
		if err != nil {
			a.app.QueueUpdateDraw(func() {
				a.updateStatusBar(fmt.Sprintf("[red]Error loading versions: %v[white]", err))
			})
			return
		}

		// Filter out deprecated versions and sort by release date
		eligibleVersions := make([]openapiclientfleet.ServicePlanSearchRecord, 0)
		for _, version := range searchRes.ServicePlanResults {
			// Only include non-deprecated versions
			if version.VersionSetStatus != "Deprecated" {
				eligibleVersions = append(eligibleVersions, version)
			}
		}

		a.app.QueueUpdateDraw(func() {
			a.populateVersionsList(eligibleVersions)
			a.showVersionsView()
		})
	}()
}

// populateVersionsList populates the versions table with eligible versions
func (a *OpsCenterApp) populateVersionsList(versions []openapiclientfleet.ServicePlanSearchRecord) {
	a.versionsView.Clear()

	// Sort versions by version number in descending order
	sortedVersions := make([]openapiclientfleet.ServicePlanSearchRecord, len(versions))
	copy(sortedVersions, versions)
	sort.Slice(sortedVersions, func(i, j int) bool {
		return compareVersions(sortedVersions[i].Version, sortedVersions[j].Version) > 0
	})

	if len(sortedVersions) == 0 {
		// Add header row
		a.versionsView.SetCell(0, 0, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetSelectable(false))
		a.versionsView.SetCell(0, 1, tview.NewTableCell("Version Name").SetAlign(tview.AlignCenter).SetSelectable(false))
		a.versionsView.SetCell(0, 2, tview.NewTableCell("Version").SetAlign(tview.AlignCenter).SetSelectable(false))
		a.versionsView.SetCell(0, 3, tview.NewTableCell("Release Date").SetAlign(tview.AlignCenter).SetSelectable(false))
		a.versionsView.SetCell(0, 4, tview.NewTableCell("Version Status").SetAlign(tview.AlignCenter).SetSelectable(false))

		// Add no data row
		a.versionsView.SetCell(1, 0, tview.NewTableCell("❌").SetAlign(tview.AlignCenter).SetSelectable(false))
		a.versionsView.SetCell(1, 1, tview.NewTableCell("No eligible versions found").SetAlign(tview.AlignCenter).SetSelectable(false).SetExpansion(4))
		return
	}

	// Add header row
	a.versionsView.SetCell(0, 0, tview.NewTableCell("Status").SetAlign(tview.AlignCenter).SetSelectable(false).SetTextColor(tcell.ColorYellow))
	a.versionsView.SetCell(0, 1, tview.NewTableCell("Version Name").SetAlign(tview.AlignCenter).SetSelectable(false).SetTextColor(tcell.ColorYellow))
	a.versionsView.SetCell(0, 2, tview.NewTableCell("Version").SetAlign(tview.AlignCenter).SetSelectable(false).SetTextColor(tcell.ColorYellow))
	a.versionsView.SetCell(0, 3, tview.NewTableCell("Release Date").SetAlign(tview.AlignCenter).SetSelectable(false).SetTextColor(tcell.ColorYellow))
	a.versionsView.SetCell(0, 4, tview.NewTableCell("Version Status").SetAlign(tview.AlignCenter).SetSelectable(false).SetTextColor(tcell.ColorYellow))

	currentVersion := a.instanceData.TierVersion
	row := 1
	for _, version := range sortedVersions {

		versionName := version.Version
		if version.VersionName != nil && *version.VersionName != "" {
			versionName = *version.VersionName
		}

		statusIcon := "📦"
		statusColor := tcell.ColorWhite
		isCurrentVersion := version.Version == currentVersion

		if isCurrentVersion {
			statusIcon = "🔷" // Current version indicator
			statusColor = tcell.ColorBlue
		} else if version.VersionSetStatus == "Preferred" {
			statusIcon = "⭐"
			statusColor = tcell.ColorGreen
		}

		releaseDate := "Unknown"
		if version.ReleasedAt != nil {
			releaseDate = *version.ReleasedAt
			// Format the date nicely if it's in RFC3339 format
			if len(releaseDate) > 10 {
				releaseDate = releaseDate[:10] // Just the date part
			}
		}

		// Add table row with current version highlighting
		statusCell := tview.NewTableCell(statusIcon).SetAlign(tview.AlignCenter).SetTextColor(statusColor)
		nameCell := tview.NewTableCell(versionName).SetAlign(tview.AlignLeft)
		versionCell := tview.NewTableCell(version.Version).SetAlign(tview.AlignLeft)
		dateCell := tview.NewTableCell(releaseDate).SetAlign(tview.AlignLeft)
		statusTextCell := tview.NewTableCell(version.VersionSetStatus).SetAlign(tview.AlignLeft)

		// Highlight current version row
		if isCurrentVersion {
			statusCell.SetBackgroundColor(tcell.ColorDarkBlue)
			nameCell.SetBackgroundColor(tcell.ColorDarkBlue)
			versionCell.SetBackgroundColor(tcell.ColorDarkBlue)
			dateCell.SetBackgroundColor(tcell.ColorDarkBlue)
			statusTextCell.SetBackgroundColor(tcell.ColorDarkBlue)
			// Also make text white for better contrast
			nameCell.SetTextColor(tcell.ColorWhite)
			versionCell.SetTextColor(tcell.ColorWhite)
			dateCell.SetTextColor(tcell.ColorWhite)
			statusTextCell.SetTextColor(tcell.ColorWhite)
		}

		// Make current version non-selectable
		if isCurrentVersion {
			statusCell.SetSelectable(false)
			nameCell.SetSelectable(false)
			versionCell.SetSelectable(false)
			dateCell.SetSelectable(false)
			statusTextCell.SetSelectable(false)
		}

		a.versionsView.SetCell(row, 0, statusCell)
		a.versionsView.SetCell(row, 1, nameCell)
		a.versionsView.SetCell(row, 2, versionCell)
		a.versionsView.SetCell(row, 3, dateCell)
		a.versionsView.SetCell(row, 4, statusTextCell)

		row++
	}

	// Set up selection handler
	a.versionsView.SetSelectedFunc(func(selectedRow, column int) {
		if selectedRow > 0 { // Skip header row
			versionCell := a.versionsView.GetCell(selectedRow, 2) // Version column
			if versionCell != nil {
				selectedVersion := versionCell.Text
				// Skip current version
				if selectedVersion == currentVersion {
					a.updateStatusBar("[yellow]Cannot select current version for one-off patch[white]")
					return
				}
				a.selectTargetVersion(selectedVersion)
			}
		}
	})

	// Set initial selection to first selectable row (skip current version if it's the first)
	if row > 1 {
		initialRow := 1
		// Skip current version for initial selection
		for r := 1; r < row; r++ {
			versionCell := a.versionsView.GetCell(r, 2)
			if versionCell != nil && versionCell.Text != currentVersion {
				initialRow = r
				break
			}
		}
		a.versionsView.Select(initialRow, 0)
	}
}

// selectTargetVersion handles target version selection and loads resource comparison
func (a *OpsCenterApp) selectTargetVersion(version string) {
	a.targetVersion = version
	a.updateStatusBar(fmt.Sprintf("Loading resource comparison for version %s...", version))

	go func() {
		// Load current version set
		currentVersionSet, err := dataaccess.DescribeVersionSet(context.Background(), a.token, a.serviceID, a.instanceData.ProductTierId, a.instanceData.TierVersion)
		if err != nil {
			a.app.QueueUpdateDraw(func() {
				a.updateStatusBar(fmt.Sprintf("[red]Error loading current version: %v[white]", err))
			})
			return
		}

		// Load target version set
		targetVersionSet, err := dataaccess.DescribeVersionSet(context.Background(), a.token, a.serviceID, a.instanceData.ProductTierId, version)
		if err != nil {
			a.app.QueueUpdateDraw(func() {
				a.updateStatusBar(fmt.Sprintf("[red]Error loading target version: %v[white]", err))
			})
			return
		}

		a.app.QueueUpdateDraw(func() {
			a.currentVersionSet = currentVersionSet
			a.targetVersionSet = targetVersionSet
			a.categorizeResources()
			a.sidebarMode = "resourceSelect"
			a.setupSidebar()
			a.showResourceSelectionView()
			a.updateStatusBar("Select resources to configure or go to Summary to finish")
		})
	}()
}

// categorizeResources categorizes resources into new, existing, and deprecated
func (a *OpsCenterApp) categorizeResources() {
	// Reset categories
	a.newResources = make([]openapiclient.ResourceSummary, 0)
	a.existingResources = make([]openapiclient.ResourceSummary, 0)
	a.deprecatedResources = make([]openapiclient.ResourceSummary, 0)

	// Create maps for easy lookup
	currentResourceMap := make(map[string]openapiclient.ResourceSummary)
	for _, resource := range a.currentVersionSet.Resources {
		currentResourceMap[resource.Id] = resource
	}

	targetResourceMap := make(map[string]openapiclient.ResourceSummary)
	for _, resource := range a.targetVersionSet.Resources {
		targetResourceMap[resource.Id] = resource
	}

	// Find new and existing resources
	for _, targetResource := range a.targetVersionSet.Resources {
		if _, exists := currentResourceMap[targetResource.Id]; exists {
			a.existingResources = append(a.existingResources, targetResource)
		} else {
			a.newResources = append(a.newResources, targetResource)
		}
	}

	// Find deprecated resources
	for _, currentResource := range a.currentVersionSet.Resources {
		if _, exists := targetResourceMap[currentResource.Id]; !exists {
			a.deprecatedResources = append(a.deprecatedResources, currentResource)
		}
	}

	// Initialize with existing resources view
	a.resourceMode = "existing"
}

// showVersionsView displays the versions selection view
func (a *OpsCenterApp) showVersionsView() {
	a.contentArea.Clear()
	a.contentArea.AddItem(a.versionsView, 0, 1, true)
	a.app.SetFocus(a.versionsView)
}

// showResourceSelectionView displays the resource selection UI for one-off patch
func (a *OpsCenterApp) showResourceSelectionView() {
	a.contentArea.Clear()

	// Create instruction text
	instructionText := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetText(fmt.Sprintf(
			"[yellow]One-Off Patch Resource Configuration[white]\n\n"+
				"Target Version: [cyan]%s[white]\n"+
				"Current Version: [gray]%s[white]\n\n"+
				"[green]Tab[white] to switch between New and Existing resources\n"+
				"[green]Enter[white] to configure selected resource\n"+
				"[green]Escape[white] to return to sidebar (finish editing)\n\n"+
				"Current Mode: [cyan]%s[white] resources",
			a.targetVersion, a.instanceData.TierVersion, strings.ToUpper(a.resourceMode)))

	// Create layout with instruction and resource list
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(instructionText, 8, 0, false).
		AddItem(a.resourcesView, 0, 1, true)

	a.contentArea.AddItem(layout, 0, 1, true)
	a.populateResourcesList()
	a.app.SetFocus(a.resourcesView)

	// Set up input handler for resource selection view
	a.resourcesView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Return to summary mode for final actions (keep the current one-off patch flow)
			a.sidebarMode = "summary"
			a.setupSidebar()
			a.showPatchSummary()
			a.app.SetFocus(a.sidebar)
			a.updateStatusBar("Select 'Summary' to review changes or 'Finish' to apply patch")
			return nil
		}
		if event.Key() == tcell.KeyTab {
			// Switch between new and existing resources
			if a.resourceMode == "new" {
				a.resourceMode = "existing"
			} else {
				a.resourceMode = "new"
			}
			a.populateResourcesList()
			a.showResourceSelectionView() // Refresh view with new mode
			return nil
		}
		return event
	})
}

// populateResourcesList populates the resources view based on current mode
func (a *OpsCenterApp) populateResourcesList() {
	a.resourcesView.Clear()

	var resources []openapiclient.ResourceSummary
	var title string
	var instructions string

	switch a.resourceMode {
	case "new":
		resources = a.newResources
		title = "New Resources (Tab: Switch to Existing)"
		instructions = "Select a resource to configure (new resources will be added)"
	case "existing":
		resources = a.existingResources
		title = "Existing Resources (Tab: Switch to New)"
		instructions = "Select a resource to modify its configuration"
	default:
		resources = a.existingResources
		title = "Existing Resources (Tab: Switch to New)"
		instructions = "Select a resource to modify its configuration"
		a.resourceMode = "existing"
	}

	a.resourcesView.SetTitle(title)
	a.updateStatusBar(instructions)

	if len(resources) == 0 {
		message := "No new resources"
		if a.resourceMode == "existing" {
			message = "No existing resources"
		}
		a.resourcesView.AddItem(fmt.Sprintf("❌ %s", message), "", 0, nil)
		return
	}

	// Add deprecated resources info if viewing existing
	if a.resourceMode == "existing" && len(a.deprecatedResources) > 0 {
		a.resourcesView.AddItem(
			fmt.Sprintf("⚠️  %d deprecated resources (will be removed)", len(a.deprecatedResources)),
			"These resources will not be available in the target version",
			0,
			nil)
	}

	for _, resource := range resources {
		resourceName := resource.Name
		resourceID := resource.Id

		// Check if resource has helm configuration (for modification eligibility)
		canModify := a.isResourceConfigurable(resource, a.resourceMode)
		var helmConfig *openapiclient.HelmChartConfiguration

		if canModify {
			// Find current helm config from instance data for existing resources
			if a.resourceMode == "existing" {
				for _, rv := range a.instanceData.ResourceVersionSummaries {
					if rv.ResourceId != nil && *rv.ResourceId == resourceID && rv.HelmDeploymentConfiguration != nil {
						// Convert fleet model to v1 model for compatibility
						helmConfig = convertFleetHelmConfigToV1(rv.HelmDeploymentConfiguration)
						break
					}
				}
			} else {
				// For new resources, create empty helm config for editing
				helmConfig = &openapiclient.HelmChartConfiguration{
					ChartValues: make(map[string]interface{}),
				}
			}
		}

		icon := "❌"
		secondaryText := a.getResourceConfigurabilityMessage(resource, canModify, resourceID)
		if canModify {
			icon = "✅"
			if _, configured := a.configChanges[resourceID]; configured {
				icon = "📝"
				secondaryText += " | Modified"
			}
		}

		// Capture variables for closure
		resourceCopy := resource
		helmConfigCopy := helmConfig
		a.resourcesView.AddItem(
			fmt.Sprintf("%s %s", icon, resourceName),
			secondaryText,
			0,
			func() {
				if canModify {
					a.selectResourceForConfig(resourceCopy, helmConfigCopy)
				} else {
					// Show specific reason why resource is not configurable
					reason := a.getResourceConfigurabilityMessage(resourceCopy, false, resourceCopy.Id)
					a.updateStatusBar(fmt.Sprintf("[red]%s[white]", reason))
				}
			})
	}
}

// convertFleetHelmConfigToV1 converts fleet SDK helm config to v1 SDK format
func convertFleetHelmConfigToV1(fleetConfig *openapiclientfleet.HelmDeploymentConfiguration) *openapiclient.HelmChartConfiguration {
	if fleetConfig == nil {
		return nil
	}

	v1Config := &openapiclient.HelmChartConfiguration{
		ChartValues: fleetConfig.Values,
	}

	return v1Config
}

// selectResourceForConfig handles resource selection for configuration
func (a *OpsCenterApp) selectResourceForConfig(resource openapiclient.ResourceSummary, helmConfig *openapiclient.HelmChartConfiguration) {
	a.selectedResource = resource.Id

	// Check if we already have changes for this resource
	if existing, ok := a.configChanges[resource.Id]; ok {
		a.originalConfig = existing.OriginalConfig
		a.modifiedConfig = existing.ModifiedConfig
	} else {
		// Extract and format helm values
		if helmConfig == nil || helmConfig.ChartValues == nil {
			a.updateStatusBar("[red]No helm values found for this resource[white]")
			return
		}

		// Convert values to YAML
		yamlData, err := yaml.Marshal(helmConfig.ChartValues)
		if err != nil {
			a.updateStatusBar(fmt.Sprintf("[red]Error formatting helm values: %v[white]", err))
			return
		}

		a.originalConfig = string(yamlData)
		a.modifiedConfig = a.originalConfig

		// Initialize config change tracking
		a.configChanges[resource.Id] = ResourceConfigChange{
			ResourceID:     resource.Id,
			ResourceName:   resource.Name,
			ResourceType:   a.resourceMode,
			OriginalConfig: a.originalConfig,
			ModifiedConfig: a.originalConfig,
			IsConfigured:   false,
		}
	}

	// Show editor
	a.showEditor(resource.Name)
}

// showPatchSummary displays the patch summary with all changes
func (a *OpsCenterApp) showPatchSummary() {
	a.sidebarMode = "summary"
	a.setupSidebar()
	a.populatePatchSummary()
	// populatePatchSummary() calls populateResourceTable() which sets up the complete layout
	a.updateStatusBar("Use Arrow Keys to navigate table, Enter to view resource diff, or select 'Finish' to apply patch")
}

// populatePatchSummary populates the patch summary view
func (a *OpsCenterApp) populatePatchSummary() {
	var summary strings.Builder

	summary.WriteString("[yellow]One-Off Patch Summary[white]\n\n")
	summary.WriteString(fmt.Sprintf("Target Version: [cyan]%s[white]\n", a.targetVersion))
	summary.WriteString(fmt.Sprintf("Current Version: [gray]%s[white]\n\n", a.instanceData.TierVersion))

	// Show resource counts
	summary.WriteString("[yellow]Resource Changes:[white]\n")
	summary.WriteString(fmt.Sprintf("• New Resources: [green]%d[white]\n", len(a.newResources)))
	summary.WriteString(fmt.Sprintf("• Existing Resources: [blue]%d[white]\n", len(a.existingResources)))
	summary.WriteString(fmt.Sprintf("• Deprecated Resources: [red]%d[white]\n\n", len(a.deprecatedResources)))

	// Show configured resources
	configuredCount := 0
	for _, change := range a.configChanges {
		if change.IsConfigured {
			configuredCount++
		}
	}

	if configuredCount > 0 {
		summary.WriteString("[yellow]Configuration Changes:[white]\n")
		summary.WriteString(fmt.Sprintf("• Modified Resources: [cyan]%d[white]\n\n", configuredCount))
	} else {
		summary.WriteString("[yellow]Configuration Changes:[white]\n")
		summary.WriteString("• [gray]No configuration changes made[white]\n")
		summary.WriteString("• [gray]Changes will be applied from target plan template[white]\n\n")
	}

	summary.WriteString("[gray]Select a resource below to view details or diff[white]\n")

	a.summaryView2.SetText(summary.String())
	a.populateResourceTable()
}

// populateResourceTable creates a combined table of all resources
func (a *OpsCenterApp) populateResourceTable() {
	// Use the summaryTable from the struct
	a.summaryTable.Clear()
	a.summaryTable.SetBorders(true)
	a.summaryTable.SetSelectable(true, false)
	a.summaryTable.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite))

	// Set headers
	headers := []string{"Type", "Name", "Status", "Configuration"}
	for i, header := range headers {
		a.summaryTable.SetCell(0, i, tview.NewTableCell(header).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	}

	// Create resource data structures for selection handling
	type rowData struct {
		resource     openapiclient.ResourceSummary
		resourceType string
		change       ResourceConfigChange
		hasChange    bool
	}

	resourceMap := make(map[int]rowData)
	row := 1

	// Add new resources
	for _, resource := range a.newResources {
		change, hasChange := a.configChanges[resource.Id]
		statusText := "NEW"
		configText := a.getResourceConfigText(resource, "new", hasChange && change.IsConfigured)

		a.summaryTable.SetCell(row, 0, tview.NewTableCell("[green]NEW[white]").SetTextColor(tcell.ColorGreen))
		a.summaryTable.SetCell(row, 1, tview.NewTableCell(resource.Name))
		a.summaryTable.SetCell(row, 2, tview.NewTableCell(statusText))
		a.summaryTable.SetCell(row, 3, tview.NewTableCell(configText))

		resourceMap[row] = rowData{
			resource:     resource,
			resourceType: "new",
			change:       change,
			hasChange:    hasChange,
		}
		row++
	}

	// Add existing resources
	for _, resource := range a.existingResources {
		change, hasChange := a.configChanges[resource.Id]
		statusText := "EXISTING"
		configText := a.getResourceConfigText(resource, "existing", hasChange && change.IsConfigured)

		a.summaryTable.SetCell(row, 0, tview.NewTableCell("[blue]EXISTING[white]").SetTextColor(tcell.ColorBlue))
		a.summaryTable.SetCell(row, 1, tview.NewTableCell(resource.Name))
		a.summaryTable.SetCell(row, 2, tview.NewTableCell(statusText))
		a.summaryTable.SetCell(row, 3, tview.NewTableCell(configText))

		resourceMap[row] = rowData{
			resource:     resource,
			resourceType: "existing",
			change:       change,
			hasChange:    hasChange,
		}
		row++
	}

	// Add deprecated resources
	for _, resource := range a.deprecatedResources {
		a.summaryTable.SetCell(row, 0, tview.NewTableCell("[red]DEPRECATED[white]").SetTextColor(tcell.ColorRed))
		a.summaryTable.SetCell(row, 1, tview.NewTableCell(resource.Name))
		a.summaryTable.SetCell(row, 2, tview.NewTableCell("TO BE REMOVED"))
		a.summaryTable.SetCell(row, 3, tview.NewTableCell("Will be deleted"))

		resourceMap[row] = rowData{
			resource:     resource,
			resourceType: "deprecated",
			change:       ResourceConfigChange{},
			hasChange:    false,
		}
		row++
	}

	// Set up selection handler for all rows
	a.summaryTable.SetSelectedFunc(func(selectedRow, selectedCol int) {
		if selectedRow == 0 {
			return // Skip header row
		}

		if data, exists := resourceMap[selectedRow]; exists {
			if data.hasChange && data.change.IsConfigured {
				a.showResourceDiff(data.change)
			} else {
				var message string
				switch data.resourceType {
				case "existing":
					message = "No configuration changes"
				case "deprecated":
					message = "This resource will be removed"
				default:
					message = "Template configuration will be applied"
				}
				a.showResourceInfo(data.resource, data.resourceType, message)
			}
		}
	})

	// Replace the summary list with the table
	a.summaryList.Clear()
	a.contentArea.Clear()

	// Create a layout with summary info and resource table
	leftPane := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.summaryView2, 0, 1, false).
		AddItem(a.summaryTable, 0, 2, true)

	summaryFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPane, 0, 1, true).
		AddItem(a.diffViewer, 0, 1, false)

	a.contentArea.AddItem(summaryFlex, 0, 1, true)

	// Set initial selection to first data row (skip header)
	if row > 1 {
		a.summaryTable.Select(1, 0)
	}
	a.app.SetFocus(a.summaryTable)
}

// showResourceInfo displays basic resource information
func (a *OpsCenterApp) showResourceInfo(resource openapiclient.ResourceSummary, resourceType, message string) {
	var info strings.Builder

	info.WriteString(fmt.Sprintf("[yellow]Resource Information - %s[white]\n\n", resource.Name))
	info.WriteString(fmt.Sprintf("Resource ID: [cyan]%s[white]\n", resource.Id))
	info.WriteString(fmt.Sprintf("Resource Type: [cyan]%s[white]\n", strings.ToUpper(resourceType)))
	info.WriteString(fmt.Sprintf("Description: [gray]%s[white]\n\n", resource.Description))

	// Check configurability and provide detailed explanation
	canConfigure := a.isResourceConfigurable(resource, resourceType)
	
	switch resourceType {
	case "new":
		info.WriteString("[green]Status: NEW RESOURCE[white]\n")
		info.WriteString("This resource will be added in the target version.\n\n")
		
		if canConfigure {
			info.WriteString("[green]✅ Configurable[white]\n")
			info.WriteString("You can customize this resource's Helm values.\n")
			if message != "" {
				info.WriteString(fmt.Sprintf("Current state: [cyan]%s[white]\n", message))
			}
		} else {
			info.WriteString("[red]❌ Not Configurable[white]\n")
			if resource.IsExternal {
				info.WriteString("This is an external resource managed outside the Kubernetes cluster.\n")
			} else if resource.ManagedResourceType != nil {
				switch *resource.ManagedResourceType {
				case "LoadBalancer":
					info.WriteString("Load balancer resources use auto-generated configuration based on service requirements.\n")
				case "Storage":
					info.WriteString("Storage resources have fixed configuration that cannot be overridden.\n")
				default:
					info.WriteString("This resource type does not support configuration overrides.\n")
				}
			} else {
				info.WriteString("This resource type does not support configuration overrides.\n")
			}
		}
		
	case "existing":
		info.WriteString("[blue]Status: EXISTING RESOURCE[white]\n")
		info.WriteString("This resource exists in both versions.\n\n")
		
		if canConfigure {
			info.WriteString("[green]✅ Configurable[white]\n")
			info.WriteString("You can modify this resource's Helm configuration values.\n")
			if message != "" {
				info.WriteString(fmt.Sprintf("Current state: [cyan]%s[white]\n", message))
			}
		} else {
			info.WriteString("[red]❌ Not Configurable[white]\n")
			info.WriteString("This resource does not have Helm configuration available for modification.\n")
			info.WriteString("It may be managed through other deployment methods (e.g., raw Kubernetes manifests, operators).\n")
		}
		
	case "deprecated":
		info.WriteString("[red]Status: DEPRECATED RESOURCE[white]\n")
		info.WriteString("This resource will be removed in the target version.\n\n")
		info.WriteString("[red]❌ Not Configurable[white]\n")
		info.WriteString("Resources being removed cannot be configured.\n")
	}

	info.WriteString("\n[gray]Press Esc to return to summary[white]")

	a.diffViewer.SetTitle(fmt.Sprintf("Resource Info - %s", resource.Name))
	a.diffViewer.SetText(info.String())
	a.app.SetFocus(a.diffViewer)
}

// showResourceDiff displays a git-style diff for a resource configuration
func (a *OpsCenterApp) showResourceDiff(change ResourceConfigChange) {
	diff := a.generateGitStyleDiff(change)
	a.diffViewer.SetTitle(fmt.Sprintf("Configuration Diff - %s", change.ResourceName))
	a.diffViewer.SetText(diff)
	a.app.SetFocus(a.diffViewer)
}

// generateGitStyleDiff generates a git-style diff for configuration changes
func (a *OpsCenterApp) generateGitStyleDiff(change ResourceConfigChange) string {
	var diff strings.Builder

	// Header
	diff.WriteString(fmt.Sprintf("[yellow]diff --git a/%s b/%s[white]\n", change.ResourceName, change.ResourceName))
	diff.WriteString("[yellow]index 0000000..1111111 100644[white]\n")
	diff.WriteString(fmt.Sprintf("[yellow]--- a/%s[white]\n", change.ResourceName))
	diff.WriteString(fmt.Sprintf("[yellow]+++ b/%s[white]\n", change.ResourceName))
	diff.WriteString("[cyan]@@ -1,1 +1,1 @@[white]\n")

	switch change.ResourceType {
	case "new":
		// New resource - show all as additions
		lines := strings.Split(change.ModifiedConfig, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				diff.WriteString(fmt.Sprintf("[green]+%s[white]\n", line))
			}
		}
	case "deprecated":
		// Deprecated resource - show all as deletions
		lines := strings.Split(change.OriginalConfig, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				diff.WriteString(fmt.Sprintf("[red]-%s[white]\n", line))
			}
		}
	default:
		// Existing resource - show actual diff
		diff.WriteString(a.generateDetailedDiff(change.OriginalConfig, change.ModifiedConfig))
	}

	diff.WriteString("\n[gray]Press Esc to return to summary[white]")
	return diff.String()
}

// generateDetailedDiff generates a detailed line-by-line diff
func (a *OpsCenterApp) generateDetailedDiff(original, modified string) string {
	var diff strings.Builder

	originalLines := strings.Split(original, "\n")
	modifiedLines := strings.Split(modified, "\n")

	// Simple line-by-line comparison
	maxLines := len(originalLines)
	if len(modifiedLines) > maxLines {
		maxLines = len(modifiedLines)
	}

	for i := 0; i < maxLines; i++ {
		var origLine, modLine string
		if i < len(originalLines) {
			origLine = originalLines[i]
		}
		if i < len(modifiedLines) {
			modLine = modifiedLines[i]
		}

		if origLine != modLine {
			if origLine != "" {
				diff.WriteString(fmt.Sprintf("[red]-%s[white]\n", origLine))
			}
			if modLine != "" {
				diff.WriteString(fmt.Sprintf("[green]+%s[white]\n", modLine))
			}
		} else if origLine != "" {
			// Unchanged lines in context
			diff.WriteString(fmt.Sprintf(" %s\n", origLine))
		}
	}

	return diff.String()
}

// finalizePatch finalizes the patch and triggers the one-off patch command
func (a *OpsCenterApp) finalizePatch() {
	// Save configuration files and show confirmation
	if err := a.saveAllConfigChanges(); err != nil {
		a.updateStatusBar(fmt.Sprintf("[red]Error saving configurations: %v[white]", err))
		return
	}

	// Show final confirmation dialog
	a.showFinalConfirmationDialog()
}

// saveAllConfigChanges saves all configuration changes to files
func (a *OpsCenterApp) saveAllConfigChanges() error {
	// Save to current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Get resource ID to key mapping
	resourceIDToKeyMap, err := a.getResourceIDToKeyMapping()
	if err != nil {
		return fmt.Errorf("failed to get resource key mapping: %w", err)
	}

	// Build resource override configuration - only include configured resources
	resourceOverrides := make(map[string]openapiclientfleet.ResourceOneOffPatchConfigurationOverride)

	for resourceID, change := range a.configChanges {
		// Only include resources that have been configured with changes
		if !change.IsConfigured {
			continue
		}

		// Convert modified config to interface{}
		var parsedConfig interface{}
		if err := yaml.Unmarshal([]byte(change.ModifiedConfig), &parsedConfig); err != nil {
			return fmt.Errorf("invalid configuration for resource %s: %w", resourceID, err)
		}

		// Create configuration in the format expected by one-off patch
		helmValues, ok := parsedConfig.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid configuration format for resource %s", resourceID)
		}

		// Get the resource key for this resource ID
		resourceKey, exists := resourceIDToKeyMap[resourceID]
		if !exists {
			return fmt.Errorf("resource key not found for resource ID %s", resourceID)
		}

		resourceOverrides[resourceKey] = openapiclientfleet.ResourceOneOffPatchConfigurationOverride{
			HelmChartValues: helmValues,
		}
	}

	// Save configuration as YAML (empty map if no changes)
	configData, err := yaml.Marshal(resourceOverrides)
	if err != nil {
		return err
	}

	// Save to current directory
	filename := filepath.Join(cwd, "instance-patch-plan.yaml")
	if err := os.WriteFile(filename, configData, 0600); err != nil {
		return err
	}

	// Store filename for later use
	a.configFilename = filename

	if len(resourceOverrides) == 0 {
		a.updateStatusBar(fmt.Sprintf("[green]Empty configuration saved (template defaults will be used): %s[white]", filename))
	} else {
		a.updateStatusBar(fmt.Sprintf("[green]Configuration saved to: %s[white]", filename))
	}
	return nil
}

// showFinalConfirmationDialog shows the final confirmation before applying patch
func (a *OpsCenterApp) showFinalConfirmationDialog() {
	configuredCount := 0
	for _, change := range a.configChanges {
		if change.IsConfigured {
			configuredCount++
		}
	}

	var message string
	if configuredCount > 0 {
		message = fmt.Sprintf("Ready to apply one-off patch!\n\nTarget Version: %s\nConfiguration Changes: %d\nNew Resources: %d\nDeprecated Resources: %d\n\nProceed with patch?",
			a.targetVersion, configuredCount, len(a.newResources), len(a.deprecatedResources))
	} else {
		message = fmt.Sprintf("Ready to apply one-off patch!\n\nTarget Version: %s\nConfiguration Changes: %d (template defaults will be used)\nNew Resources: %d\nDeprecated Resources: %d\n\nProceed with patch?",
			a.targetVersion, configuredCount, len(a.newResources), len(a.deprecatedResources))
	}

	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Apply Patch", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// Return to main UI first
			a.app.SetRoot(a.mainFlex, true)

			if buttonLabel == "Apply Patch" {
				a.showCommandExecution()
			}
		})

	a.app.SetRoot(modal, true)
}

// showCommandExecution shows the command execution screen and runs the patch
func (a *OpsCenterApp) showCommandExecution() {
	// Clear content area and show command viewer
	a.contentArea.Clear()
	a.contentArea.AddItem(a.commandViewer, 0, 1, true)
	a.app.SetFocus(a.commandViewer)

	// Build and display the command
	configFlag := ""
	if a.configFilename != "" {
		configFlag = fmt.Sprintf(" --configuration-override %s", a.configFilename)
	}

	command := fmt.Sprintf("omctl instance patch %s --target-tier-version %s%s", a.instanceID, a.targetVersion, configFlag)

	var output strings.Builder
	output.WriteString("[yellow]Executing One-Off Patch[white]\n")
	output.WriteString("=================================\n\n")
	output.WriteString(fmt.Sprintf("[cyan]Command:[white] %s\n\n", command))
	output.WriteString("[yellow]Output:[white]\n")
	output.WriteString("--------\n")

	a.commandViewer.SetText(output.String())
	a.updateStatusBar("[yellow]Executing command... Press 'q' to return to main menu[white]")

	// Execute the command in a goroutine
	go a.executePatchCommand(&output)
}

// executePatchCommand executes the patch command and shows real-time output
func (a *OpsCenterApp) executePatchCommand(output *strings.Builder) {
	// Build command args securely
	args := []string{"instance", "patch", a.instanceID, "--target-tier-version", a.targetVersion}
	if a.configFilename != "" {
		args = append(args, "--configuration-override", a.configFilename)
	}

	cmd := exec.Command("omctl", args...)

	// Set up pipes for output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			_, _ = fmt.Fprintf(output, "[red]Error setting up stdout: %v[white]\n", err)
			a.commandViewer.SetText(output.String())
		})
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.app.QueueUpdateDraw(func() {
			_, _ = fmt.Fprintf(output, "[red]Error setting up stderr: %v[white]\n", err)
			a.commandViewer.SetText(output.String())
		})
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		a.app.QueueUpdateDraw(func() {
			_, _ = fmt.Fprintf(output, "[red]Error starting command: %v[white]\n", err)
			a.commandViewer.SetText(output.String())
		})
		return
	}

	// Read output in real-time
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			a.app.QueueUpdateDraw(func() {
				output.WriteString(line + "\n")
				a.commandViewer.SetText(output.String())
				a.commandViewer.ScrollToEnd()
			})
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			a.app.QueueUpdateDraw(func() {
				_, _ = fmt.Fprintf(output, "[red]%s[white]\n", line)
				a.commandViewer.SetText(output.String())
				a.commandViewer.ScrollToEnd()
			})
		}
	}()

	// Wait for command to finish
	err = cmd.Wait()
	a.app.QueueUpdateDraw(func() {
		if err != nil {
			_, _ = fmt.Fprintf(output, "\n[red]Command failed with error: %v[white]\n", err)
		} else {
			output.WriteString("\n[green]Command completed successfully![white]\n")
		}
		output.WriteString("\n[gray]Press 'q' to return to main menu[white]")
		a.commandViewer.SetText(output.String())
		a.commandViewer.ScrollToEnd()
		a.updateStatusBar("[green]Command execution finished. Press 'q' to return to main menu[white]")
	})
}

// getResourceIDToKeyMapping creates a mapping from resource ID to resource key (UrlKey)
func (a *OpsCenterApp) getResourceIDToKeyMapping() (map[string]string, error) {
	resourceIDToKeyMap := make(map[string]string)

	// Get the target version set to find product tier ID
	if a.targetVersionSet == nil {
		return nil, fmt.Errorf("target version set not available")
	}

	// Get service offering details for the target version
	offering, err := dataaccess.DescribeServiceOffering(
		context.Background(),
		a.token,
		a.serviceID,
		a.targetVersionSet.ProductTierId,
		a.targetVersion,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get service offering: %w", err)
	}

	// Extract resource parameter details to get UrlKey
	if offering.ConsumptionDescribeServiceOfferingResult == nil ||
		len(offering.ConsumptionDescribeServiceOfferingResult.Offerings) == 0 {
		return nil, fmt.Errorf("no service offering found")
	}

	// Map resource ID to UrlKey
	for _, resourceEntity := range offering.ConsumptionDescribeServiceOfferingResult.Offerings[0].ResourceParameters {
		resourceIDToKeyMap[resourceEntity.ResourceId] = resourceEntity.UrlKey
	}

	return resourceIDToKeyMap, nil
}

// getResourceConfigurabilityMessage returns an appropriate message explaining why a resource is or isn't configurable
func (a *OpsCenterApp) getResourceConfigurabilityMessage(resource openapiclient.ResourceSummary, canModify bool, resourceID string) string {
	if canModify {
		return fmt.Sprintf("ID: %s | Configurable", resourceID)
	}

	// Determine why the resource is not configurable
	switch a.resourceMode {
	case "existing":
		// Check if resource exists in instance but has no helm configuration
		found := false
		for _, rv := range a.instanceData.ResourceVersionSummaries {
			if rv.ResourceId != nil && *rv.ResourceId == resourceID {
				found = true
				if rv.HelmDeploymentConfiguration == nil {
					return "Not configurable: No Helm configuration available"
				}
				break
			}
		}
		if !found {
			return "Not configurable: Resource not found in current instance"
		}
		return "Not configurable: Unknown reason"

	case "new":
		// New resources should generally be configurable, but check for specific cases
		if resource.IsExternal {
			return "Not configurable: External resource (managed outside cluster)"
		}
		// Check resource type for other non-configurable types
		if resource.ManagedResourceType != nil {
			switch *resource.ManagedResourceType {
			case "LoadBalancer":
				return "Not configurable: Load balancer resources use auto-generated configuration"
			case "Storage":
				return "Not configurable: Storage resources have fixed configuration"
			}
		}
		return "Not configurable: Resource type does not support configuration overrides"

	default:
		// For deprecated resources
		return "Not configurable: Resource will be removed in target version"
	}
}

// getResourceConfigText returns appropriate configuration text for the summary table
func (a *OpsCenterApp) getResourceConfigText(resource openapiclient.ResourceSummary, resourceType string, isConfigured bool) string {
	// Check if resource is configurable first
	canModify := a.isResourceConfigurable(resource, resourceType)
	
	if !canModify {
		// Return a brief reason why it's not configurable
		switch resourceType {
		case "existing":
			return "Not configurable (no Helm config)"
		case "new":
			if resource.IsExternal {
				return "Not configurable (external)"
			}
			if resource.ManagedResourceType != nil {
				switch *resource.ManagedResourceType {
				case "LoadBalancer":
					return "Not configurable (auto-generated)"
				case "Storage":
					return "Not configurable (fixed config)"
				}
			}
			return "Not configurable"
		default:
			return "Will be deleted"
		}
	}
	
	// Resource is configurable
	if isConfigured {
		switch resourceType {
		case "existing":
			return "Modified"
		case "new":
			return "Custom configured"
		default:
			return "Configured"
		}
	}
	
	// Resource is configurable but not modified
	switch resourceType {
	case "existing":
		return "No changes"
	case "new":
		return "Template default"
	default:
		return "Unchanged"
	}
}

// isResourceConfigurable checks if a resource can be configured
func (a *OpsCenterApp) isResourceConfigurable(resource openapiclient.ResourceSummary, resourceType string) bool {
	switch resourceType {
	case "existing":
		// Check if resource has helm configuration in instance data
		for _, rv := range a.instanceData.ResourceVersionSummaries {
			if rv.ResourceId != nil && *rv.ResourceId == resource.Id && rv.HelmDeploymentConfiguration != nil {
				return true
			}
		}
		return false
	case "new":
		// Check for specific non-configurable resource types
		if resource.IsExternal {
			return false
		}
		if resource.ManagedResourceType != nil {
			switch *resource.ManagedResourceType {
			case "LoadBalancer", "Storage":
				return false
			}
		}
		return true // New resources are generally configurable
	default:
		return false // Deprecated resources are not configurable
	}
}
