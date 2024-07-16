package service

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/pkg/errors"

	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	serviceExample = `
		# Delete the service with name
		omnistrate-ctl delete service <name>`
)

// ServiceCmd represents the delete command
var ServiceCmd = &cobra.Command{
	Use:          "service <name>",
	Short:        "Delete a service",
	Long:         ``,
	Example:      serviceExample,
	RunE:         run,
	SilenceUsage: true,
}

func run(cmd *cobra.Command, args []string) error {
	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// List services
	listRes, err := dataaccess.ListServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Filter services by name
	var serviceID string
	var found bool
	for _, s := range listRes.Services {
		if s.Name == args[0] {
			serviceID = string(s.ID)
			found = true
			break
		}
	}

	if !found {
		utils.PrintError(errors.New("service not found"))
		return nil
	}

	// Delete service
	err = dataaccess.DeleteService(serviceID, token)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess("Service deleted successfully")

	return nil
}
