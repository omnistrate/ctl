package instance
import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "strings"
    "time"

    // WebSocket client for log streaming
    "github.com/gorilla/websocket"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"

    "github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
    "github.com/omnistrate-oss/omnistrate-ctl/internal/config"
    "github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
    "github.com/omnistrate-oss/omnistrate-ctl/internal/utils"
    "github.com/spf13/cobra"
)
// LogStream represents a pod or log source for TUI
type LogStream struct {
    PodName string
    LogsURL string
}

// launchLogsTUI displays a TUI for selecting pods and viewing logs
func launchLogsTUI(instanceID string, logStreams []LogStream) error {
    app := tview.NewApplication()

    // Left panel: list of pods (use tview.List for AddItem/SetSelectedFunc)
    leftPanel := tview.NewList()
    leftPanel.SetBorder(true)
    leftPanel.SetTitle("Pods")

    // Right panel: log output
    rightPanel := tview.NewTextView()
    rightPanel.SetBorder(true)
    rightPanel.SetTitle("Logs")
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

    // Populate left panel with pods
    for i, stream := range logStreams {
        leftPanel.AddItem(stream.PodName, "", rune('a'+i), nil)
    }

    // When a pod is selected, show logs in right panel
    leftPanel.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
        if index < 0 || index >= len(logStreams) {
            return
        }
        rightPanel.SetText(fmt.Sprintf("Connecting to logs for pod: %s...", logStreams[index].PodName))
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
        }(logStreams[index].LogsURL)
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

    app.SetFocus(leftPanel)
    if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
        return fmt.Errorf("failed to run logs TUI: %w", err)
    }
    return nil
}

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

// GetResourceInstanceLogs streams or snapshots logs for a given instance.
// If outputMode == "json", fetches a snapshot and returns as []byte.
// Otherwise, connects to logsSocketURL and streams logs to stdout.
func GetResourceInstanceLogs(ctx context.Context, token, serviceID, environmentID, instanceID, outputMode string) ([]byte, error) {
    instance, err := dataaccess.DescribeResourceInstance(ctx, token, serviceID, environmentID, instanceID)
    if err != nil {
        return nil, fmt.Errorf("failed to describe resource instance: %w", err)
    }

    // Debug: print the instance object
    log.Printf("[DEBUG] instance: %+v\n", instance)
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

    var logsURL string

    // Try to construct logsSocketURL from omnistrateobserv resource
    topologyMap := instance.ConsumptionResourceInstanceResult.DetailedNetworkTopology
    // Find podName from the main resource's nodes list if available
    podName := ""
    for _, topologyRaw := range topologyMap {
        if topologyEntry, ok := topologyRaw.(map[string]interface{}); ok {
            if mainVal, ok := topologyEntry["main"].(bool); ok && mainVal {
                if nodes, ok := topologyEntry["nodes"].([]interface{}); ok && len(nodes) > 0 {
                    if node, ok := nodes[0].(map[string]interface{}); ok {
                        if id, ok := node["id"].(string); ok && id != "" {
                            podName = id
                        }
                    }
                }
            }
        }
    }

    for _, topologyRaw := range topologyMap {
        if topologyEntry, ok := topologyRaw.(map[string]interface{}); ok {
            if rk, ok := topologyEntry["resourceKey"].(string); ok && rk == "omnistrateobserv" {
                if ce, ok := topologyEntry["clusterEndpoint"].(string); ok && ce != "" {
                    parts := strings.SplitN(ce, "@", 2)
                    if len(parts) == 2 {
                        userPass := parts[0]
                        baseURL := parts[1]
                        creds := strings.SplitN(userPass, ":", 2)
                        if len(creds) == 2 {
                            username := creds[0]
                            password := creds[1]
                            // If podName not set from main resource, fallback to nodes.id or podName in this entry
                            if podName == "" {
                                if nid, ok := topologyEntry["nodes.id"].(string); ok && nid != "" {
                                    podName = nid
                                } else if pn, ok := topologyEntry["podName"].(string); ok && pn != "" {
                                    podName = pn
                                }
                            }
                            // Build query string
                            logsURL = fmt.Sprintf("wss://%s/logs?username=%s&password=%s", baseURL, username, password)
                            if podName != "" {
                                logsURL += fmt.Sprintf("&podName=%s", podName)
                            }
                            logsURL += fmt.Sprintf("&instanceId=%s", instanceID)
                            break
                        }
                    }
                }
            }
        }
    }


    log.Printf("[DEBUG] logsURL: %s\n", logsURL)
    if logsURL == "" {
        return nil, errors.New("logsSocketURL not available for this instance")
    }

    if outputMode == "json" {
        // For JSON, connect, read a snapshot, and return as JSON array
        return fetchLogsSnapshot(logsURL)
    } else {
        // For other modes, stream logs live
        return nil, streamLogsToStdout(logsURL)
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

// streamLogsToStdout connects to the websocket and streams logs to stdout
func streamLogsToStdout(logsURL string) error {
    c, _, err := websocket.DefaultDialer.Dial(logsURL, nil)
    if err != nil {
        return fmt.Errorf("failed to connect to logs websocket: %w", err)
    }
    defer c.Close()

    log.Println("Streaming logs. Press Ctrl+C to exit.")
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            return err
        }
        fmt.Println(string(message))
    }
}