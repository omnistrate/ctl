package dataaccess

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
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

func ListServices(token string) (*serviceapi.ListServiceResult, error) {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceapi.List{
		Token: token,
	}

	res, err := service.ListService(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DescribeService(token, serviceID string) (*serviceapi.DescribeServiceResult, error) {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceapi.DescribeServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceID),
	}

	res, err := service.DescribeService(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteService(token, serviceID string) error {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := serviceapi.DeleteServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceID),
	}

	err = service.DeleteService(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}
