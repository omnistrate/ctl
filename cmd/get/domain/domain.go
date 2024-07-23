package domain

import (
	"fmt"
	saasportalapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/saas_portal_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var (
	domainExample = `  # Get all domains
  omnistrate-ctl get domain

  # Get domain with name
  omnistrate-ctl get domain <name>

  # Get multiple domains
  omnistrate-ctl get domain <name1> <name2> <name3>`
)

// DomainCmd represents the describe command
var DomainCmd = &cobra.Command{
	Use:     "domain <name>",
	Short:   "Display one or more domains",
	Long:    `The get domain command displays basic information about one or more domains.`,
	Example: domainExample,
	RunE:    Run,
	PostRun: func(cmd *cobra.Command, args []string) {
		dataaccess.AskVerifyDomainIfAny()
	},
	SilenceUsage: true,
}

func Run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var domains []*saasportalapi.CustomDomain

	// List all domains
	var listRes *saasportalapi.ListSaaSPortalCustomDomainsResult
	listRes, err = dataaccess.ListDomains(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	allDomains := listRes.CustomDomains

	// Print domains table if no domain name is provided
	if len(args) == 0 {
		utils.PrintSuccess(fmt.Sprintf("%d domain(s) found", len(allDomains)))
		if len(allDomains) > 0 {
			printTable(allDomains)
		}
		return nil
	}

	// Format listRes.Domains into a map
	domainMap := make(map[string]*saasportalapi.CustomDomain)
	for _, domain := range allDomains {
		domainMap[domain.Name] = domain
	}

	// Filter domains by name
	for _, name := range args {
		domain, ok := domainMap[name]
		if !ok {
			utils.PrintError(fmt.Errorf("domain '%s' not found", name))
			continue
		}
		domains = append(domains, domain)
	}

	printTable(domains)

	return nil
}

func printTable(domains []*saasportalapi.CustomDomain) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Environment Type\tName\tDomain\tStatus\tCluster Endpoint")

	for _, domain := range domains {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			domain.EnvironmentType,
			domain.Name,
			domain.CustomDomain,
			domain.Status,
			domain.ClusterEndpoint)
	}

	err := w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
