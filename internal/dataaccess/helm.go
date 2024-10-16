package dataaccess

import (
	"context"
	"net/http"

	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	helmpackageapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/helm_package_api"
	"github.com/omnistrate/ctl/internal/config"
)

func SaveHelmChart(
	ctx context.Context,
	token string,
	chartName string,
	chartVersion string,
	namespace string,
	repoName string,
	repoURL string,
	values map[string]any,
) (
	helmPackage *helmpackageapi.HelmPackage,
	err error,
) {

	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()

	r, err := apiClient.HelmPackageApiAPI.
		HelmPackageApiSaveHelmPackage(ctxWithToken).
		SaveHelmPackageRequestBody(openapiclient.SaveHelmPackageRequestBody{
			HelmPackage: openapiclient.HelmPackage{
				ChartName:     chartName,
				ChartVersion:  chartVersion,
				Namespace:     namespace,
				ChartRepoName: repoName,
				ChartRepoUrl:  repoURL,
				ChartValues:   values,
			},
		}).Execute()

	if err != nil {
		return nil, handleV1Error(err)
	}

	r.Body.Close()
	return
}

func ListHelmCharts(ctx context.Context, token string) (helmPackages *openapiclient.ListHelmPackagesResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()

	var r *http.Response
	helmPackages, r, err = apiClient.HelmPackageApiAPI.HelmPackageApiListHelmPackages(ctxWithToken).Execute()
	if err != nil {
		return nil, handleV1Error(err)
	}

	r.Body.Close()
	return
}

func DescribeHelmChart(ctx context.Context, token, chartName, chartVersion string) (helmPackage *openapiclient.HelmPackage, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()

	var r *http.Response
	helmPackage, r, err = apiClient.HelmPackageApiAPI.HelmPackageApiDescribeHelmPackage(ctxWithToken, chartName, chartVersion).Execute()
	if err != nil {
		return nil, handleV1Error(err)
	}

	r.Body.Close()
	return
}

func ListHelmChartInstallations(ctx context.Context, token string, hostClusterID *helmpackageapi.HostClusterID) (helmPackageInstallations *helmpackageapi.ListHelmPackageInstallationsResult, err error) {
	helmPackageService := httpclientwrapper.NewHelmPackage(config.GetHostScheme(), config.GetHost())

	request := &helmpackageapi.ListHelmPackageInstallationsRequest{
		Token:         token,
		HostClusterID: hostClusterID,
	}

	if helmPackageInstallations, err = helmPackageService.ListHelmPackageInstallations(ctx, request); err != nil {
		return
	}
	return
}

func DeleteHelmChart(ctx context.Context, token, chartName, chartVersion string) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	var r *http.Response
	r, err = apiClient.HelmPackageApiAPI.HelmPackageApiDeleteHelmPackage(ctxWithToken, chartName, chartVersion).Execute()
	if err != nil {
		return handleV1Error(err)
	}

	r.Body.Close()
	return
}
