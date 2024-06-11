package utils

import (
	"fmt"
	"github.com/fatih/color"
)

func PrintError(err error) {
	errorMsg := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Println(errorMsg("Error: "), err)
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

func PrintBold(msg string) {
	bold := color.New(color.Bold).SprintFunc()
	fmt.Println(bold(msg))
}
