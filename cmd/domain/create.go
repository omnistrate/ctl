package domain

import (
	"fmt"
	saasportalapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/saas_portal_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	createExample = `  # Create a custom domain for dev environment
  omnistrate-ctl domain create dev --domain abc.dev --env dev

  # Create a custom domain for prod environment
  omnistrate-ctl domain create abc.cloud --domain abc.cloud --env prod`
)

var createCmd = &cobra.Command{
	Use:          "create [flags]",
	Short:        "Create a domain",
	Long:         `Create a domain with the specified name and custom domain. The domain will be created for the specified environment type.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	createCmd.Flags().String("domain", "", "Custom domain")
	createCmd.Flags().String("environment-type", "", "Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private'")

	err := createCmd.MarkFlagRequired("domain")
	if err != nil {
		return
	}
	err = createCmd.MarkFlagRequired("environment-type")
	if err != nil {
		return
	}
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Get flags
	domain, _ := cmd.Flags().GetString("domain")
	environmentType, _ := cmd.Flags().GetString("environment-type")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	domains, err := dataaccess.ListDomains(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	for _, d := range domains.CustomDomains {
		// Check if domain with the same name already exists
		if d.Name == args[0] {
			err = errors.New("domain with the same name already exists, please choose a different name. You can use 'omnistrate-ctl get domain' to list all existing domains.")
			utils.PrintError(err)
			return err
		}

		// Check if domain is registered
		if d.CustomDomain == domain {
			err = errors.New(fmt.Sprintf("%s is already registered with another domain, please choose a different domain. You can use 'omnistrate-ctl get domain' to list all existing domains.", domain))
			utils.PrintError(err)
			return err
		}

		// Check if domain of the same environment type already exists
		if d.EnvironmentType == saasportalapi.EnvironmentType(strings.ToUpper(environmentType)) {
			err = errors.New("domain with the same environment type already exists, please choose a different environment type. You can use 'omnistrate-ctl get domain' to list all existing domains.")
			utils.PrintError(err)
			return err
		}
	}

	// Create domain
	request := &saasportalapi.CreateSaaSPortalCustomDomainRequest{
		Token:           token,
		Name:            args[0],
		Description:     "Custom domain for " + environmentType + " environment",
		EnvironmentType: saasportalapi.EnvironmentType(strings.ToUpper(environmentType)),
		CustomDomain:    domain,
	}

	err = dataaccess.CreateDomain(request)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess("Domain created successfully")

	domains, err = dataaccess.ListDomains(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var customDomain *saasportalapi.CustomDomain
	for _, d := range domains.CustomDomains {
		if d.Name == args[0] {
			customDomain = d
			break
		}
	}

	dataaccess.PrintNextStepVerifyDomainMsg(customDomain.ClusterEndpoint)

	return nil
}
