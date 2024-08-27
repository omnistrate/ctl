package instance

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	createExample = `  # Create an instance deployment
  omctl instance create --service=mysql --environment=dev --plan=mysql --version=latest --resource=mySQL --cloud-provider=aws --region=ca-central-1 --param '{"databaseName":"default","password":"a_secure_password","rootPassword":"a_secure_root_password","username":"user"}'`
)

var InstanceID string

var createCmd = &cobra.Command{
	Use:          "create --service=[service] --environment=[environment] --plan=[plan] --version=[version] --resource=[resource] --cloud-provider=[aws|gcp] --region=[region] [--param=param] [--param-file=file-path]",
	Short:        "Create an instance deployment",
	Long:         `This command helps you create an instance deployment for your service.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Flags().String("service", "", "Service name")
	createCmd.Flags().String("environment", "", "Environment name")
	createCmd.Flags().String("plan", "", "Service plan name")
	createCmd.Flags().String("version", "preferred", "Service plan version (latest|preferred|1.0 etc.)")
	createCmd.Flags().String("resource", "", "Resource name")
	createCmd.Flags().String("cloud-provider", "", "Cloud provider (aws|gcp)")
	createCmd.Flags().String("region", "", "Region code (e.g. us-east-2, us-central1)")
	createCmd.Flags().String("param", "", "Parameters for the instance deployment")
	createCmd.Flags().String("param-file", "", "Json file containing parameters for the instance deployment")
	createCmd.Flags().StringP("subscription-id", "", "", "Subscription ID to use for the instance deployment. If not provided, instance deployment will be created in your own subscription.")
	createCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")

	if err := createCmd.MarkFlagRequired("service"); err != nil {
		return
	}
	if err := createCmd.MarkFlagRequired("environment"); err != nil {
		return
	}
	if err := createCmd.MarkFlagRequired("plan"); err != nil {
		return
	}
	if err := createCmd.MarkFlagRequired("resource"); err != nil {
		return
	}
	if err := createCmd.MarkFlagRequired("cloud-provider"); err != nil {
		return
	}
	if err := createCmd.MarkFlagRequired("region"); err != nil {
		return
	}
	if err := createCmd.MarkFlagFilename("param-file"); err != nil {
		return
	}
	createCmd.MarkFlagsMutuallyExclusive("param", "param-file")

	createCmd.Args = cobra.NoArgs // Require no arguments
}

func runCreate(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	service, err := cmd.Flags().GetString("service")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	environment, err := cmd.Flags().GetString("environment")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	plan, err := cmd.Flags().GetString("plan")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	version, err := cmd.Flags().GetString("version")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	resource, err := cmd.Flags().GetString("resource")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	cloudProvider, err := cmd.Flags().GetString("cloud-provider")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	param, err := cmd.Flags().GetString("param")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	paramFile, err := cmd.Flags().GetString("param-file")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	subscriptionID, err := cmd.Flags().GetString("subscription-id")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Creating instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if resource exists
	serviceID, _, productTierID, _, err := getResource(token, service, environment, plan, resource)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get the version
	switch version {
	case "latest":
		version, err = dataaccess.FindLatestVersion(token, serviceID, productTierID)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
	case "preferred":
		version, err = dataaccess.FindPreferredVersion(token, serviceID, productTierID)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
	}

	// Check if the version exists
	_, err = dataaccess.DescribeVersionSet(token, serviceID, productTierID, version)
	if err != nil {
		if strings.Contains(err.Error(), "Version set not found") {
			err = errors.New(fmt.Sprintf("version %s not found", version))
		}
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe service offering
	res, err := dataaccess.DescribeServiceOffering(token, serviceID, productTierID, version)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	offering := res.ConsumptionDescribeServiceOfferingResult.Offerings[0]

	// Format parameters
	formattedParams, err := formatParams(param, paramFile)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	var resourceKey string
	found := false
	for _, resourceEntity := range offering.ResourceParameters {
		if strings.EqualFold(resourceEntity.Name, resource) {
			found = true
			resourceKey = resourceEntity.URLKey
		}
	}

	if !found {
		err = fmt.Errorf("resource %s not found in the service offering", resource)
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	request := inventoryapi.FleetCreateResourceInstanceRequest{
		ServiceProviderID:     inventoryapi.ServiceProviderID(res.ConsumptionDescribeServiceOfferingResult.ServiceProviderID),
		ServiceKey:            res.ConsumptionDescribeServiceOfferingResult.ServiceURLKey,
		ServiceAPIVersion:     offering.ServiceAPIVersion,
		ServiceEnvironmentKey: offering.ServiceEnvironmentURLKey,
		ServiceModelKey:       offering.ServiceModelURLKey,
		ProductTierKey:        offering.ProductTierURLKey,
		ProductTierVersion:    &version,
		ResourceKey:           resourceKey,
		CloudProvider:         &cloudProvider,
		Region:                &region,
		RequestParams:         formattedParams,
		NetworkType:           nil,
	}
	if subscriptionID != "" {
		request.SubscriptionID = (*inventoryapi.SubscriptionID)(commonutils.ToPtr(subscriptionID))
	}
	instance, err := dataaccess.CreateInstance(token, request)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	if res == nil || instance.ID == nil {
		err = errors.New("failed to create instance")
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully created instance")

	// Search for the instance
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", *instance.ID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if len(searchRes.ResourceInstanceResults) == 0 {
		err = errors.New("failed to find the created instance")
		utils.PrintError(err)
		return err
	}

	// Format instance
	formattedInstance := formatInstance(searchRes.ResourceInstanceResults[0], false)
	InstanceID = formattedInstance.InstanceID

	// Marshal instance to JSON
	data, err := json.MarshalIndent(formattedInstance, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	if err = utils.PrintTextTableJsonOutput(output, string(data)); err != nil {
		return err
	}

	return nil
}

// Helper functions

func getResource(token, serviceNameArg, environmentArg, planNameArg, resourceNameArg string) (serviceID, environmentID, productTierID, resourceID string, err error) {
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resource:%s", resourceNameArg))
	if err != nil {
		return
	}

	found := false
	for _, res := range searchRes.ResourceResults {
		if res == nil {
			continue
		}
		if strings.EqualFold(res.Name, resourceNameArg) &&
			strings.EqualFold(res.ServiceName, serviceNameArg) &&
			strings.EqualFold(res.ProductTierName, planNameArg) &&
			strings.EqualFold(res.ServiceEnvironmentName, environmentArg) {
			found = true
			serviceID = string(res.ServiceID)
			environmentID = string(res.ServiceEnvironmentID)
			productTierID = string(res.ProductTierID)
			resourceID = res.ID
			break
		}
	}

	if !found {
		err = fmt.Errorf("target resource not found. Please check input values and try again")
		return
	}

	return
}

func formatInstance(instance *inventoryapi.ResourceInstanceSearchRecord, truncateNames bool) model.Instance {
	planName := ""
	if instance.ProductTierName != nil {
		planName = *instance.ProductTierName
	}
	planVersion := ""
	if instance.ProductTierVersion != nil {
		planVersion = *instance.ProductTierVersion
	}
	serviceName := instance.ServiceName
	if truncateNames {
		serviceName = utils.TruncateString(serviceName, defaultMaxNameLength)
		planName = utils.TruncateString(planName, defaultMaxNameLength)
	}
	subscriptionID := ""
	if instance.SubscriptionID != nil {
		subscriptionID = string(*instance.SubscriptionID)
	}

	formattedInstance := model.Instance{
		InstanceID:     instance.ID,
		Service:        serviceName,
		Environment:    instance.ServiceEnvironmentName,
		Plan:           planName,
		Version:        planVersion,
		Resource:       instance.ResourceName,
		CloudProvider:  string(instance.CloudProvider),
		Region:         instance.RegionCode,
		Status:         string(instance.Status),
		SubscriptionID: subscriptionID,
	}

	return formattedInstance
}

func formatParams(param, paramFile string) (formattedParams map[string]string, err error) {
	// Read parameters from file if provided
	if paramFile != "" {
		var fileContent []byte
		fileContent, err = os.ReadFile(paramFile)
		if err != nil {
			return
		}
		param = string(fileContent)
	}

	// Extract parameters from json format param
	if param != "" {
		err = json.Unmarshal([]byte(param), &formattedParams)
		if err != nil {
			return
		}
	}

	return
}
