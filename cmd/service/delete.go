package service

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	deleteExample = `  # Delete service with name
  omctl service delete [service-name]

  # Delete service with ID
  omctl service delete --id=[service-ID]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [service-name] [flags]",
	Short:        "Delete a service",
	Long:         `This command helps you delete a service using its name or ID.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Args = cobra.ExactArgs(1) // Require at least one argument

	deleteCmd.Flags().String("id", "", "Service ID")
}

func runDelete(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	var name string
	if len(args) > 0 {
		name = args[0]
	}

	// Retrieve flags
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate input args
	err = validateDeleteInput(name, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Check if service exists
	id, name, err = getService(token, name, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Delete service
	err = dataaccess.DeleteService(token, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	utils.PrintSuccess("Service deleted successfully.")

	return nil
}

func validateDeleteInput(serviceNameArg, serviceIDArg string) error {
	if serviceNameArg == "" && serviceIDArg == "" {
		return errors.New("service name or ID must be provided")
	}

	if serviceNameArg != "" && serviceIDArg != "" {
		return errors.New("only one of service name or ID can be provided")
	}

	return nil
}
