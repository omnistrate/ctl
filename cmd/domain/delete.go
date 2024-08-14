package domain

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

var (
	deleteExample = `  # Delete domain with name
  omnistrate-ctl delete domain <name>

  # Delete multiple domains with names
  omnistrate-ctl delete domain <name1> <name2> <name3>`
)

var deleteCmd = &cobra.Command{
	Use:          "delete",
	Short:        "Delete one or more domains",
	Long:         `Delete domain by specifying name or environment type. Use --env to specify environment type. If not specified, name is assumed. If multiple domains are found with the same name, all of them will be deleted.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Args = cobra.MinimumNArgs(1) // Require at least one argument
}

func runDelete(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	environmentTypes := make([]string, 0)

	// List domains
	listRes, err := dataaccess.ListDomains(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Filter domains by name
	found := make(map[string]int)
	for _, name := range args {
		found[name] = 0
	}

	for _, d := range listRes.CustomDomains {
		if slices.Contains(args, d.Name) {
			environmentTypes = append(environmentTypes, string(d.EnvironmentType))
			found[d.Name] += 1
		}
	}

	domainsNotFound := make([]string, 0)
	for name, count := range found {
		if count == 0 {
			domainsNotFound = append(domainsNotFound, name)
		}
	}

	if len(domainsNotFound) > 0 {
		err = errors.New("domain(s) not found: " + strings.Join(domainsNotFound, ", "))
		utils.PrintError(err)
		return err
	}

	for name, count := range found {
		if count > 1 {
			utils.PrintWarning("Multiple domains found with name: " + name + ". Deleting all of them.")
		}
	}

	// Delete domain
	for _, environmentType := range environmentTypes {
		err = dataaccess.DeleteDomain(environmentType, token)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}

	utils.PrintSuccess("Domain(s) deleted successfully")

	return nil
}
