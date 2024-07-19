package dataaccess

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	saasportalapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/saas_portal_api"
	"github.com/omnistrate/ctl/utils"
)

func ListDomains(token string) (*saasportalapi.ListSaaSPortalCustomDomainsResult, error) {
	domain, err := httpclientwrapper.NewSaaSPortal(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := saasportalapi.ListSaaSPortalCustomDomainsRequest{
		Token: token,
	}

	res, err := domain.ListSaaSPortalCustomDomains(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteDomain(environmentType, token string) error {
	service, err := httpclientwrapper.NewSaaSPortal(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := saasportalapi.DeleteSaaSPortalCustomDomainRequest{
		Token:           token,
		EnvironmentType: saasportalapi.EnvironmentType(environmentType),
	}

	err = service.DeleteSaaSPortalCustomDomain(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}

func CreateDomain(request *saasportalapi.CreateSaaSPortalCustomDomainRequest) error {
	service, err := httpclientwrapper.NewSaaSPortal(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	err = service.CreateSaaSPortalCustomDomain(context.Background(), request)
	if err != nil {
		return err
	}

	return nil
}

const (
	DomainNotVerifiedWarningMsgTemplate = `
WARNING! Domain %s not verified. Need to verify ownership before use.

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

func AskVerifyDomainIfAny() {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return
	}

	// List all domains
	listRes, err := ListDomains(token)
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
