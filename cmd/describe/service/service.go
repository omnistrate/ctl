package service

import (
	"encoding/json"
	"fmt"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

var (
	serviceExample = `  # Describe the service with name
  omnistrate-ctl describe service <name>

  # Describe the service with ID
  omnistrate-ctl describe service <id> --id`
)

// ServiceCmd represents the describe command
var ServiceCmd = &cobra.Command{
	Use:          "service <name>",
	Short:        "Display details for one or more services",
	Long:         "Display detailed information about the service by specifying the service name or ID.",
	Example:      serviceExample,
	RunE:         Run,
	SilenceUsage: true,
}

func init() {
	ServiceCmd.Args = cobra.MinimumNArgs(1) // Require at least one argument

	ServiceCmd.Flags().Bool("id", false, "Specify service ID instead of name")
}

func Run(cmd *cobra.Command, args []string) error {
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

	var services []*serviceapi.DescribeServiceResult
	if ID {
		for _, name := range args {
			service, err := dataaccess.DescribeService(name, token)
			if err != nil {
				utils.PrintError(err)
				return err
			}
			services = append(services, service)
		}
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
				services = append(services, s)
				found[s.Name] = true
			}
		}

		namesNotFound := make([]string, 0)
		for name, ok := range found {
			if !ok {
				namesNotFound = append(namesNotFound, name)
			}
		}

		if len(namesNotFound) > 0 {
			err = fmt.Errorf("service(s) not found: %s", strings.Join(namesNotFound, ", "))
			utils.PrintError(err)
			return err
		}
	}

	// Print service details
	for _, service := range services {
		data, err := json.MarshalIndent(service, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
