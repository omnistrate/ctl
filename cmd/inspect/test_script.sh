#!/bin/bash
# Test script for inspect command

# Instructions:
# 1. Build the CLI with: make build
# 2. Test the text mode version (non-interactive)
echo "Testing text mode (non-interactive)..."
./dist/omnistrate-ctl-$(go env GOOS)-$(go env GOARCH) inspect my-test-namespace --text

# 3. Test with a real Kubernetes namespace
echo ""
echo "To inspect a real Kubernetes namespace, make sure you have 'kubectl' configured and run:"
echo "./dist/omnistrate-ctl-$(go env GOOS)-$(go env GOARCH) inspect <namespace-name>"
echo ""
echo "Optional Kubernetes flags:"
echo "--kubeconfig   Path to kubeconfig file (default is ~/.kube/config)"
echo "--context      Kubernetes context to use"

# 4. Test the interactive TUI mode
echo ""
echo "To test interactive colorized TUI mode, run:"
echo "./dist/omnistrate-ctl-$(go env GOOS)-$(go env GOARCH) inspect <namespace-name>"
echo ""
echo "Navigation instructions:"
echo "- Use TAB to switch between Workload and Infrastructure views"
echo "- Use arrow keys (↑/↓) to navigate through the tree"
echo "- Use ENTER to expand/collapse nodes (with visual feedback)"
echo "- Press 'q' to quit the TUI"
echo ""
echo "Features:"
echo "- Color-coded resources based on type and status"
echo "- Selectable and navigable elements"
echo "- Interactive visual feedback when expanding/collapsing"
echo "- Status indicators for pods (✅ Running, ⏳ Pending, ❌ Failed)"
echo "- Detailed VM information with instance types and resources"