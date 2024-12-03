package dataaccess

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	openapiclientv1 "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/pkg/errors"
)

// Configure registration api client
func getV1Client() *openapiclientv1.APIClient {
	configuration := openapiclientv1.NewConfiguration()
	configuration.Host = config.GetHost()
	configuration.Scheme = config.GetHostScheme()

	var servers openapiclientv1.ServerConfigurations
	for _, server := range configuration.Servers {
		server.URL = fmt.Sprintf("%s://%s", config.GetHostScheme(), config.GetHost())
		servers = append(servers, server)
	}
	configuration.Servers = servers

	configuration.HTTPClient = getRetryableHttpClient()

	configuration.Debug = config.IsDebugLogLevel()

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

	var servers openapiclientfleet.ServerConfigurations
	for _, server := range configuration.Servers {
		server.URL = fmt.Sprintf("%s://%s", config.GetHostScheme(), config.GetHost())
		servers = append(servers, server)
	}
	configuration.Servers = servers

	configuration.HTTPClient = getRetryableHttpClient()

	configuration.Debug = config.IsDebugLogLevel()

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

// Configure retryable http client
// retryablehttp gives us automatic retries with exponential backoff.
func getRetryableHttpClient() *http.Client {
	// retryablehttp gives us automatic retries with exponential backoff.
	httpClient := retryablehttp.NewClient()
	// HTTP requests are logged at DEBUG level.
	httpClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	httpClient.CheckRetry = retryablehttp.DefaultRetryPolicy
	httpClient.HTTPClient.Timeout = config.GetClientTimeout()
	httpClient.Logger
	return httpClient.StandardClient()
}
