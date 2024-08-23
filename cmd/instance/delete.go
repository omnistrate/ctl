package instance

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/input"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

const (
	deleteExample = `  # Delete instance
  omctl instance delete instance-abcd1234`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [instance-id] [flags]",
	Short:        "Delete an instance",
	Long:         `This command helps you delete an instance from your account.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	deleteCmd.Flags().BoolP("yes", "y", false, "Pre-approve the deletion of the instance without prompting for confirmation")
	deleteCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func runDelete(cmd *cobra.Command, args []string) error {
	instanceId := args[0]

	// Get flags
	output, _ := cmd.Flags().GetString("output")
	yes, _ := cmd.Flags().GetBool("yes")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Confirm deletion
	if !yes {
		ok, err := prompt.New().Ask("Are you sure you want to delete this instance? (y/n)").
			Input("", input.WithValidateFunc(
				func(input string) error {
					if slices.Contains([]string{"y", "yes", "n", "no"}, strings.ToLower(input)) {
						return nil
					} else {
						return errors.New("invalid input")
					}
				}))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if !slices.Contains([]string{"y", "yes"}, strings.ToLower(ok)) {
			return nil
		}
	}

	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Deleting instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if the instance exists
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", instanceId))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var found bool
	var serviceId, environmentId, resourceId string
	for _, instance := range searchRes.ResourceInstanceResults {
		if instance.ID == instanceId {
			serviceId = string(instance.ServiceID)
			environmentId = string(instance.ServiceEnvironmentID)
			if instance.ResourceID == nil {
				err = fmt.Errorf("resource ID not returned for instance %s", instanceId)
				utils.PrintError(err)
				return err
			}
			resourceId = string(*instance.ResourceID)
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceId)
		utils.PrintError(err)
		return nil
	}

	err = dataaccess.DeleteInstance(token, serviceId, environmentId, resourceId, instanceId)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if output != "json" {
		spinner.UpdateMessage("Successfully deleted instance")
		spinner.Complete()
		sm.Stop()
	}

	return nil
}
