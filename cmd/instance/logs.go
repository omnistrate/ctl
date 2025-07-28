package instance

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "time"

    // WebSocket client for log streaming
    "github.com/gdamore/tcell/v2"
    "github.com/gorilla/websocket"
    "github.com/rivo/tview"

    "github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
    "github.com/omnistrate-oss/omnistrate-ctl/internal/config"
    "github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
    "github.com/omnistrate-oss/omnistrate-ctl/internal/utils"
    "github.com/spf13/cobra"
)



const (
    logsExample = `# Stream logs for an instance deployment
omnistrate-ctl instance logs instance-abcd1234

# Get a snapshot of logs in JSON format
omnistrate-ctl instance logs instance-abcd1234 -o json`
)

var logsCmd = &cobra.Command{
    Use:          "logs [instance-id]",
    Short:        "Fetch logs for an instance deployment",
    Long:         `This command streams or snapshots logs for the instance of your service.`,
    Example:      logsExample,
    RunE:         runLogs,
    SilenceUsage: true,
}

func init() {
    logsCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
    // Add logsCmd to your instance command in your main CLI setup
}

// runLogs fetches or streams logs for a given instance using GetResourceInstanceLogs.
func runLogs(cmd *cobra.Command, args []string) error {
    defer config.CleanupArgsAndFlags(cmd, &args)

    instanceID := args[0]

    // Retrieve output flag
    output, err := cmd.Flags().GetString("output")
    if err != nil {
        utils.PrintError(err)
        return err
    }

    // Validate user login
    token, err := common.GetTokenWithLogin()
    if err != nil {
        utils.PrintError(err)
        return err
    }

    // Get serviceID and environmentID for the instance
    serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
    if err != nil {
        utils.PrintError(err)
        return err
    }

    // Fetch or stream logs
    logs, err := GetResourceInstanceLogs(cmd.Context(), token, serviceID, environmentID, instanceID, output)
    if err != nil {
        utils.PrintError(err)
        return err
    }

    // If output is json, print the logs as JSON
    if output == "json" && logs != nil {
        fmt.Println(string(logs))
    }

    return nil
}

// LogStream represents a pod or log source for TUI
type LogStream struct {
    PodName string
    LogsURL string
}

type InstanceDetail struct {
   
    InstanceID       string
    Status           string
}


// GetResourceInstanceLogs streams or snapshots logs for a given instance.
// If outputMode == "json", fetches a snapshot and returns as []byte.
// Otherwise, connects to logsSocketURL and streams logs to stdout.
func GetResourceInstanceLogs(ctx context.Context, token, serviceID, environmentID, instanceID, outputMode string) ([]byte, error) {
    instance, err := dataaccess.DescribeResourceInstance(ctx, token, serviceID, environmentID, instanceID)
    if err != nil {
        return nil, fmt.Errorf("failed to describe resource instance: %w", err)
    }

 
    // Check if logs are enabled via LOGS#INTERNAL feature
    isLogsEnabled := false
    features := instance.ConsumptionResourceInstanceResult.ProductTierFeatures
    if features != nil {
        if featRaw, ok := features["LOGS#INTERNAL"]; ok && featRaw != nil {
            if featMap, ok := featRaw.(map[string]interface{}); ok {
                if enabled, ok := featMap["enabled"].(bool); ok && enabled {
                    isLogsEnabled = true
                }
            }
        }
    }
    if !isLogsEnabled {
        return nil, errors.New("logs are not enabled for this instance")
    }

    topologyMap := instance.ConsumptionResourceInstanceResult.DetailedNetworkTopology
    var logStreams []LogStream

    // Build InstanceDetail from instance.ConsumptionResourceInstanceResult (struct field access)
    instanceDetail := InstanceDetail{
        InstanceID: "",
        Status:     "",
    }
    if instance.ConsumptionResourceInstanceResult.Id != nil {
        instanceDetail.InstanceID = *instance.ConsumptionResourceInstanceResult.Id
    }
    if instance.ConsumptionResourceInstanceResult.Status != nil {
        instanceDetail.Status = *instance.ConsumptionResourceInstanceResult.Status
    }

    // Find omnistrateobserv resource for log endpoint
    var baseURL, username, password string
    for _, topologyRaw := range topologyMap {
        if topologyEntry, ok := topologyRaw.(map[string]interface{}); ok {
            if rk, ok := topologyEntry["resourceKey"].(string); ok && rk == "omnistrateobserv" {
                if ce, ok := topologyEntry["clusterEndpoint"].(string); ok && ce != "" {
                    parts := strings.SplitN(ce, "@", 2)
                    if len(parts) == 2 {
                        userPass := parts[0]
                        baseURL = parts[1]
                        creds := strings.SplitN(userPass, ":", 2)
                        if len(creds) == 2 {
                            username = creds[0]
                            password = creds[1]
                        }
                    }
                }
            }
        }
    }

    // Find all pods in the main resource and build log URLs
    for _, topologyRaw := range topologyMap {
        if topologyEntry, ok := topologyRaw.(map[string]interface{}); ok {
            if mainVal, ok := topologyEntry["main"].(bool); ok && mainVal {
                if nodes, ok := topologyEntry["nodes"].([]interface{}); ok {
                    for _, n := range nodes {
                        if node, ok := n.(map[string]interface{}); ok {
                            if podName, ok := node["id"].(string); ok && podName != "" && baseURL != "" && username != "" && password != "" {
                                logsURL := fmt.Sprintf("wss://%s/logs?username=%s&password=%s&podName=%s&instanceId=%s", baseURL, username, password, podName, instanceID)
                                logStreams = append(logStreams, LogStream{PodName: podName, LogsURL: logsURL})
                            }
                        }
                    }
                }
            }
        }
    }

    // Detect and add generic/simple resources for websocket log streaming
    for _, topologyRaw := range topologyMap {
        if topologyEntry, ok := topologyRaw.(map[string]interface{}); ok {
            // Generic resource detection: not helm, not terraform, but has logs endpoint
            if resourceType, ok := topologyEntry["type"].(string); ok && resourceType == "generic" {
                if logsEndpoint, ok := topologyEntry["logsEndpoint"].(string); ok && logsEndpoint != "" {
                    // If credentials are needed, use omnistrateobserv creds if available
                    logsURL := logsEndpoint
                    if baseURL != "" && username != "" && password != "" {
                        // If logsEndpoint is a path, prepend baseURL
                        if !strings.HasPrefix(logsEndpoint, "ws") {
                            logsURL = fmt.Sprintf("wss://%s%s?username=%s&password=%s&instanceId=%s", baseURL, logsEndpoint, username, password, instanceID)
                        }
                    }
                    name := topologyEntry["resourceKey"].(string)
                    logStreams = append(logStreams, LogStream{PodName: name, LogsURL: logsURL})
                }
            }
        }
    }

    if len(logStreams) == 0 {
        return nil, errors.New("No log streams available for this instance")
    }

    if outputMode == "json" {
        // For JSON, connect, read a snapshot from the first pod/resource
        return fetchLogsSnapshot(logStreams[0].LogsURL)
    } else {
        // For other modes, launch the TUI for interactive log selection
        err := launchLogsTUI(instanceDetail, logStreams)
        return nil, err
    }
}

// fetchLogsSnapshot connects to the websocket, reads available logs, and returns as JSON array
func fetchLogsSnapshot(logsURL string) ([]byte, error) {
    c, _, err := websocket.DefaultDialer.Dial(logsURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to logs websocket: %w", err)
    }
    defer c.Close()

    var logs []string
    c.SetReadDeadline(time.Now().Add(5 * time.Second))
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            break // likely timeout or EOF
        }
        logs = append(logs, string(message))
    }
    return json.Marshal(logs)
}



// launchLogsTUI displays a TUI for selecting pods and viewing logs
func launchLogsTUI(instanceDetail InstanceDetail, logStreams []LogStream) error {
    app := tview.NewApplication()

    // Left panel: show instance and resource names, colorized
    leftPanel := tview.NewList()
    leftPanel.SetBorder(true)
    leftPanel.SetTitle("Resources")

    // Right panel: log output or resource info
    rightPanel := tview.NewTextView()
    rightPanel.SetBorder(true)
    rightPanel.SetTitle("Content")
    rightPanel.SetDynamicColors(true)
    rightPanel.SetWrap(true)
    rightPanel.SetScrollable(true)

    // Help text
    helpText := tview.NewTextView().
        SetText("Navigate: ↑/↓ | Enter: view logs | Esc: go back | q: quit").
        SetTextAlign(tview.AlignCenter).
        SetDynamicColors(true)

    // Layout
    flex := tview.NewFlex()
    flex.AddItem(leftPanel, 0, 1, true)
    flex.AddItem(rightPanel, 0, 2, false)

    mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
    mainFlex.AddItem(flex, 0, 1, true)
    mainFlex.AddItem(helpText, 1, 0, false)

    
    
    leftPanel.AddItem(fmt.Sprintf("[yellow]%s[-]", strings.TrimSpace(instanceDetail.InstanceID)), "", 0, nil)

    // Populate left panel with resource names (pods or generic)
    for i, stream := range logStreams {
        color := "[blue]"
        if i == 0 {
            color = "[green]"
        }
        leftPanel.AddItem(fmt.Sprintf("%s%s[-]", color, stream.PodName), "", rune('a'+i), nil)
    }

    // When a resource is selected, show logs or instance info in right panel
    leftPanel.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
        // index 0 is instance ID, show info/status
        if index == 0 {
            rightPanel.SetTitle("Instance Info")
            rightPanel.SetText(fmt.Sprintf("[yellow]Instance ID:[white] %s\n[green]Status:[white] %s\nSelect a resource to view logs.", instanceDetail.InstanceID, instanceDetail.Status))
            return
        }
        idx := index - 1
        if idx < 0 || idx >= len(logStreams) {
            rightPanel.SetTitle("Content")
            rightPanel.SetText("Select a resource to view logs.")
            return
        }
        rightPanel.SetTitle(fmt.Sprintf("Logs: %s", logStreams[idx].PodName))
        rightPanel.SetText(fmt.Sprintf("Connecting to logs for resource: %s...", logStreams[idx].PodName))
        go func(logsURL string) {
            c, _, err := websocket.DefaultDialer.Dial(logsURL, nil)
            if err != nil {
                app.QueueUpdateDraw(func() {
                    rightPanel.SetText(fmt.Sprintf("Failed to connect: %v", err))
                })
                return
            }
            defer c.Close()
            app.QueueUpdateDraw(func() {
                rightPanel.SetText("")
            })
            for {
                _, message, err := c.ReadMessage()
                if err != nil {
                    app.QueueUpdateDraw(func() {
                        rightPanel.SetText(fmt.Sprintf("Connection closed: %v", err))
                    })
                    break
                }
                app.QueueUpdateDraw(func() {
                    rightPanel.Write([]byte(string(message) + "\n"))
                })
            }
        }(logStreams[idx].LogsURL)
        app.SetFocus(rightPanel)
    })

    // Keyboard navigation
    app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyCtrlC, tcell.KeyRune:
            if event.Rune() == 'q' || event.Rune() == 'Q' {
                app.Stop()
                return nil
            }
        case tcell.KeyEscape:
            app.SetFocus(leftPanel)
            return nil
        }
        return event
    })

    leftPanel.SetCurrentItem(0)
    app.SetFocus(leftPanel)
    // Manually trigger selected function for index 0 so instance info is shown immediately
    if leftPanel.GetItemCount() > 0 {
        if selFunc := leftPanel.GetSelectedFunc(); selFunc != nil {
            mainText, secondaryText := leftPanel.GetItemText(0)
            selFunc(0, mainText, secondaryText, 0)
        }
    }
    if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
        return fmt.Errorf("failed to run logs TUI: %w", err)
    }
    return nil
}
