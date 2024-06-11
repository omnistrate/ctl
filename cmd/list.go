package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
	"sort"
	"time"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List service",
	Long:         `List service. The service must be created before it can be listed.`,
	Example:      `  ./omnistrate-ctl list`,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(listCmd)
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

	// Parse dates and sort services by CreatedAt in descending order
	sort.Slice(res.Services, func(i, j int) bool {
		timeI, errI := time.Parse(time.RFC3339, res.Services[i].CreatedAt)
		timeJ, errJ := time.Parse(time.RFC3339, res.Services[j].CreatedAt)
		if errI != nil || errJ != nil {
			return false
		}
		return timeI.After(timeJ)
	})

	// Print service details
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("%s %d\n\n", green("Total Services:"), len(res.Services))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{bold("Created At"), bold("ID"), bold("Name"), bold("Description")})

	for _, service := range res.Services {
		t.AppendRow(table.Row{service.CreatedAt, green(service.ID), service.Name, service.Description})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, WidthMax: 30},
		{Number: 2, Align: text.AlignLeft, WidthMax: 30},
		{Number: 3, Align: text.AlignLeft, WidthMax: 50},
		{Number: 4, Align: text.AlignLeft, WidthMax: 50},
	})

	t.Render()

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
