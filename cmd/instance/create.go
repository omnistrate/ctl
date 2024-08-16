package instance

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

const (
	createExample = `# Create an instance deployment
omnistrate instance create --service=mysql --environment=dev --plan=mysql --version=latest --resource=mySQL --cloud-provider=aws --region=ca-central-1 --param '{"databaseName":"default","password":"a_secure_password","rootPassword":"a_secure_root_password","username":"user"}'`
)

var InstanceID string

var createCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create an instance deployment",
	Long:         `This command helps you create an instance deployment for your service.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Flags().String("service", "", "Service name")
	createCmd.Flags().String("environment", "", "Environment type")
	createCmd.Flags().String("plan", "", "Service plan name")
	createCmd.Flags().String("version", "preferred", "Service plan version (latest|preferred|1.0 etc.)")
	createCmd.Flags().String("resource", "", "Resource name")
	createCmd.Flags().String("cloud-provider", "", "Cloud provider (aws|gcp)")
	createCmd.Flags().String("region", "", "Region code (e.g. us-east-2, us-central1)")
	createCmd.Flags().String("param", "", "Parameters for the instance deployment")
	createCmd.Flags().String("param-file", "", "File containing parameters for the instance deployment")
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
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Get flags
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
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Creating instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Get resource
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resource:%s", resource))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var serviceID, productTierID string
	found := false
	for _, res := range searchRes.ResourceResults {
		if res == nil {
			continue
		}
		if strings.ToLower(res.Name) == strings.ToLower(resource) &&
			strings.ToLower(res.ServiceName) == strings.ToLower(service) &&
			strings.ToLower(res.ProductTierName) == strings.ToLower(plan) &&
			res.ServiceEnvironmentType != nil && strings.ToLower(string(*res.ServiceEnvironmentType)) == strings.ToLower(environment) {
			found = true
			serviceID = string(res.ServiceID)
			productTierID = string(res.ProductTierID)
			break
		}
	}

	if !found {
		err = fmt.Errorf("target resource not found. Please check input values and try again")
		utils.PrintError(err)
		return err
	}

	// Get the version
	switch version {
	case "latest":
		version, err = dataaccess.FindLatestVersion(token, serviceID, productTierID)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	case "preferred":
		version, err = dataaccess.FindPreferredVersion(token, serviceID, productTierID)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}

	// Check if the version exists
	_, err = dataaccess.DescribeVersionSet(token, serviceID, productTierID, version)
	if err != nil {
		if strings.Contains(err.Error(), "Version set not found") {
			err = errors.New(fmt.Sprintf("version %s not found", version))
		}
		utils.PrintError(err)
		return err
	}

	// Describe service offering
	res, err := dataaccess.DescribeServiceOffering(token, serviceID, productTierID, version)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	offering := res.ConsumptionDescribeServiceOfferingResult.Offerings[0]

	// Read parameters from file if provided
	if paramFile != "" {
		fileContent, err := ioutil.ReadFile(paramFile)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		param = string(fileContent)
	}

	// Extract parameters from json format param
	var formattedParams map[string]interface{}
	if param != "" {
		err = json.Unmarshal([]byte(param), &formattedParams)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}

	var resourceKey string
	found = false
	for _, resourceEntity := range offering.ResourceParameters {
		if strings.ToLower(resourceEntity.Name) == strings.ToLower(resource) {
			found = true
			resourceKey = resourceEntity.URLKey
		}
	}

	if !found {
		err = fmt.Errorf("resource %s not found in the service offering", resource)
		utils.PrintError(err)
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
	instance, err := dataaccess.CreateInstance(token, request)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if res == nil || instance.ID == nil {
		err = errors.New("failed to create instance")
		utils.PrintError(err)
		return err
	}

	if spinner != nil {
		spinner.UpdateMessage("Successfully created instance")
		spinner.Complete()
		sm.Stop()
	}

	formattedInstance := model.Instance{
		ID:            *instance.ID,
		Service:       service,
		Environment:   environment,
		Plan:          plan,
		Version:       version,
		Resource:      resource,
		CloudProvider: cloudProvider,
		Region:        region,
		Status:        "DEPLOYING",
	}
	InstanceID = formattedInstance.ID

	var jsonData []string
	data, err := json.MarshalIndent(formattedInstance, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	jsonData = append(jsonData, string(data))

	// Print output
	switch output {
	case "text":
		err = utils.PrintText(jsonData)
		if err != nil {
			return err
		}
	case "table":
		err = utils.PrintTable(jsonData)
		if err != nil {
			return err
		}
	case "json":
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%+v\n", jsonData[0])
		if err != nil {
			utils.PrintError(err)
			return err
		}
	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
