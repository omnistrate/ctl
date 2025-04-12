package inspect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"

	"github.com/omnistrate/ctl/internal/dataaccess"
)

var (
	outputFlag      string
	textMode        bool // Flag for text mode output
	kubeconfig      string
	kubeContext     string
	defaultKubeConf string
)

func init() {
	// Get default kubeconfig location from environment or standard path
	defaultKubeConf = os.Getenv("KUBECONFIG")
	if defaultKubeConf == "" {
		// If KUBECONFIG is not set, use the standard path
		if home := homedir.HomeDir(); home != "" {
			defaultKubeConf = filepath.Join(home, ".kube", "config")
		}
	}

	// Add flags for the command
	Cmd.Flags().StringVarP(&outputFlag, "output", "o", "table", "Output format (table|text|json)")
	Cmd.Flags().BoolVar(&textMode, "text", false, "Output text representation (shorthand for --output=text)")
	Cmd.Flags().StringVar(&kubeconfig, "kubeconfig", defaultKubeConf, "Path to the kubeconfig file")
	Cmd.Flags().StringVar(&kubeContext, "context", "", "Kubernetes context to use")
}

func runInspect(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	instanceID := args[0]

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		fmt.Printf("Warning: Kubeconfig file not found at %s\n", kubeconfig)
		fmt.Println("Trying to use in-cluster config or default settings...")
	}

	// If text mode flag is set, override output format
	if textMode {
		outputFlag = "text"
	}

	// Process based on output format
	switch outputFlag {
	case "text", "json":
		// For both text and JSON formats, we need to fetch the real data first
		inspectClient := dataaccess.NewK8sInspectClient(dataaccess.K8sClientConfig{
			Kubeconfig:  kubeconfig,
			KubeContext: kubeContext,
		})

		workloadItems, azItems, storageData, err := inspectClient.GetClusterData(ctx, instanceID)
		if err != nil {
			fmt.Printf("Error fetching cluster data: %v\n", err)
			return err
		}

		// For text format, generate a text representation of the data
		if outputFlag == "text" {
			fmt.Println(generateTextOutput(instanceID, workloadItems, azItems, storageData))
			return nil
		}

		// Convert data to JSON format
		type jsonOutput struct {
			InstanceID   string                               `json:"instanceId"`
			Workloads    []dataaccess.InspectWorkloadItem     `json:"workloads"`
			AZs          []dataaccess.InspectAZItem           `json:"availabilityZones"`
			StorageClass []dataaccess.InspectStorageClassItem `json:"storageClasses"`
		}

		output := jsonOutput{
			InstanceID:   instanceID,
			Workloads:    workloadItems,
			AZs:          azItems,
			StorageClass: storageData,
		}

		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Printf("Error converting to JSON: %v\n", err)
			return err
		}

		fmt.Println(string(jsonData))
		return nil
	case "table":
		// Default to interactive TUI mode
		// Continues to the code below
	default:
		return fmt.Errorf("unsupported output format: %s. Supported formats are table, text, and json", outputFlag)
	}

	// Create a K8sInspectClient
	inspectClient := dataaccess.NewK8sInspectClient(dataaccess.K8sClientConfig{
		Kubeconfig:  kubeconfig,
		KubeContext: kubeContext,
	})

	// Get data for inspection from Kubernetes cluster
	workloadItems, azItems, storageData, err := inspectClient.GetClusterData(ctx, instanceID)
	if err != nil {
		fmt.Printf("Warning: Error fetching cluster data: %v\nFalling back to sample data...\n", err)
		// Fall back to sample data if there's an error
		workloadItems, azItems, storageData = inspectClient.GetSampleData(instanceID)
	}

	// Validate data before launching TUI
	if len(azItems) == 0 {
		return fmt.Errorf("no data found for instance-id: %s", instanceID)
	}

	// Launch TUI with the data
	return launchTUI(instanceID, workloadItems, azItems, storageData)
}

// generateTextOutput creates a text representation of the cluster data
func generateTextOutput(instanceID string, workloadItems []dataaccess.InspectWorkloadItem, azItems []dataaccess.InspectAZItem, storageClasses []dataaccess.InspectStorageClassItem) string {
	var sb strings.Builder

	// Calculate summary stats
	totalVMs := make(map[string]bool) // Map to track unique VM names
	totalCPU := 0.0                   // Track total CPUs
	totalMemory := 0.0                // Track total memory in GB
	totalStorage := 0.0               // Track total storage in GiB
	storageClassCounts := make(map[string]float64)

	// First identify pods related to workloads
	workloadPods := make(map[string]bool)
	for _, workload := range workloadItems {
		for _, pods := range workload.AZs {
			for _, pod := range pods {
				workloadPods[pod.Name] = true
			}
		}
	}

	// Check if we have any workload pods
	hasAnyWorkloadPods := len(workloadPods) > 0

	// Count VMs and resources
	for _, az := range azItems {
		for _, vm := range az.VMs {
			if len(vm.Pods) == 0 {
				continue
			}

			if hasAnyWorkloadPods {
				hasWorkloadPod := false
				for _, pod := range vm.Pods {
					if workloadPods[pod.Name] {
						hasWorkloadPod = true
						break
					}
				}

				if !hasWorkloadPod {
					continue
				}
			}

			totalVMs[vm.Name] = true
			totalCPU += float64(vm.VCPUs)
			totalMemory += vm.MemoryGB
		}
	}

	// Identify PVCs attached to relevant pods
	relevantPVCs := make(map[string]bool)
	if hasAnyWorkloadPods {
		for _, workload := range workloadItems {
			for _, pods := range workload.AZs {
				for _, pod := range pods {
					for _, pvc := range pod.PVCs {
						relevantPVCs[pvc.Name] = true
					}
				}
			}
		}
	} else {
		for _, az := range azItems {
			for _, vm := range az.VMs {
				for _, pod := range vm.Pods {
					for _, pvc := range pod.PVCs {
						relevantPVCs[pvc.Name] = true
					}
				}
			}
		}
	}

	// Calculate storage totals
	for _, sc := range storageClasses {
		for _, pv := range sc.PVs {
			if !relevantPVCs[pv.PVCName] {
				continue
			}

			var size float64
			if strings.HasSuffix(pv.Size, "Gi") {
				_, err := fmt.Sscanf(pv.Size, "%f", &size)
				if err != nil {
					log.Fatal().Err(err).Msg("Error parsing PV size")
					return ""
				}
			} else if strings.HasSuffix(pv.Size, "Mi") {
				_, err := fmt.Sscanf(pv.Size, "%f", &size)
				if err != nil {
					log.Fatal().Err(err).Msg("Error parsing PV size")
					return ""
				}
				size = size / 1024.0
			}

			totalStorage += size
			storageClassCounts[sc.Name] += size
		}
	}

	// Title with enhanced styling
	sb.WriteString("╔══════════════════════════════════════════════════════════╗\n")
	sb.WriteString(fmt.Sprintf("║  Kubernetes Resource Inspector - Namespace: %-10s ║\n", instanceID))
	sb.WriteString("╚══════════════════════════════════════════════════════════╝\n\n")

	// Workload View summary
	sb.WriteString("📊 WORKLOAD VIEW SUMMARY\n")
	sb.WriteString("══════════════════════\n")
	sb.WriteString(fmt.Sprintf("VMs: %d, vCPUs: %.0f, Memory: %.1f GB RAM\n\n", len(totalVMs), totalCPU, totalMemory))

	// Workload details
	sb.WriteString("WORKLOAD DETAILS\n")
	sb.WriteString("═══════════════\n")
	if len(workloadItems) == 0 {
		sb.WriteString("No StatefulSets or Deployments found. Listing standalone pods:\n\n")

		// Group standalone pods by AZ
		podsByAZ := make(map[string][]dataaccess.InspectPodItem)
		for _, azItem := range azItems {
			for _, vm := range azItem.VMs {
				podsByAZ[azItem.Name] = append(podsByAZ[azItem.Name], vm.Pods...)
			}
		}

		// Show pods by AZ
		for az, pods := range podsByAZ {
			sb.WriteString(fmt.Sprintf("🌐 AZ: %s\n", az))
			for _, pod := range pods {
				sb.WriteString(fmt.Sprintf("  ⎈ Pod: %s (%s)\n", pod.Name, pod.Status))
				sb.WriteString(fmt.Sprintf("    Node: %s\n", pod.NodeName))

				// Show attached PVCs
				if len(pod.PVCs) > 0 {
					sb.WriteString("    Attached PVCs:\n")
					for _, pvc := range pod.PVCs {
						sb.WriteString(fmt.Sprintf("      💾 %s (%s, %s)\n", pvc.Name, pvc.Size, pvc.Status))
					}
				}
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
	} else {
		// Show workloads (StatefulSets and Deployments)
		for _, workload := range workloadItems {
			icon := "💾"
			if workload.Type == "Deployment" {
				icon = "🚀"
			}
			sb.WriteString(fmt.Sprintf("%s %s: %s\n", icon, workload.Type, workload.Name))

			// Show pods by AZ
			for az, pods := range workload.AZs {
				sb.WriteString(fmt.Sprintf("  🌐 AZ: %s\n", az))
				for _, pod := range pods {
					sb.WriteString(fmt.Sprintf("    ⎈ Pod: %s (%s)\n", pod.Name, pod.Status))
					sb.WriteString(fmt.Sprintf("      Node: %s\n", pod.NodeName))

					// Show attached PVCs
					if len(pod.PVCs) > 0 {
						sb.WriteString("      Attached PVCs:\n")
						for _, pvc := range pod.PVCs {
							sb.WriteString(fmt.Sprintf("        💾 %s (%s, %s)\n", pvc.Name, pvc.Size, pvc.Status))
						}
					}
				}
			}
			sb.WriteString("\n")
		}

		// Show standalone pods if any
		standalonePods := make(map[string][]dataaccess.InspectPodItem)
		for _, azItem := range azItems {
			for _, vm := range azItem.VMs {
				for _, pod := range vm.Pods {
					if workloadPods[pod.Name] {
						continue
					}
					standalonePods[azItem.Name] = append(standalonePods[azItem.Name], pod)
				}
			}
		}

		if len(standalonePods) > 0 {
			sb.WriteString("🔷 STANDALONE PODS\n")
			for az, pods := range standalonePods {
				sb.WriteString(fmt.Sprintf("  🌐 AZ: %s\n", az))
				for _, pod := range pods {
					sb.WriteString(fmt.Sprintf("    ⎈ Pod: %s (%s)\n", pod.Name, pod.Status))
					sb.WriteString(fmt.Sprintf("      Node: %s\n", pod.NodeName))

					// Show attached PVCs
					if len(pod.PVCs) > 0 {
						sb.WriteString("      Attached PVCs:\n")
						for _, pvc := range pod.PVCs {
							sb.WriteString(fmt.Sprintf("        💾 %s (%s, %s)\n", pvc.Name, pvc.Size, pvc.Status))
						}
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	// Infrastructure View
	sb.WriteString("\n🏢 INFRASTRUCTURE VIEW SUMMARY\n")
	sb.WriteString("═══════════════════════════\n")
	sb.WriteString(fmt.Sprintf("VMs: %d, vCPUs: %.0f, Memory: %.1f GB RAM\n", len(totalVMs), totalCPU, totalMemory))

	// Show instance types
	instanceTypes := make(map[string]int)
	for _, az := range azItems {
		for _, vm := range az.VMs {
			if len(vm.Pods) == 0 {
				continue
			}

			if hasAnyWorkloadPods {
				hasWorkloadPod := false
				for _, pod := range vm.Pods {
					if workloadPods[pod.Name] {
						hasWorkloadPod = true
						break
					}
				}

				if !hasWorkloadPod {
					continue
				}
			}

			instanceTypes[vm.InstanceType]++
		}
	}

	if len(instanceTypes) > 0 {
		sb.WriteString("Instance Types:\n")
		for instanceType, count := range instanceTypes {
			sb.WriteString(fmt.Sprintf("  %s: %d\n", instanceType, count))
		}
	}
	sb.WriteString("\n")

	// VM details
	sb.WriteString("INFRASTRUCTURE DETAILS\n")
	sb.WriteString("═════════════════════\n")
	for _, az := range azItems {
		sb.WriteString(fmt.Sprintf("🌐 AZ: %s\n", az.Name))

		for _, vm := range az.VMs {
			hasRelevantPods := false
			if hasAnyWorkloadPods {
				for _, pod := range vm.Pods {
					if workloadPods[pod.Name] {
						hasRelevantPods = true
						break
					}
				}
			} else if len(vm.Pods) > 0 {
				hasRelevantPods = true
			}

			if !hasRelevantPods {
				continue
			}

			sb.WriteString(fmt.Sprintf("  💻 VM: %s\n", vm.Name))
			sb.WriteString(fmt.Sprintf("    Type: %s, vCPUs: %d, Memory: %.1f GB\n", vm.InstanceType, vm.VCPUs, vm.MemoryGB))
			sb.WriteString("    Pods:\n")

			for _, pod := range vm.Pods {
				if hasAnyWorkloadPods && !workloadPods[pod.Name] {
					continue
				}
				sb.WriteString(fmt.Sprintf("      ⎈ %s (%s)\n", pod.Name, pod.Status))
			}
		}
		sb.WriteString("\n")
	}

	// Storage View
	sb.WriteString("\n💿 STORAGE VIEW SUMMARY\n")
	sb.WriteString("═════════════════════\n")
	sb.WriteString(fmt.Sprintf("Total Storage: %.1f GiB\n", totalStorage))

	if len(storageClassCounts) > 0 {
		sb.WriteString("Storage Classes:\n")
		for sc, size := range storageClassCounts {
			sb.WriteString(fmt.Sprintf("  %s: %.1f GiB\n", sc, size))
		}
	}
	sb.WriteString("\n")

	// Storage details
	sb.WriteString("STORAGE DETAILS\n")
	sb.WriteString("═══════════════\n")

	// Find all pods with PVCs
	podsWithStorage := make(map[string][]dataaccess.InspectPodItem)
	for _, azItem := range azItems {
		for _, vm := range azItem.VMs {
			for _, pod := range vm.Pods {
				if len(pod.PVCs) > 0 {
					if hasAnyWorkloadPods && !workloadPods[pod.Name] {
						continue
					}

					// Find workload for this pod (if any)
					workloadKey := "Standalone"
					for _, workload := range workloadItems {
						for _, pods := range workload.AZs {
							for _, wPod := range pods {
								if wPod.Name == pod.Name {
									workloadKey = fmt.Sprintf("%s: %s", workload.Type, workload.Name)
									break
								}
							}
						}
					}
					podsWithStorage[workloadKey] = append(podsWithStorage[workloadKey], pod)
				}
			}
		}
	}

	if len(podsWithStorage) == 0 {
		sb.WriteString("No persistent volumes found.\n")
	} else {
		for workload, pods := range podsWithStorage {
			if workload == "Standalone" {
				sb.WriteString("🔷 Standalone Pods:\n")
			} else {
				if strings.HasPrefix(workload, "StatefulSet") {
					sb.WriteString(fmt.Sprintf("💾 %s\n", workload))
				} else {
					sb.WriteString(fmt.Sprintf("🚀 %s\n", workload))
				}
			}

			for _, pod := range pods {
				sb.WriteString(fmt.Sprintf("  ⎈ Pod: %s (%s)\n", pod.Name, pod.Status))
				sb.WriteString("    PVCs:\n")

				for _, pvc := range pod.PVCs {
					sb.WriteString(fmt.Sprintf("      💾 %s\n", pvc.Name))
					sb.WriteString(fmt.Sprintf("        Size: %s, Status: %s\n", pvc.Size, pvc.Status))
					if pvc.StorageClass != "" {
						sb.WriteString(fmt.Sprintf("        Storage Class: %s\n", pvc.StorageClass))
					}
					if pvc.PVName != "" {
						sb.WriteString(fmt.Sprintf("        Bound to PV: %s\n", pvc.PVName))
					}
					if len(pvc.AccessModes) > 0 {
						sb.WriteString(fmt.Sprintf("        Access Modes: %s\n", strings.Join(pvc.AccessModes, ", ")))
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// launchTUI creates and runs the terminal UI
func launchTUI(instanceID string, workloadItems []dataaccess.InspectWorkloadItem, azItems []dataaccess.InspectAZItem, storageClasses []dataaccess.InspectStorageClassItem) error {
	// Setup TUI application
	app := tview.NewApplication()

	// Create the layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Create tabs for different views
	tabs := tview.NewPages()

	// Calculate view summary stats that show resources related to all relevant pods
	// First identify all pods related to workloads
	workloadPods := make(map[string]bool)
	for _, workload := range workloadItems {
		for _, pods := range workload.AZs {
			for _, pod := range pods {
				workloadPods[pod.Name] = true
			}
		}
	}

	// Map to track unique VM names that host relevant pods (workload or standalone)
	totalVMs := make(map[string]bool)
	totalCPU := 0.0    // Track total CPUs
	totalMemory := 0.0 // Track total memory in GB

	// Count unique VMs and get total resources
	// If there are workload pods, only count VMs hosting those
	// If there are no workload pods, count all VMs with any pods
	hasAnyWorkloadPods := len(workloadPods) > 0

	for _, az := range azItems {
		for _, vm := range az.VMs {
			if len(vm.Pods) == 0 {
				continue // Skip VMs with no pods at all
			}

			if hasAnyWorkloadPods {
				// Only count VMs that host workload pods
				hasWorkloadPod := false
				for _, pod := range vm.Pods {
					if workloadPods[pod.Name] {
						hasWorkloadPod = true
						break
					}
				}

				if !hasWorkloadPod {
					continue // Skip VMs with no workload pods
				}
			}

			// Count this VM
			totalVMs[vm.Name] = true
			totalCPU += float64(vm.VCPUs)
			totalMemory += vm.MemoryGB
		}
	}

	// Create workload tree with enhanced styling and summary header
	workloadTree := tview.NewTreeView()

	// Create summary node with instance count and resource info
	workloadSummary := fmt.Sprintf("📊 Workload View - %d VMs, %.0f vCPUs, %.1f GB RAM",
		len(totalVMs), totalCPU, totalMemory)

	workloadRoot := tview.NewTreeNode(workloadSummary).
		SetColor(tcell.ColorYellow).
		SetSelectable(true)
	workloadTree.SetRoot(workloadRoot)
	workloadTree.SetCurrentNode(workloadRoot)
	workloadTree.SetBorder(true).SetTitle(" Workloads ").SetTitleColor(tcell.ColorYellow)
	workloadTree.SetGraphics(true)

	// Add workload data to tree with enhanced colors
	if len(workloadItems) == 0 {
		// If no workloads found, show pods directly
		podsByAZ := make(map[string][]dataaccess.InspectPodItem)

		// Get all pods organized by AZ
		for _, azItem := range azItems {
			azName := azItem.Name
			for _, vm := range azItem.VMs {
				podsByAZ[azName] = append(podsByAZ[azName], vm.Pods...)
			}
		}

		// Show pods directly under workload view
		for az, pods := range podsByAZ {
			azNode := tview.NewTreeNode(fmt.Sprintf("🌐 AZ: %s", az)).
				SetColor(tcell.ColorBlue).
				SetReference(az).
				SetSelectable(true)

			for _, pod := range pods {
				// Color pods based on status
				var podColor tcell.Color

				switch pod.Status {
				case "Running":
					podColor = tcell.ColorGreen
				case "Pending":
					podColor = tcell.ColorYellow
				case "Failed":
					podColor = tcell.ColorRed
				default:
					podColor = tcell.ColorGray
				}

				// Use Kubernetes logo for pods
				podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
					SetColor(podColor).
					SetReference(pod).
					SetSelectable(true)
				azNode.AddChild(podNode)
			}

			if len(azNode.GetChildren()) > 0 {
				workloadRoot.AddChild(azNode)
			}
		}
	} else {
		// Normal case - show workloads
		for _, workload := range workloadItems {
			// Choose color based on workload type
			var workloadColor tcell.Color
			var workloadIcon string

			if workload.Type == "StatefulSet" {
				workloadColor = tcell.ColorGreen
				workloadIcon = "💾"
			} else {
				workloadColor = tcell.ColorDarkCyan
				workloadIcon = "🚀"
			}

			workloadNode := tview.NewTreeNode(fmt.Sprintf("%s %s: %s", workloadIcon, workload.Type, workload.Name)).
				SetColor(workloadColor).
				SetReference(workload).
				SetSelectable(true)

			for az, pods := range workload.AZs {
				azNode := tview.NewTreeNode(fmt.Sprintf("🌐 AZ: %s", az)).
					SetColor(tcell.ColorBlue).
					SetReference(az).
					SetSelectable(true)

				for _, pod := range pods {
					// Color pods based on status
					var podColor tcell.Color

					switch pod.Status {
					case "Running":
						podColor = tcell.ColorGreen
					case "Pending":
						podColor = tcell.ColorYellow
					case "Failed":
						podColor = tcell.ColorRed
					default:
						podColor = tcell.ColorGray
					}

					// Use Kubernetes logo for pods
					podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
						SetColor(podColor).
						SetReference(pod).
						SetSelectable(true)
					azNode.AddChild(podNode)
				}

				workloadNode.AddChild(azNode)
			}

			workloadRoot.AddChild(workloadNode)
		}

		// Now add standalone pods (pods not part of any deployment or statefulset)
		standalonePods := tview.NewTreeNode("🔷 Standalone Pods").
			SetColor(tcell.ColorBlue).
			SetSelectable(true)

		// Find standalone pods by checking all pods in azItems
		hasStandalonePods := false

		// Create a map of all pods that are part of workloads for quick lookup
		workloadPodNames := make(map[string]bool)
		for _, workload := range workloadItems {
			for _, pods := range workload.AZs {
				for _, pod := range pods {
					workloadPodNames[pod.Name] = true
				}
			}
		}

		// Create map to organize standalone pods by AZ
		standalonePodsByAZ := make(map[string][]dataaccess.InspectPodItem)

		// Find pods that aren't part of any workload
		for _, azItem := range azItems {
			azName := azItem.Name
			for _, vm := range azItem.VMs {
				for _, pod := range vm.Pods {
					// Skip pods that are part of workloads
					if workloadPodNames[pod.Name] {
						continue
					}

					// Add to standalone pods map
					standalonePodsByAZ[azName] = append(standalonePodsByAZ[azName], pod)
					hasStandalonePods = true
				}
			}
		}

		// Add standalone pods to the tree by AZ
		for az, pods := range standalonePodsByAZ {
			azNode := tview.NewTreeNode(fmt.Sprintf("🌐 AZ: %s", az)).
				SetColor(tcell.ColorBlue).
				SetReference(az).
				SetSelectable(true)

			for _, pod := range pods {
				// Color pods based on status
				var podColor tcell.Color

				switch pod.Status {
				case "Running":
					podColor = tcell.ColorGreen
				case "Pending":
					podColor = tcell.ColorYellow
				case "Failed":
					podColor = tcell.ColorRed
				default:
					podColor = tcell.ColorGray
				}

				// Use Kubernetes logo for pods
				podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
					SetColor(podColor).
					SetReference(pod).
					SetSelectable(true)
				azNode.AddChild(podNode)
			}

			standalonePods.AddChild(azNode)
		}

		// Only add standalone pods section if there are any
		if hasStandalonePods {
			workloadRoot.AddChild(standalonePods)
		}
	}

	// Create infrastructure tree with enhanced styling
	infraTree := tview.NewTreeView()

	// Count instance types for VMs with relevant pods
	instanceTypes := make(map[string]int)
	for _, az := range azItems {
		for _, vm := range az.VMs {
			if len(vm.Pods) == 0 {
				continue // Skip VMs with no pods at all
			}

			if hasAnyWorkloadPods {
				// Only count VMs that host workload pods when workloads exist
				hasWorkloadPod := false
				for _, pod := range vm.Pods {
					if workloadPods[pod.Name] {
						hasWorkloadPod = true
						break
					}
				}

				if !hasWorkloadPod {
					continue // Skip VMs with no workload pods
				}
			}

			// Count this VM's instance type
			instanceTypes[vm.InstanceType]++
		}
	}

	// Create instance type summary string
	var instanceTypeSummary strings.Builder
	instanceTypeSummary.WriteString(fmt.Sprintf("🏢 Infrastructure View - %d VMs, %.0f vCPUs, %.1f GB RAM\n",
		len(totalVMs), totalCPU, totalMemory))

	instanceTypeSummary.WriteString("Instance Types: ")
	i := 0
	for instanceType, count := range instanceTypes {
		if i > 0 {
			instanceTypeSummary.WriteString(", ")
		}
		instanceTypeSummary.WriteString(fmt.Sprintf("%s (%d)", instanceType, count))
		i++
	}

	infraRoot := tview.NewTreeNode(instanceTypeSummary.String()).
		SetColor(tcell.ColorYellow).
		SetSelectable(true)
	infraTree.SetRoot(infraRoot)
	infraTree.SetCurrentNode(infraRoot)
	infraTree.SetBorder(true).SetTitle(" Infrastructure ").SetTitleColor(tcell.ColorYellow)
	infraTree.SetGraphics(true)

	// Add AZ data to tree with enhanced colors
	for _, az := range azItems {
		azNode := tview.NewTreeNode(fmt.Sprintf("🌐 AZ: %s", az.Name)).
			SetColor(tcell.ColorGreen).
			SetReference(az).
			SetSelectable(true)

		// Track if we've added any VMs to this AZ
		hasVMs := false

		for _, vm := range az.VMs {
			// Skip VMs with no pods in this namespace
			if len(vm.Pods) == 0 {
				continue
			}

			// Different colors based on instance type
			var vmColor tcell.Color

			if strings.Contains(vm.InstanceType, "xlarge") {
				vmColor = tcell.ColorDarkCyan
			} else {
				vmColor = tcell.ColorBlue
			}

			vmNode := tview.NewTreeNode(fmt.Sprintf("💻 VM: %s (Type: %s, vCPUs: %d, Memory: %.1fGB)",
				vm.Name, vm.InstanceType, vm.VCPUs, vm.MemoryGB)).
				SetColor(vmColor).
				SetReference(vm).
				SetSelectable(true)

			for _, pod := range vm.Pods {
				// Color pods based on status
				var podColor tcell.Color

				switch pod.Status {
				case "Running":
					podColor = tcell.ColorGreen
				case "Pending":
					podColor = tcell.ColorYellow
				case "Failed":
					podColor = tcell.ColorRed
				default:
					podColor = tcell.ColorGray
				}

				// Use Kubernetes logo for pods
				podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
					SetColor(podColor).
					SetReference(pod).
					SetSelectable(true)
				vmNode.AddChild(podNode)
			}

			azNode.AddChild(vmNode)
			hasVMs = true
		}

		// Only add AZs that have VMs with pods from this namespace
		if hasVMs {
			infraRoot.AddChild(azNode)
		}
	}

	// Create storage tree with enhanced styling
	storageTree := tview.NewTreeView()

	// Calculate storage view summary stats
	totalStorage := 0.0 // in GiB
	storageClassMap := make(map[string]float64)

	// Build a map of all relevant PVCs
	relevantPVCs := make(map[string]bool)

	// First include PVCs from workload pods if there are any
	if hasAnyWorkloadPods {
		// PVCs attached to workload pods
		for _, workload := range workloadItems {
			for _, pods := range workload.AZs {
				for _, pod := range pods {
					for _, pvc := range pod.PVCs {
						relevantPVCs[pvc.Name] = true
					}
				}
			}
		}
	} else {
		// If no workloads, include all PVCs from all pods
		for _, az := range azItems {
			for _, vm := range az.VMs {
				for _, pod := range vm.Pods {
					for _, pvc := range pod.PVCs {
						relevantPVCs[pvc.Name] = true
					}
				}
			}
		}
	}

	// Process storage classes to get total storage and per-class breakdown
	for _, sc := range storageClasses {
		for _, pv := range sc.PVs {
			// Skip PVs not bound to relevant PVCs
			if !relevantPVCs[pv.PVCName] {
				continue
			}

			// Parse size (simple version)
			sizeStr := pv.Size
			var size float64

			// Very basic parsing of common formats
			if strings.HasSuffix(sizeStr, "Gi") {
				_, err := fmt.Sscanf(sizeStr, "%f", &size)
				if err != nil {
					return err
				}
			} else if strings.HasSuffix(sizeStr, "Mi") {
				_, err := fmt.Sscanf(sizeStr, "%f", &size)
				if err != nil {
					return err
				}
				size = size / 1024.0
			}

			totalStorage += size
			storageClassMap[sc.Name] += size
		}
	}

	// Create storage summary
	var storageSummary strings.Builder
	storageSummary.WriteString(fmt.Sprintf("💿 Storage View - %.1f GiB Total Storage\n", totalStorage))

	// Add storage class breakdown
	storageSummary.WriteString("Storage Classes: ")
	j := 0
	for sc, size := range storageClassMap {
		if j > 0 {
			storageSummary.WriteString(", ")
		}
		storageSummary.WriteString(fmt.Sprintf("%s (%.1f GiB)", sc, size))
		j++
	}

	storageRoot := tview.NewTreeNode(storageSummary.String()).
		SetColor(tcell.ColorYellow).
		SetSelectable(true)
	storageTree.SetRoot(storageRoot)
	storageTree.SetCurrentNode(storageRoot)
	storageTree.SetBorder(true).SetTitle(" Storage ").SetTitleColor(tcell.ColorYellow)
	storageTree.SetGraphics(true)

	// Create a map to store Storage Class details for pop-ups
	storageClassDetails := make(map[string]dataaccess.InspectStorageClassItem)
	for _, sc := range storageClasses {
		storageClassDetails[sc.Name] = sc
	}

	// Create a map to store PV details for pop-ups
	pvDetails := make(map[string]dataaccess.InspectPVItem)
	for _, sc := range storageClasses {
		for _, pv := range sc.PVs {
			pvDetails[pv.Name] = pv
		}
	}

	// Create organization by StatefulSets and Deployments
	statefulSetNode := tview.NewTreeNode("💾 StatefulSets").
		SetColor(tcell.ColorGreen).
		SetSelectable(true)
	storageRoot.AddChild(statefulSetNode)

	deploymentNode := tview.NewTreeNode("🚀 Deployments").
		SetColor(tcell.ColorDarkCyan).
		SetSelectable(true)
	storageRoot.AddChild(deploymentNode)

	// Create modal for PV details
	pvModal := tview.NewModal()
	pvModal.SetBackgroundColor(tcell.ColorBlack)
	pvModal.SetTextColor(tcell.ColorWhite)
	// Modal doesn't support text alignment
	pvModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Close" {
			// Return to the source page stored in the title
			sourcePage := pvModal.GetTitle()
			if sourcePage == "" {
				// Default to storage view
				sourcePage = "storage"
			}
			tabs.SwitchToPage(sourcePage)
		}
	})

	// Create modal for Pod details
	podModal := tview.NewModal()
	podModal.SetBackgroundColor(tcell.ColorBlack)
	podModal.SetTextColor(tcell.ColorWhite)
	// Modal doesn't support text alignment
	podModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Close" {
			// Get source page from the modal title
			sourcePage := podModal.GetTitle()
			if sourcePage == "" {
				// Default to workload view if we can't determine the source page
				sourcePage = "workload"
			}
			tabs.SwitchToPage(sourcePage)
		}
	})

	// Add pages for modals
	tabs.AddPage("pvModal", pvModal, false, false)
	tabs.AddPage("podModal", podModal, false, false)

	// Function to show PV details modal
	showPVDetails := func(pvName string) {
		pv, exists := pvDetails[pvName]
		if !exists {
			return
		}

		// Save current page
		sourcePage, _ := tabs.GetFrontPage()

		// Define standard indentation
		const indent = "  "

		sc, scExists := storageClassDetails[pv.StorageClass]

		var details strings.Builder
		details.WriteString(fmt.Sprintf("PV Name: %s\n", pv.Name))
		details.WriteString(fmt.Sprintf("Size: %s\n", pv.Size))
		details.WriteString(fmt.Sprintf("Status: %s\n", pv.Status))
		details.WriteString(fmt.Sprintf("Access Modes: %s\n", strings.Join(pv.AccessModes, ", ")))
		details.WriteString(fmt.Sprintf("Volume Type: %s\n", pv.VolumeType))

		if pv.PVCName != "" {
			details.WriteString(fmt.Sprintf("Bound to PVC: %s\n", pv.PVCName))
			details.WriteString(fmt.Sprintf("PVC Namespace: %s\n", pv.PVCNamespace))
		}

		details.WriteString("\nStorage Class Details:\n")
		if scExists {
			details.WriteString(fmt.Sprintf("%sName: %s\n", indent, sc.Name))
			details.WriteString(fmt.Sprintf("%sProvisioner: %s\n", indent, sc.Provisioner))

			if len(sc.Parameters) > 0 {
				details.WriteString("\nParameters:\n")

				// Sort keys for consistent display
				keys := make([]string, 0, len(sc.Parameters))
				for k := range sc.Parameters {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				for _, k := range keys {
					details.WriteString(fmt.Sprintf("%s%s: %s\n", indent, k, sc.Parameters[k]))
				}
			}
		} else {
			details.WriteString(fmt.Sprintf("%sName: %s\n", indent, pv.StorageClass))
			details.WriteString(fmt.Sprintf("%sDetails not available\n", indent))
		}

		pvModal.SetText(details.String())
		// Store the source page in the title
		pvModal.SetTitle(sourcePage)
		pvModal.ClearButtons().AddButtons([]string{"Close"})
		tabs.SwitchToPage("pvModal")
	}

	// Function to show Pod details modal
	showPodDetails := func(pod dataaccess.InspectPodItem) {
		// Remember which page we're on by storing it in the modal title
		sourcePage, _ := tabs.GetFrontPage()

		// Define standard indentation
		const indent = "  "
		const subIndent = "    "

		var details strings.Builder
		details.WriteString(fmt.Sprintf("Pod Name: %s\n", pod.Name))
		details.WriteString(fmt.Sprintf("Status: %s\n", pod.Status))
		details.WriteString(fmt.Sprintf("Node: %s\n", pod.NodeName))
		details.WriteString(fmt.Sprintf("Namespace: %s\n", pod.Namespace))

		// Add labels with proper indentation
		if len(pod.Labels) > 0 {
			details.WriteString("\nLabels:\n")
			keys := make([]string, 0, len(pod.Labels))
			for k := range pod.Labels {
				keys = append(keys, k)
			}
			sort.Strings(keys) // Sort keys for consistent display

			for _, k := range keys {
				details.WriteString(fmt.Sprintf("%s%s: %s\n", indent, k, pod.Labels[k]))
			}
		}

		// Add resource limits and requests with proper indentation
		details.WriteString("\nResource Limits:\n")
		if len(pod.Resources.Limits) > 0 {
			keys := make([]string, 0, len(pod.Resources.Limits))
			for k := range pod.Resources.Limits {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				details.WriteString(fmt.Sprintf("%s%s: %s\n", indent, k, pod.Resources.Limits[k]))
			}
		} else {
			details.WriteString(fmt.Sprintf("%sNone\n", indent))
		}

		details.WriteString("\nResource Requests:\n")
		if len(pod.Resources.Requests) > 0 {
			keys := make([]string, 0, len(pod.Resources.Requests))
			for k := range pod.Resources.Requests {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				details.WriteString(fmt.Sprintf("%s%s: %s\n", indent, k, pod.Resources.Requests[k]))
			}
		} else {
			details.WriteString(fmt.Sprintf("%sNone\n", indent))
		}

		// Add attached PVCs with proper indentation
		if len(pod.PVCs) > 0 {
			details.WriteString("\nAttached PVCs:\n")
			for i, pvc := range pod.PVCs {
				details.WriteString(fmt.Sprintf("%s%s\n", indent, pvc.Name))
				details.WriteString(fmt.Sprintf("%sSize: %s\n", subIndent, pvc.Size))
				details.WriteString(fmt.Sprintf("%sStatus: %s\n", subIndent, pvc.Status))
				if pvc.StorageClass != "" {
					details.WriteString(fmt.Sprintf("%sStorage Class: %s\n", subIndent, pvc.StorageClass))
				}
				if len(pvc.AccessModes) > 0 {
					details.WriteString(fmt.Sprintf("%sAccess Modes: %s\n", subIndent, strings.Join(pvc.AccessModes, ", ")))
				}
				if i < len(pod.PVCs)-1 {
					details.WriteString("\n")
				}
			}
		}

		podModal.SetText(details.String())
		// Store the source page in the title for returning later
		podModal.SetTitle(sourcePage)
		podModal.ClearButtons().AddButtons([]string{"Close"})
		tabs.SwitchToPage("podModal")
	}

	// Process workloads and organize by type
	// Add a fallback node if no PVCs are found
	noStorageNode := tview.NewTreeNode("No persistent volumes found").
		SetColor(tcell.ColorOrange).
		SetSelectable(true)

	// Track if we found any PVCs
	foundPVCs := false

	if len(workloadItems) == 0 {
		// Show pods directly in storage view if no workloads found
		podsNode := tview.NewTreeNode("🔷 Pods").
			SetColor(tcell.ColorBlue). // Use ColorBlue instead of ColorCyan
			SetSelectable(true)
		storageRoot.AddChild(podsNode)

		// Get all pods with PVCs
		for _, azItem := range azItems {
			for _, vm := range azItem.VMs {
				for _, pod := range vm.Pods {
					// Skip pods without PVCs
					if len(pod.PVCs) == 0 {
						continue
					}

					foundPVCs = true

					// Color pods based on status
					var podColor tcell.Color

					switch pod.Status {
					case "Running":
						podColor = tcell.ColorGreen
					case "Pending":
						podColor = tcell.ColorYellow
					case "Failed":
						podColor = tcell.ColorRed
					default:
						podColor = tcell.ColorGray
					}

					podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
						SetColor(podColor).
						SetReference(pod).
						SetSelectable(true)

					// Add PVCs for this pod
					for _, pvc := range pod.PVCs {
						var pvcColor tcell.Color

						switch pvc.Status {
						case "Bound":
							pvcColor = tcell.ColorGreen
						case "Pending":
							pvcColor = tcell.ColorYellow
						case "Lost":
							pvcColor = tcell.ColorRed
						default:
							pvcColor = tcell.ColorGray
						}

						pvcNode := tview.NewTreeNode(fmt.Sprintf("💾 PVC: %s (Size: %s, Status: %s)",
							pvc.Name, pvc.Size, pvc.Status)).
							SetColor(pvcColor).
							SetReference(pvc).
							SetSelectable(true)

						// Add PVC details
						if pvc.StorageClass != "" {
							scNode := tview.NewTreeNode(fmt.Sprintf("StorageClass: %s, Access: %s",
								pvc.StorageClass, strings.Join(pvc.AccessModes, ", "))).
								SetColor(tcell.ColorDarkGray).
								SetSelectable(true)
							pvcNode.AddChild(scNode)
						}

						// Add PV reference if bound
						if pvc.PVName != "" {
							pvRefNode := tview.NewTreeNode(fmt.Sprintf("📁 PV: %s", pvc.PVName)).
								SetColor(tcell.ColorDarkCyan).
								SetReference(pvc.PVName). // Store PV name as reference
								SetSelectable(true)
							pvcNode.AddChild(pvRefNode)
						}

						podNode.AddChild(pvcNode)
					}

					podsNode.AddChild(podNode)
				}
			}
		}

		// If no pods with PVCs were found, show a message
		if !foundPVCs {
			podsNode.AddChild(tview.NewTreeNode("No pods with persistent volumes").
				SetColor(tcell.ColorYellow))
		}
	} else {
		// Normal case - organize by workload types
		for _, workload := range workloadItems {
			// Choose parent node based on workload type
			var parentNode *tview.TreeNode
			var workloadColor tcell.Color
			var workloadIcon string

			if workload.Type == "StatefulSet" {
				parentNode = statefulSetNode
				workloadColor = tcell.ColorGreen
				workloadIcon = "💾"
			} else {
				parentNode = deploymentNode
				workloadColor = tcell.ColorDarkCyan
				workloadIcon = "🚀"
			}

			// Create workload node
			workloadNode := tview.NewTreeNode(fmt.Sprintf("%s %s: %s", workloadIcon, workload.Type, workload.Name)).
				SetColor(workloadColor).
				SetReference(workload).
				SetSelectable(true)

			// Add pods for this workload
			for _, pods := range workload.AZs {
				for _, pod := range pods {
					// Skip pods without PVCs
					if len(pod.PVCs) == 0 {
						continue
					}

					foundPVCs = true

					// Color pods based on status
					var podColor tcell.Color

					switch pod.Status {
					case "Running":
						podColor = tcell.ColorGreen
					case "Pending":
						podColor = tcell.ColorYellow
					case "Failed":
						podColor = tcell.ColorRed
					default:
						podColor = tcell.ColorGray
					}

					podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
						SetColor(podColor).
						SetReference(pod).
						SetSelectable(true)

					// Add PVCs for this pod
					for _, pvc := range pod.PVCs {
						var pvcColor tcell.Color

						switch pvc.Status {
						case "Bound":
							pvcColor = tcell.ColorGreen
						case "Pending":
							pvcColor = tcell.ColorYellow
						case "Lost":
							pvcColor = tcell.ColorRed
						default:
							pvcColor = tcell.ColorGray
						}

						pvcNode := tview.NewTreeNode(fmt.Sprintf("💾 PVC: %s (Size: %s, Status: %s)",
							pvc.Name, pvc.Size, pvc.Status)).
							SetColor(pvcColor).
							SetReference(pvc).
							SetSelectable(true)

						// Add PVC details
						if pvc.StorageClass != "" {
							scNode := tview.NewTreeNode(fmt.Sprintf("StorageClass: %s, Access: %s",
								pvc.StorageClass, strings.Join(pvc.AccessModes, ", "))).
								SetColor(tcell.ColorDarkGray).
								SetSelectable(true)
							pvcNode.AddChild(scNode)
						}

						// Add PV reference if bound
						if pvc.PVName != "" {
							pvRefNode := tview.NewTreeNode(fmt.Sprintf("📁 PV: %s", pvc.PVName)).
								SetColor(tcell.ColorDarkCyan).
								SetReference(pvc.PVName). // Store PV name as reference
								SetSelectable(true)
							pvcNode.AddChild(pvRefNode)
						}

						podNode.AddChild(pvcNode)
					}

					workloadNode.AddChild(podNode)
				}
			}

			// Only add workload nodes if they have pods with PVCs
			if len(workloadNode.GetChildren()) > 0 {
				parentNode.AddChild(workloadNode)
			}
		}

		// Add standalone pods with PVCs to storage view
		standalonePods := tview.NewTreeNode("🔷 Standalone Pods").
			SetColor(tcell.ColorBlue).
			SetSelectable(true)

		// Find standalone pods with PVCs
		hasStandalonePVCs := false

		// Create a map of all pods that are part of workloads for quick lookup
		workloadPodNames := make(map[string]bool)
		for _, workload := range workloadItems {
			for _, pods := range workload.AZs {
				for _, pod := range pods {
					workloadPodNames[pod.Name] = true
				}
			}
		}

		// Find all pods with PVCs that aren't part of workloads
		for _, azItem := range azItems {
			for _, vm := range azItem.VMs {
				for _, pod := range vm.Pods {
					// Skip pods that are part of workloads or have no PVCs
					if workloadPodNames[pod.Name] || len(pod.PVCs) == 0 {
						continue
					}

					hasStandalonePVCs = true
					foundPVCs = true

					// Color pods based on status
					var podColor tcell.Color

					switch pod.Status {
					case "Running":
						podColor = tcell.ColorGreen
					case "Pending":
						podColor = tcell.ColorYellow
					case "Failed":
						podColor = tcell.ColorRed
					default:
						podColor = tcell.ColorGray
					}

					podNode := tview.NewTreeNode(fmt.Sprintf("⎈ Pod: %s (%s)", pod.Name, pod.Status)).
						SetColor(podColor).
						SetReference(pod).
						SetSelectable(true)

					// Add PVCs for this pod
					for _, pvc := range pod.PVCs {
						var pvcColor tcell.Color

						switch pvc.Status {
						case "Bound":
							pvcColor = tcell.ColorGreen
						case "Pending":
							pvcColor = tcell.ColorYellow
						case "Lost":
							pvcColor = tcell.ColorRed
						default:
							pvcColor = tcell.ColorGray
						}

						pvcNode := tview.NewTreeNode(fmt.Sprintf("💾 PVC: %s (Size: %s, Status: %s)",
							pvc.Name, pvc.Size, pvc.Status)).
							SetColor(pvcColor).
							SetReference(pvc).
							SetSelectable(true)

						// Add PVC details
						if pvc.StorageClass != "" {
							scNode := tview.NewTreeNode(fmt.Sprintf("StorageClass: %s, Access: %s",
								pvc.StorageClass, strings.Join(pvc.AccessModes, ", "))).
								SetColor(tcell.ColorDarkGray).
								SetSelectable(true)
							pvcNode.AddChild(scNode)
						}

						// Add PV reference if bound
						if pvc.PVName != "" {
							pvRefNode := tview.NewTreeNode(fmt.Sprintf("📁 PV: %s", pvc.PVName)).
								SetColor(tcell.ColorDarkCyan).
								SetReference(pvc.PVName). // Store PV name as reference
								SetSelectable(true)
							pvcNode.AddChild(pvRefNode)
						}

						podNode.AddChild(pvcNode)
					}

					standalonePods.AddChild(podNode)
				}
			}
		}

		// Only add standalone pods section if there are any with PVCs
		if hasStandalonePVCs {
			storageRoot.AddChild(standalonePods)
		}

		// If no deployments with PVCs were found, show a message
		if len(deploymentNode.GetChildren()) == 0 {
			deploymentNode.AddChild(tview.NewTreeNode("No deployments with persistent volumes").
				SetColor(tcell.ColorYellow))
		}

		// If no statefulsets with PVCs were found, show a message
		if len(statefulSetNode.GetChildren()) == 0 {
			statefulSetNode.AddChild(tview.NewTreeNode("No statefulsets with persistent volumes").
				SetColor(tcell.ColorYellow))
		}
	}

	// If no PVCs were found at all, show a message
	if !foundPVCs {
		storageRoot.AddChild(noStorageNode)
	}

	// Set up special event handler for storage tree
	storageTree.SetSelectedFunc(func(node *tview.TreeNode) {
		// Check if node is a PV reference
		ref := node.GetReference()
		if pvName, ok := ref.(string); ok {
			// If it's a PV node, show modal with details
			showPVDetails(pvName)
			return
		}

		// Check if it's a pod
		if pod, ok := ref.(dataaccess.InspectPodItem); ok {
			// If it's a pod, show the pod details modal
			showPodDetails(pod)
			return
		}

		// Regular expand/collapse for other nodes
		expanded := node.IsExpanded()
		node.SetExpanded(!expanded)

		// Change color briefly to indicate action
		originalColor := node.GetColor()
		node.SetColor(tcell.ColorLightCyan)

		// Schedule reset of color
		go func() {
			time.Sleep(150 * time.Millisecond)
			app.QueueUpdateDraw(func() {
				node.SetColor(originalColor)
			})
		}()
	})

	// Create event handler function for tree expand/collapse and pod details with visual feedback
	treeExpandCollapseAndDetailsFunc := func(node *tview.TreeNode) {
		// Check reference type to see if it's a pod
		ref := node.GetReference()
		if pod, ok := ref.(dataaccess.InspectPodItem); ok {
			// If it's a pod, show the pod details modal
			showPodDetails(pod)
			return
		}

		// Default behavior - toggle expand/collapse with visual feedback
		expanded := node.IsExpanded()
		node.SetExpanded(!expanded)

		// Change color briefly to indicate action
		originalColor := node.GetColor()
		node.SetColor(tcell.ColorLightCyan)

		// Schedule reset of color
		go func() {
			time.Sleep(150 * time.Millisecond)
			app.QueueUpdateDraw(func() {
				node.SetColor(originalColor)
			})
		}()
	}

	// Setup event handlers for tree expand/collapse with additional feedback
	workloadTree.SetSelectedFunc(treeExpandCollapseAndDetailsFunc)
	infraTree.SetSelectedFunc(treeExpandCollapseAndDetailsFunc)
	// Storage tree has its own custom handler for PV detail popups

	// Add pages for each view
	tabs.AddPage("workload", workloadTree, true, true)
	tabs.AddPage("infra", infraTree, true, false)
	tabs.AddPage("storage", storageTree, true, false)

	// Create status bar with more detailed info
	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[yellow::b]TAB[white]: Switch Views | [yellow::b]↑/↓[white]: Navigate | [yellow::b]ENTER[white]: Expand/Collapse | [yellow::b]q[white]: Quit")

	// Add components to layout with enhanced title
	titleBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[::b]Kubernetes Resource Inspector[white] - Namespace: [green::b]%s", instanceID))

	// Create help text
	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[green]Active view: Workload[white] - Use TAB to toggle views")

	// Create a horizontal flex for the help and view indicator
	helpFlex := tview.NewFlex().
		AddItem(helpText, 0, 1, false)

	// Add all components to the main layout
	flex.AddItem(titleBar, 1, 1, false).
		AddItem(helpFlex, 1, 1, false).
		AddItem(tabs, 0, 1, true).
		AddItem(statusBar, 1, 1, false)

	// Handle keyboard input with enhanced feedback
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Switch between views cyclically
			if tabs.HasPage("workload") {
				currentPage, _ := tabs.GetFrontPage()
				switch currentPage {
				case "workload":
					tabs.SwitchToPage("infra")
					helpText.SetText("[blue]Active view: Infrastructure[white] - Use TAB to toggle views")
				case "infra":
					tabs.SwitchToPage("storage")
					helpText.SetText("[purple]Active view: Storage[white] - Use TAB to toggle views")
				default:
					tabs.SwitchToPage("workload")
					helpText.SetText("[green]Active view: Workload[white] - Use TAB to toggle views")
				}
			}
		} else if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	// Start the application with all TUI enhancements
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	return nil
}
