package deprecated

import (
	"encoding/json"
	"fmt"

	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

var (
	describeServiceID string
)

var DescribeCmd = &cobra.Command{
	Use:          "describe [flags]",
	Short:        "Describe a Service (deprecated)",
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	// Deprecated flags. Kept for backwards compatibility.
	DescribeCmd.Flags().StringVarP(&describeServiceID, "service-id", "", "", "this flag is deprecated.")
}

func run(cmd *cobra.Command, args []string) (err error) {
	if describeServiceID == "" {
		err = cmd.Help()
		if err != nil {
			return
		}
	} else {
		defer func() {
			describeServiceID = ""
		}()

		// Validate user is currently logged in
		var token string
		token, err = config.GetToken()
		if err != nil {
			utils.PrintError(err)
			return
		}

		// Describe object
		var svc *openapiclient.DescribeServiceResult
		svc, err = dataaccess.DescribeService(cmd.Context(), token, describeServiceID)
		if err != nil {
			utils.PrintError(err)
			return
		}

		// Print service details
		var data []byte
		data, err = json.MarshalIndent(svc, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return
		}
		fmt.Println(string(data))
	}

	return
}
