package service

import (
	"encoding/json"
	"errors"
	"fmt"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	serviceExample = `  # Describe the service with the name
  omnistrate-ctl describe service <name>`
)

// ServiceCmd represents the describe command
var ServiceCmd = &cobra.Command{
	Use:          "service <name>",
	Short:        "Describe service",
	Long:         `The describe service command displays detailed information about a service.`,
	Example:      serviceExample,
	RunE:         Run,
	SilenceUsage: true,
}

func init() {
	ServiceCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func Run(cmd *cobra.Command, args []string) error {
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
	var service *serviceapi.DescribeServiceResult
	var found bool
	for _, s := range listRes.Services {
		if s.Name == args[0] {
			service = s
			found = true
			break
		}
	}

	if !found {
		utils.PrintError(errors.New("service not found"))
		return nil
	}

	// Print service details
	data, err := json.MarshalIndent(service, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
