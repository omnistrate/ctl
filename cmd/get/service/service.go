package service

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	serviceExample = `
		# List all services
		kubectl get services

		# List the service with the name 'my-service'
		kubectl get service my-service`
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
	listRes, err := listServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print services table if no service name is provided
	if len(args) == 0 {
		printTable(listRes.Services)
		return nil
	}

	// Format listRes.Services into a map
	serviceMap := make(map[string]*serviceapi.DescribeServiceResult)
	for _, service := range listRes.Services {
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

func listServices(token string) (*serviceapi.ListServiceResult, error) {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceapi.List{
		Token: token,
	}

	res, err := service.ListService(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func printTable(services []*serviceapi.DescribeServiceResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Service ID\tName\tCreated At\tEnvironments")

	for _, service := range services {
		envNames := []string{}
		for _, env := range service.ServiceEnvironments {
			envNames = append(envNames, env.Name)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			service.ID,
			service.Name,
			service.CreatedAt,
			strings.Join(envNames, ","))
	}

	w.Flush()
}
