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

const (
	describeExample = `  # Describe service with name
  omctl service describe <name>

  # Describe service with ID
  omctl service describe <id> --id

  # Describe multiple services with names
  omctl service describe <name1> <name2> <name3>

  # Describe multiple services with IDs
  omctl service describe <id1> <id2> <id3> --id`
)

var describeCmd = &cobra.Command{
	Use:          "describe [flags]",
	Short:        "Display details for one or more services",
	Long:         "Display detailed information about the service by specifying the service name or ID",
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.MinimumNArgs(1) // Require at least one argument

	describeCmd.Flags().Bool("id", false, "Specify service ID instead of name")
}

func runDescribe(cmd *cobra.Command, args []string) error {
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
			service, err := dataaccess.DescribeService(token, name)
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
				service, err := dataaccess.DescribeService(token, string(s.ID))
				if err != nil {
					utils.PrintError(err)
					return err
				}
				services = append(services, service)
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
