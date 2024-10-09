package dataaccess

import (
	"fmt"

	"github.com/omnistrate/ctl/config"
	openapiclientv1 "github.com/omnistrate/omnistrate-sdk-go/v1"
	"github.com/pkg/errors"
)

// Configure registration api client with retries
func getV1Client() *openapiclientv1.APIClient {
	configuration := openapiclientv1.NewConfiguration()
	configuration.Host = config.GetHost()
	configuration.Scheme = config.GetHostScheme()
	configuration.Debug = config.GetDebug()

	apiClient := openapiclientv1.NewAPIClient(configuration)
	return apiClient
}

func handleV1Error(err error) error {
	if err != nil {
		var serviceErr *openapiclientv1.GenericOpenAPIError
		ok := errors.As(err, &serviceErr)
		if !ok {
			return err
		}
		apiError, ok := serviceErr.Model().(openapiclientv1.Error)
		if !ok {
			return fmt.Errorf("%s\nDetail: %s", serviceErr.Error(), string(serviceErr.Body()))
		}
		return fmt.Errorf("%s\nDetail: %s", apiError.Name, apiError.Message)
	}
	return err
}
