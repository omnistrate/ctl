package deprecated

import (
	"encoding/json"
	"fmt"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	describeServiceID string
)

var DescribeCmd = &cobra.Command{
	Use:          "describe [flags]",
	Short:        "Describe a service (deprecated)",
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
		token, err = utils.GetToken()
		if err != nil {
			utils.PrintError(err)
			return
		}

		// Describe object
		var svc *serviceapi.DescribeServiceResult
		svc, err = dataaccess.DescribeService(token, describeServiceID)
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
