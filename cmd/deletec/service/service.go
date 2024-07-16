package service

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/pkg/errors"
	"slices"
	"strings"

	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	serviceExample = `  # Delete service with name
  omnistrate-ctl delete service <name>

  # Delete service with ID
  omnistrate-ctl delete service <ID> --id

  # Delete multiple services with names
  omnistrate-ctl delete service <name1> <name2> <name3>

  # Delete multiple services with IDs
  omnistrate-ctl delete service <ID1> <ID2> <ID3> --id`
)

// ServiceCmd represents the delete command
var ServiceCmd = &cobra.Command{
	Use:          "service <name>",
	Short:        "Delete one or more services",
	Long:         `Delete service with name or ID. Use --id to specify ID. If not specified, name is assumed.`,
	Example:      serviceExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	ServiceCmd.Args = cobra.MinimumNArgs(1) // Require at least one argument

	ServiceCmd.Flags().Bool("id", false, "Specify service ID instead of name")
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var ID bool
	ID, err = cmd.Flags().GetBool("id")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var serviceIDs []string
	if ID {
		serviceIDs = args
	} else {
		// List services
		listRes, err := dataaccess.ListServices(token)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		found := make(map[string]bool)
		for _, name := range args {
			found[name] = false
		}

		// Filter services by name
		for _, s := range listRes.Services {
			if slices.Contains(args, s.Name) {
				serviceIDs = append(serviceIDs, string(s.ID))
				found[s.Name] = true
			}
		}

		// Check if all services were found
		servicesNotFound := make([]string, 0)
		for name, ok := range found {
			if !ok {
				servicesNotFound = append(servicesNotFound, name)
			}
		}

		if len(servicesNotFound) > 0 {
			err = errors.New("services not found: " + strings.Join(servicesNotFound, ", "))
			utils.PrintError(err)
			return err
		}
	}

	// Delete service
	for _, serviceID := range serviceIDs {
		err = dataaccess.DeleteService(serviceID, token)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}

	utils.PrintSuccess("Service(s) deleted successfully")

	return nil
}
