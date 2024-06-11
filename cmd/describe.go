package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	describeServiceID string
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:     "describe [--service-id SERVICE_ID]",
	Short:   "Describe service",
	Long:    `Describe service for a given service id.`,
	Example: `  ./omnistrate-ctl describe --service-id SERVICE_ID`,
	RunE:    runDescribe,
}

func init() {
	rootCmd.AddCommand(describeCmd)

	describeCmd.Flags().StringVarP(&describeServiceID, "service-id", "", "", "service id")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer resetDescribe()

	// Validate input arguments
	if len(describeServiceID) == 0 {
		return fmt.Errorf("must provide --service-id")
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		return fmt.Errorf("unable to retrieve authentication credentials, %s", err.Error())
	}

	// Describe service
	service, err := describeService(describeServiceID, token)
	if err != nil {
		fmt.Println("Error describing service:", err.Error())
		return err
	}

	// Print service details
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Println(bold("Service Details:"))
	fmt.Printf("%s %s\n", green("Service ID:"), service.ID)
	fmt.Printf("%s %s\n", green("Service Name:"), service.Name)
	fmt.Printf("%s %s\n", green("Service Created At:"), service.CreatedAt)
	fmt.Printf("%s %s\n", green("Service Description:"), service.Description)
	if service.ServiceLogoURL != nil {
		fmt.Printf("%s %s\n", green("Service Logo URL:"), *service.ServiceLogoURL)
	}

	fmt.Println(bold("\nEnvironments:"))

	for _, env := range service.ServiceEnvironments {
		fmt.Printf("\n%s %s\n", green("Environment ID:"), env.ID)
		fmt.Printf("%s %s\n", green("Environment Name:"), env.Name)
		fmt.Printf("%s %s\n", green("Environment Visibility:"), env.Visibility)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.Style().Options.SeparateRows = true

		t.AppendHeader(table.Row{bold("Product Tier ID"), bold("Name"), bold("Description"), bold("Tenancy Type"), bold("Deployment Type")})

		for _, plan := range env.ServicePlans {
			t.AppendRow(table.Row{green(plan.ProductTierID), plan.Name, plan.Description, plan.TierType, plan.ModelType})
		}

		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, Align: text.AlignLeft, WidthMax: 30},
			{Number: 2, Align: text.AlignLeft, WidthMax: 30},
			{Number: 3, Align: text.AlignLeft, WidthMax: 50},
			{Number: 4, Align: text.AlignLeft, WidthMax: 30},
			{Number: 5, Align: text.AlignLeft, WidthMax: 30},
		})

		t.Render()
	}

	return nil
}

func describeService(serviceId, token string) (*serviceapi.DescribeServiceResult, error) {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, fmt.Errorf("unable to describe service, %s", err.Error())
	}

	request := serviceapi.DescribeServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceId),
	}

	res, err := service.DescribeService(context.Background(), &request)
	if err != nil {
		return nil, fmt.Errorf("unable to describe service, %s", err.Error())
	}
	return res, nil
}

func resetDescribe() {
	describeServiceID = ""
}
