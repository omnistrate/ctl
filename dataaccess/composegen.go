package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	composegenapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/compose_gen_api"
	"github.com/omnistrate/ctl/utils"
)

func CheckIfContainerImageAccessible(token string, request *composegenapi.CheckIfContainerImageAccessibleRequest) (res *composegenapi.CheckIfContainerImageAccessibleResult, err error) {
	request.Token = token

	service, err := httpclientwrapper.NewComposeGen(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	res, err = service.CheckIfContainerImageAccessible(context.Background(), request)
	if err != nil {
		return
	}
	return
}

func GenerateComposeSpecFromContainerImage(token string, request *composegenapi.GenerateComposeSpecFromContainerImageRequest) (res *composegenapi.GenerateComposeSpecFromContainerImageResult, err error) {
	request.Token = token

	service, err := httpclientwrapper.NewComposeGen(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	res, err = service.GenerateComposeSpecFromContainerImage(context.Background(), request)
	if err != nil {
		return
	}
	return
}
