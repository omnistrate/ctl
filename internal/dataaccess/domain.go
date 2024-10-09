package dataaccess

import (
	"context"
	"fmt"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	saasportalapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/saas_portal_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/utils"
)

func ListDomains(ctx context.Context, token string) (*saasportalapi.ListSaaSPortalCustomDomainsResult, error) {
	domain, err := httpclientwrapper.NewSaaSPortal(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	request := saasportalapi.ListSaaSPortalCustomDomainsRequest{
		Token: token,
	}

	res, err := domain.ListSaaSPortalCustomDomains(ctx, &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteDomain(ctx context.Context, token, environmentType string) error {
	service, err := httpclientwrapper.NewSaaSPortal(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	request := saasportalapi.DeleteSaaSPortalCustomDomainRequest{
		Token:           token,
		EnvironmentType: saasportalapi.EnvironmentType(environmentType),
	}

	err = service.DeleteSaaSPortalCustomDomain(ctx, &request)
	if err != nil {
		return err
	}
	return nil
}

func CreateDomain(ctx context.Context, request *saasportalapi.CreateSaaSPortalCustomDomainRequest) error {
	service, err := httpclientwrapper.NewSaaSPortal(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	err = service.CreateSaaSPortalCustomDomain(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

const (
	DomainNotVerifiedWarningMsgTemplate = `
WARNING! Domain %s is not verified. Need to verify ownership before use.

Please create the following DNS section for your domain and add two CNAME records as follows.

- Type: CNAME
  Name: @
  Target: %s

- Type: CNAME
  Name: www
  Target: %s
`

	NextStepVerifyDomainMsgTemplate = `
Next step:

Verify domain ownership.

Please create the following DNS section for your domain and add two CNAME records as follows.

- Type: CNAME
  Name: @
  Target: %s

- Type: CNAME
  Name: www
  Target: %s
`
)

func PrintNextStepVerifyDomainMsg(clusterEndpoint string) {
	fmt.Println(fmt.Sprintf(NextStepVerifyDomainMsgTemplate, clusterEndpoint, clusterEndpoint))
}

func PrintDomainNotVerifiedWarningMsg(domain, clusterEndpoint string) {
	utils.PrintWarning(fmt.Sprintf(DomainNotVerifiedWarningMsgTemplate, domain, clusterEndpoint, clusterEndpoint))
}

func AskVerifyDomainIfAny(ctx context.Context) {
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return
	}

	// List all domains
	listRes, err := ListDomains(ctx, token)
	if err != nil {
		utils.PrintError(err)
		return
	}

	// Warn if any domains are not verified
	for _, domain := range listRes.CustomDomains {
		if domain.Status == "PENDING" {
			PrintDomainNotVerifiedWarningMsg(domain.CustomDomain, domain.ClusterEndpoint)
		}
	}
}
