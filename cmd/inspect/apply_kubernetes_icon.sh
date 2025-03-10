#\!/bin/bash

# Update pod icons in analyze.go
sed -i '' -e 's/podIcon = "✅"/podColor = tcell.ColorGreen/g' \
          -e 's/podIcon = "⏳"/podColor = tcell.ColorYellow/g' \
          -e 's/podIcon = "❌"/podColor = tcell.ColorRed/g' \
          -e 's/podIcon = "❔"/podColor = tcell.ColorGray/g' \
          -e 's/%s Pod: %s (%s)", podIcon, pod.Name, pod.Status/⎈ Pod: %s (%s)", pod.Name, pod.Status/g' \
          cmd/inspect/analyze.go

# Update pod icons in mockup.go
sed -i '' -e 's/✅ Pod:/⎈ Pod:/g' \
          -e 's/⏳ Pod:/⎈ Pod:/g' \
          -e 's/❌ Pod:/⎈ Pod:/g' \
          -e 's/❔ Pod:/⎈ Pod:/g' \
          cmd/inspect/mockup.go

echo "Kubernetes pod icons updated to '⎈' in all files."
