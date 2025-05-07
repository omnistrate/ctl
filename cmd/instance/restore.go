package instance

import (
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	restoreExample = `# Restore to a new instance from a snapshot
omctl instance restore --service-id service-abc123 --environment-id env-xyz789 --snapshot-id snapshot-123def --param '{"key": "value"}'`
)

var restoreCmd = &cobra.Command{
	Use:          "restore --service-id <service-id> --environment-id <environment-id> --snapshot-id <snapshot-id> [--param=param] [--param-file=file-path]",
	Short:        "Create a new instance by restoring from a snapshot",
	Long:         `This command helps you create a new instance by restoring from a snapshot.`,
	Example:      restoreExample,
	RunE:         runRestore,
	SilenceUsage: true,
}

func init() {
	restoreCmd.Args = cobra.NoArgs
	restoreCmd.Flags().String("service-id", "", "The ID of the service")
	restoreCmd.Flags().String("environment-id", "", "The ID of the environment")
	restoreCmd.Flags().String("snapshot-id", "", "The ID of the snapshot to restore from")
	restoreCmd.Flags().String("param", "", "Parameters override for the instance deployment")
	restoreCmd.Flags().String("param-file", "", "Json file containing parameters override for the instance deployment")
	if err := restoreCmd.MarkFlagRequired("service-id"); err != nil {
		return
	}
	if err := restoreCmd.MarkFlagRequired("environment-id"); err != nil {
		return
	}
	if err := restoreCmd.MarkFlagRequired("snapshot-id"); err != nil {
		return
	}
	if err := restoreCmd.MarkFlagFilename("param-file"); err != nil {
		return
	}
}

func runRestore(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	serviceID, err := cmd.Flags().GetString("service-id")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	environmentID, err := cmd.Flags().GetString("environment-id")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	snapshotID, err := cmd.Flags().GetString("snapshot-id")
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

	// Validate user login
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Creating new instance from snapshot..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Format parameters
	formattedParams, err := common.FormatParams(param, paramFile)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Restore from snapshot
	result, err := dataaccess.RestoreResourceInstanceSnapshot(cmd.Context(), token, serviceID, environmentID, snapshotID, formattedParams)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully initiated restore operation from snapshot")

	// Print output
	if err = utils.PrintTextTableJsonOutput(output, result); err != nil {
		return err
	}

	return nil
}
