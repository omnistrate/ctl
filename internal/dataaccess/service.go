package dataaccess

import (
	"context"
	"fmt"

	openapiclient "github.com/omnistrate/omnistrate-sdk-go/v1"
)

const (
	NextStepsAfterBuildMsgTemplate = `
Next steps:
- Customize domain name for SaaS offer: check 'omnistrate-ctl create domain' command
- Update the service configuration: check 'omnistrate-ctl build' command`
)

func PrintNextStepsAfterBuildMsg() {
	fmt.Println(NextStepsAfterBuildMsgTemplate)
}

func ListServices(ctx context.Context, token string) (*openapiclient.ListServiceResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	resp, r, err := apiClient.ServiceApiAPI.ServiceApiListService(ctxWithToken).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return resp, nil
}

func DescribeService(ctx context.Context, token, serviceID string) (*openapiclient.DescribeServiceResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	resp, r, err := apiClient.ServiceApiAPI.ServiceApiDescribeService(ctxWithToken, serviceID).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return resp, nil
}

func DeleteService(ctx context.Context, token, serviceID string) error {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	r, err := apiClient.ServiceApiAPI.ServiceApiDeleteService(ctxWithToken, serviceID).Execute()

	err = handleV1Error(err)
	if err != nil {
		return err
	}
	r.Body.Close()

	return nil
}
