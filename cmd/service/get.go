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
	getExample = `  # Get all services
  omnistrate-ctl service get

  # Get service with name
  omnistrate-ctl service get <name>

  # Get multiple services with names
  omnistrate-ctl service get <name1> <name2> <name3>

  # Get service with ID
  omnistrate-ctl service get <id> --id

  # Get multiple services with IDs
  omnistrate-ctl service get <id1> <id2> <id3> --id`
)

var getCmd = &cobra.Command{
	Use:          "get",
	Short:        "Display one or more services (deprecated. Please use 'service list' instead)",
	Long:         `The service get command displays basic information about one or more services.`,
	Example:      getExample,
	RunE:         runGet,
	SilenceUsage: true,
}

func init() {
	getCmd.Flags().Bool("id", false, "Specify service ID instead of name")
}

func runGet(cmd *cobra.Command, args []string) error {
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
			service, err = dataaccess.DescribeService(token, id)
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

	_, err := fmt.Fprintln(w, "ID\tName\tEnvironments")
	if err != nil {
		return
	}

	for _, service := range services {
		envNames := []string{}
		for _, env := range service.ServiceEnvironments {
			envNames = append(envNames, env.Name)
		}
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\n",
			service.ID,
			service.Name,
			strings.Join(envNames, ","))
		if err != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
