package service

import (
	"fmt"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	serviceExample = `  # Get all services
  omnistrate-ctl get service

  # Get service with name
  omnistrate-ctl get service <name>

  # Get multiple services with names
  omnistrate-ctl get service <name1> <name2> <name3>

  # Get service with ID
  omnistrate-ctl get service <id> --id

  # Get multiple services with IDs
  omnistrate-ctl get service <id1> <id2> <id3> --id`
)

// ServiceCmd represents the describe command
var ServiceCmd = &cobra.Command{
	Use:          "service <name>",
	Short:        "Display one or more services",
	Long:         `The get service command displays basic information about one or more services.`,
	Example:      serviceExample,
	RunE:         Run,
	SilenceUsage: true,
}

func init() {
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
		for _, id := range args {
			var service *serviceapi.DescribeServiceResult
			service, err = dataaccess.DescribeService(id, token)
			if err != nil {
				utils.PrintError(err)
				continue
			}
			services = append(services, service)
		}
	} else {
		// List services
		var listRes *serviceapi.ListServiceResult
		listRes, err = dataaccess.ListServices(token)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		allServices := listRes.Services

		// Print services table if no service name is provided
		if len(args) == 0 {
			utils.PrintSuccess(fmt.Sprintf("%d service(s) found", len(allServices)))
			if len(allServices) > 0 {
				printTable(allServices)
			}
			return nil
		}

		// Format allServices into a map
		serviceMap := make(map[string]*serviceapi.DescribeServiceResult)
		for _, service := range allServices {
			serviceMap[service.Name] = service
		}

		// Filter services by name
		for _, name := range args {
			service, ok := serviceMap[name]
			if !ok {
				utils.PrintError(fmt.Errorf("service '%s' not found", name))
				continue
			}
			services = append(services, service)
		}
	}

	printTable(services)

	return nil
}

func printTable(services []*serviceapi.DescribeServiceResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "ID\tName\tEnvironments")

	for _, service := range services {
		envNames := []string{}
		for _, env := range service.ServiceEnvironments {
			envNames = append(envNames, env.Name)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			service.ID,
			service.Name,
			strings.Join(envNames, ","))
	}

	w.Flush()
}
