package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	composegenapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/compose_gen_api"
	"github.com/omnistrate/ctl/internal/config"
)

func CheckIfContainerImageAccessible(ctx context.Context, token string, request *composegenapi.CheckIfContainerImageAccessibleRequest) (res *composegenapi.CheckIfContainerImageAccessibleResult, err error) {
	request.Token = token

	service, err := httpclientwrapper.NewComposeGen(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	res, err = service.CheckIfContainerImageAccessible(ctx, request)
	if err != nil {
		return
	}
	return
}

func GenerateComposeSpecFromContainerImage(ctx context.Context, token string, request *composegenapi.GenerateComposeSpecFromContainerImageRequest) (res *composegenapi.GenerateComposeSpecFromContainerImageResult, err error) {
	request.Token = token

	service, err := httpclientwrapper.NewComposeGen(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	res, err = service.GenerateComposeSpecFromContainerImage(ctx, request)
	if err != nil {
		return
	}
	return
}
