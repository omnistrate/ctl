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

var (
	domainExample = `  # Create a custom domain for dev environment
  create domain dev --domain abc.dev --env dev

  # Create a custom domain for prod environment
  create domain abc.cloud --domain abc.cloud --env prod`
)

var DomainCmd = &cobra.Command{
	Use:          "domain <name> [flags]",
	Short:        "Create a domain",
	Long:         ``,
	Example:      domainExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	DomainCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	DomainCmd.Flags().String("domain", "", "Custom domain")
	DomainCmd.Flags().String("env", "", "Type of environment. Valid options include: 'prod', 'canary', 'staging', 'qa', 'dev'")

	err := DomainCmd.MarkFlagRequired("domain")
	if err != nil {
		return
	}
	err = DomainCmd.MarkFlagRequired("env")
	if err != nil {
		return
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Get flags
	domain, _ := cmd.Flags().GetString("domain")
	env, _ := cmd.Flags().GetString("env")

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
		if d.EnvironmentType == saasportalapi.EnvironmentType(strings.ToUpper(env)) {
			err = errors.New("domain with the same environment type already exists, please choose a different environment type. You can use 'omnistrate-ctl get domain' to list all existing domains.")
			utils.PrintError(err)
			return err
		}
	}

	// Create domain
	request := &saasportalapi.CreateSaaSPortalCustomDomainRequest{
		Token:           token,
		Name:            args[0],
		Description:     "Custom domain for " + env + " environment",
		EnvironmentType: saasportalapi.EnvironmentType(strings.ToUpper(env)),
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
