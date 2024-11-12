package service

import (
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	initExample = `# Initialize service with a template
omctl service init [service-name] --template=[template-name] --model=[model-name] --vectordb=[vectordb]
`
)

var initCmd = &cobra.Command{
	Use:          "init [service-name] [flags]",
	Short:        "Initialize a service with a template",
	Long:         `This command helps you initialize a service with a SaaS template.`,
	Example:      initExample,
	RunE:         runInit,
	SilenceUsage: true,
}

func init() {
	initCmd.Flags().StringP("template", "", "", "Template name. Valid options include: 'ai-chatbot', 'ai-image-generator'")
	initCmd.Flags().StringP("model", "", "openai-gpt-4o", "Model name. Valid options include: 'openai-gpt-4o'")
	initCmd.Flags().StringP("vectordb", "", "weaviate", "VectorDB name. Valid options include: 'weaviate'")

	err := initCmd.MarkFlagRequired("template")
	if err != nil {
		return
	}

	initCmd.Args = cobra.ExactArgs(1)
}

func runInit(cmd *cobra.Command, args []string) (err error) {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Get flags
	template, _ := cmd.Flags().GetString("template")
	model, _ := cmd.Flags().GetString("model")
	output, _ := cmd.Flags().GetString("output")
	vectordb, _ := cmd.Flags().GetString("vectordb")
	serviceName := args[0]

	// Validate user is logged in
	var token string
	token, err = config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		spinner = sm.AddSpinner("Initializing project...")
		sm.Start()
	}

	// Check if a spec file already exists in the current directory
	if err = utils.CheckSpecFileExists(spinner, sm, "omnistrate-spec.yaml"); err != nil {
		return
	}

	// Initialize service from template

}
