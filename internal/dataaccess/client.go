package dataaccess

import (
	"fmt"
	"net/http"

	openapiclientv1 "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/ctl/internal/config"
	openapiclientfleet "github.com/omnistrate/omnistrate-sdk-go/fleet"
	"github.com/pkg/errors"
)

// Configure registration api client
func getV1Client() *openapiclientv1.APIClient {
	configuration := openapiclientv1.NewConfiguration()
	configuration.Host = config.GetHost()
	configuration.Scheme = config.GetHostScheme()
	configuration.Debug = config.GetDebug()

	configuration.HTTPClient = &http.Client{
		Timeout: config.GetClientTimeout(),
	}

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

// Configure fleet api client
func getFleetClient() *openapiclientfleet.APIClient {
	configuration := openapiclientfleet.NewConfiguration()
	configuration.Host = config.GetHost()
	configuration.Scheme = config.GetHostScheme()
	configuration.Debug = config.GetDebug()

	configuration.HTTPClient = &http.Client{
		Timeout: config.GetClientTimeout(),
	}

	apiClient := openapiclientfleet.NewAPIClient(configuration)
	return apiClient
}

func handleFleetError(err error) error {
	if err != nil {
		var serviceErr *openapiclientfleet.GenericOpenAPIError
		ok := errors.As(err, &serviceErr)
		if !ok {
			return err
		}
		apiError, ok := serviceErr.Model().(openapiclientfleet.Error)
		if !ok {
			return fmt.Errorf("%s\nDetail: %s", serviceErr.Error(), string(serviceErr.Body()))
		}
		return fmt.Errorf("%s\nDetail: %s", apiError.Name, apiError.Message)
	}
	return err
}
