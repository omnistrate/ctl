package dataaccess

import (
	"context"
	"github.com/pkg/errors"

	"fmt"
	"io"
	"net/http"
)

func GetInstanceDeploymentEntity(ctx context.Context, token string, instanceID string, deploymentType string, deploymentName string) (output string, err error) {
	httpClient := getRetryableHttpClient()

	urlPath := fmt.Sprintf("http://localhost:80/2022-09-01-00/%s/%s/%s", deploymentType, instanceID, deploymentName)
	request, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return "", err
	}
	request = request.WithContext(ctx)
	request.Header.Add("Authorization", token)

	var response *http.Response
	defer func() {
		if response != nil {
			_ = response.Body.Close()
		}
	}()

	response, err = httpClient.Do(request)
	if err != nil {
		err = errors.Wrap(err, "failed to get instance deployment entity, you need to run it on dataplane agent")
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get instance deployment entity: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func PauseInstanceDeploymentEntity(ctx context.Context, token string, instanceID string, deploymentType string, deploymentName string) (err error) {
	httpClient := getRetryableHttpClient()

	urlPath := fmt.Sprintf("http://localhost:80/2022-09-01-00/%s/pause/%s/%s", deploymentType, instanceID, deploymentName)
	request, err := http.NewRequest(http.MethodPost, urlPath, nil)
	if err != nil {
		return err
	}
	request = request.WithContext(ctx)
	request.Header.Add("Authorization", token)

	var response *http.Response
	defer func() {
		if response != nil {
			_ = response.Body.Close()
		}
	}()
	
	response, err = httpClient.Do(request)
	if err != nil {
		err = errors.Wrap(err, "failed to pause instance deployment entity, you need to run it on dataplane agent")
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pause instance deployment entity: %s", response.Status)
	}

	return nil
}
