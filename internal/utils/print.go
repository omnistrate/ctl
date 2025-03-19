package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/model"
)

func PrintError(err error) {
	errorMsg := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Println(errorMsg("Error: "), err.Error())
	if !config.IsDryRun() {
		os.Exit(1)
	}
}

func PrintSuccess(msg string) {
	successMsg := color.New(color.FgGreen, color.Bold).SprintFunc()
	fmt.Println(successMsg(msg))
}

func PrintInfo(msg string) {
	infoMsg := color.New(color.FgCyan).SprintFunc()
	fmt.Println(infoMsg(msg))
}

func PrintWarning(msg string) {
	warningMsg := color.New(color.FgYellow).SprintFunc()
	fmt.Println(warningMsg(msg))
}

func PrintURL(label, url string) {
	urlMsg := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("%s: %s\n", label, urlMsg(url))
}

func PrintJSON(res interface{}) {
	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		PrintError(err)
		return
	}
	fmt.Println(string(data))
}

func PrintTextTableJsonArrayOutput[T any](output string, objects []T) error {
	switch output {
	case "text":
		dataArray := make([]string, 0)
		for _, obj := range objects {
			data, err := json.MarshalIndent(obj, "", "    ")
			if err != nil {
				return err
			}
			dataArray = append(dataArray, string(data))
		}
		return PrintText(dataArray)
	case "table":
		dataArray := make([]string, 0)
		for _, obj := range objects {
			data, err := json.MarshalIndent(obj, "", "    ")
			if err != nil {
				return err
			}
			dataArray = append(dataArray, string(data))
		}
		return PrintTable(dataArray)
	case "json":
		data, err := json.MarshalIndent(objects, "", "    ")
		if err != nil {
			return err
		}
		if len(objects) == 0 {
			fmt.Println("[]")
		} else {
			fmt.Printf("%s\n", data)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", output)
	}
	return nil
}

func PrintTextTableJsonOutput[T any](output string, object T) error {
	data, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		return err
	}

	switch output {
	case "text":
		return PrintText([]string{string(data)})
	case "table":
		return PrintTable([]string{string(data)})
	case "json":
		fmt.Printf("%s\n", data)
	default:
		return fmt.Errorf("unsupported output format: %s", output)
	}
	return nil
}

// PrintUpgradeStatuses prints upgrade statuses in a table format
func PrintUpgradeStatuses(statuses []*model.UpgradeStatus) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"UPGRADE ID", "TOTAL", "PENDING", "IN PROGRESS", "COMPLETED", "FAILED", "SCHEDULED", "SKIPPED", "STATUS"})

	for _, s := range statuses {
		scheduled := ""
		if s.Scheduled != nil {
			scheduled = fmt.Sprintf("%d", *s.Scheduled)
		}
		t.AppendRow(table.Row{
			s.UpgradeID,
			s.Total,
			s.Pending,
			s.InProgress,
			s.Completed,
			s.Failed,
			scheduled,
			s.Skipped,
			s.Status,
		})
	}

	t.Render()
}

// FromInt64Ptr converts int64 pointer to int pointer
func FromInt64Ptr(val *int64) *int {
	if val == nil {
		return nil
	}
	result := int(*val)
	return &result
}
