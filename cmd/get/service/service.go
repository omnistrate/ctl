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
	serviceExample = `  # List all services
  omnistrate-ctl get services

  # List the service with the name
  omnistrate-ctl get service <name>`
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
	allServices := listRes.Services

	// Print services table if no service name is provided
	if len(args) == 0 {
		utils.PrintSuccess(fmt.Sprintf("%d services found", len(allServices)))
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
	var services []*serviceapi.DescribeServiceResult
	for _, name := range args {
		service, ok := serviceMap[name]
		if !ok {
			utils.PrintError(fmt.Errorf("service '%s' not found", name))
			continue
		}
		services = append(services, service)
	}

	// Print service details per service if service name is provided
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
