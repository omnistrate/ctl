package utils

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/omnistrate/commons/pkg/utils"
	"os"
)

func PrintError(err error) {
	errorMsg := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Println(errorMsg("Error: "), err.Error())
	if !utils.IsDryRun() {
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

func PrintTextTableJsonArrayOutput(output string, jsonData []string) error {
	switch output {
	case "text":
		return PrintText(jsonData)
	case "table":
		return PrintTable(jsonData)
	case "json":
		fmt.Printf("%+v\n", jsonData)
	default:
		return fmt.Errorf("unsupported output format: %s", output)
	}
	return nil
}

func PrintTextTableJsonOutput(output string, jsonData string) error {
	switch output {
	case "text":
		return PrintText([]string{jsonData})
	case "table":
		return PrintTable([]string{jsonData})
	case "json":
		fmt.Printf("%+v\n", jsonData)
	default:
		return fmt.Errorf("unsupported output format: %s", output)
	}
	return nil
}
