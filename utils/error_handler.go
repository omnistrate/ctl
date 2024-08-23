package utils

import "github.com/chelnak/ysmrr"

func HandleSpinnerError(spinner *ysmrr.Spinner, sm ysmrr.SpinnerManager, err error) {
	if spinner != nil {
		spinner.Error()
		sm.Stop()
	}
	PrintError(err)
}

func HandleSpinnerSuccess(spinner *ysmrr.Spinner, sm ysmrr.SpinnerManager, message string) {
	if spinner != nil {
		spinner.UpdateMessage(message)
		spinner.Complete()
		sm.Stop()
	}
}
