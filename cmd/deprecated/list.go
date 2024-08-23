package deprecated

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"time"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List all available services (deprecated)",
	Long:         `The list command retrieves and displays a list of all available services that have been created.`,
	Example:      `  omctl list`,
	RunE:         runList,
	SilenceUsage: true,
}

func runList(cmd *cobra.Command, args []string) error {
	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// List service
	res, err := listServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print services
	printServicesTable(res.Services)

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

func printServicesTable(services []*serviceapi.DescribeServiceResult) {
	// Parse dates and sort services by CreatedAt in descending order
	sort.Slice(services, func(i, j int) bool {
		timeI, errI := time.Parse(time.RFC3339, services[i].CreatedAt)
		timeJ, errJ := time.Parse(time.RFC3339, services[j].CreatedAt)
		if errI != nil || errJ != nil {
			return false
		}
		return timeI.After(timeJ)
	})

	// Print service details
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("%s %d\n\n", green("Total Services:"), len(services))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{bold("Created At"), bold("ID"), bold("Name"), bold("Description")})

	for _, service := range services {
		t.AppendRow(table.Row{service.CreatedAt, green(service.ID), service.Name, service.Description})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, WidthMax: 30},
		{Number: 2, Align: text.AlignLeft, WidthMax: 30},
		{Number: 3, Align: text.AlignLeft, WidthMax: 50},
		{Number: 4, Align: text.AlignLeft, WidthMax: 50},
	})

	t.Render()
}
