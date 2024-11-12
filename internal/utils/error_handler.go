package utils

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"os"
)

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

func CheckSpecFileExists(spinner *ysmrr.Spinner, sm ysmrr.SpinnerManager, specFile string) (err error) {
	if _, err = os.Stat(specFile); !os.IsNotExist(err) {
		err = fmt.Errorf("unable to access spec file: %s", specFile)
		HandleSpinnerError(spinner, sm, err)
		return err
	}

	if err == nil {
		// Spec file exists
		err = fmt.Errorf("current directory is already initialized with an existing project")
		HandleSpinnerError(spinner, sm, err)
		return err
	}

	return nil
}
