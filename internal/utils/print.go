package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/omnistrate/ctl/internal/config"
)

var (
	LastPrintedString string
)

func PrintError(err error) {
	errorMsg := color.New(color.FgRed, color.Bold).SprintFunc()
	msg := fmt.Sprintf("%s %s", errorMsg("Error: "), err.Error())
	fmt.Println(msg)
	LastPrintedString = msg
	if !config.IsDryRun() {
		os.Exit(1)
	}
}

func PrintSuccess(msg string) {
	successMsg := color.New(color.FgGreen, color.Bold).SprintFunc()
	formatted := successMsg(msg)
	fmt.Println(formatted)
	LastPrintedString = formatted
}

func PrintInfo(msg string) {
	infoMsg := color.New(color.FgCyan).SprintFunc()
	formatted := infoMsg(msg)
	fmt.Println(formatted)
	LastPrintedString = formatted
}

func PrintWarning(msg string) {
	warningMsg := color.New(color.FgYellow).SprintFunc()
	formatted := warningMsg(msg)
	fmt.Println(formatted)
	LastPrintedString = formatted
}

func PrintURL(label, url string) {
	urlMsg := color.New(color.FgCyan).SprintFunc()
	formatted := fmt.Sprintf("%s: %s", label, urlMsg(url))
	fmt.Println(formatted)
	LastPrintedString = formatted
}

func PrintJSON(res interface{}) {
	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		PrintError(err)
		return
	}
	formatted := string(data)
	fmt.Println(formatted)
	LastPrintedString = formatted
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
		err := PrintText(dataArray)
		if err == nil {
			LastPrintedString = fmt.Sprintf("%v", dataArray)
		}
		return err
	case "table":
		dataArray := make([]string, 0)
		for _, obj := range objects {
			data, err := json.MarshalIndent(obj, "", "    ")
			if err != nil {
				return err
			}
			dataArray = append(dataArray, string(data))
		}
		err := PrintTable(dataArray)
		if err == nil {
			LastPrintedString = fmt.Sprintf("%v", dataArray)
		}
		return err
	case "json":
		data, err := json.MarshalIndent(objects, "", "    ")
		if err != nil {
			return err
		}
		if len(objects) == 0 {
			formatted := "[]"
			fmt.Println(formatted)
			LastPrintedString = formatted
		} else {
			formatted := string(data)
			fmt.Printf("%s\n", formatted)
			LastPrintedString = formatted
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
		err := PrintText([]string{string(data)})
		if err == nil {
			LastPrintedString = string(data)
		}
		return err
	case "table":
		err := PrintTable([]string{string(data)})
		if err == nil {
			LastPrintedString = string(data)
		}
		return err
	case "json":
		formatted := string(data)
		fmt.Printf("%s\n", formatted)
		LastPrintedString = formatted
	default:
		return fmt.Errorf("unsupported output format: %s", output)
	}
	return nil
}
