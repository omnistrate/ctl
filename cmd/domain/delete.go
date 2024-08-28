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
  omctl delete domain [domain-name]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [name]",
	Short:        "Delete a Custom Domain",
	Long:         `This command helps you delete a Custom Domain.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func runDelete(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

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
		err = dataaccess.DeleteDomain(token, environmentType)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}

	if output != "json" {
		utils.PrintSuccess("Domain deleted successfully")
	}

	return nil
}
