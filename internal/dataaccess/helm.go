package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	helmpackageapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/helm_package_api"
	"github.com/omnistrate/ctl/internal/config"
)

func SaveHelmChart(
	token string,
	chartName string,
	chartVersion string,
	namespace string,
	repoURL string,
	values map[string]any,
) (
	helmPackage *helmpackageapi.HelmPackage,
	err error,
) {
	helmPackageService := httpclientwrapper.NewHelmPackage(config.GetHostScheme(), config.GetHost())

	request := &helmpackageapi.SaveHelmPackageRequest{
		Token: token,
		HelmPackage: &helmpackageapi.HelmPackage{
			ChartName:    chartName,
			ChartVersion: chartVersion,
			Namespace:    namespace,
			RepoURL:      repoURL,
			Values:       values,
		},
	}

	if helmPackage, err = helmPackageService.SaveHelmPackage(context.Background(), request); err != nil {
		return
	}
	return
}

func ListHelmCharts(token string) (helmPackages *helmpackageapi.ListHelmPackagesResult, err error) {
	helmPackageService := httpclientwrapper.NewHelmPackage(config.GetHostScheme(), config.GetHost())

	request := &helmpackageapi.ListHelmPackagesRequest{
		Token: token,
	}

	if helmPackages, err = helmPackageService.ListHelmPackages(context.Background(), request); err != nil {
		return
	}
	return
}

func DescribeHelmChart(token, chartName, chartVersion string) (helmPackage *helmpackageapi.HelmPackage, err error) {
	helmPackageService := httpclientwrapper.NewHelmPackage(config.GetHostScheme(), config.GetHost())

	request := &helmpackageapi.DescribeHelmPackageRequest{
		Token:        token,
		ChartName:    chartName,
		ChartVersion: chartVersion,
	}

	if helmPackage, err = helmPackageService.DescribeHelmPackage(context.Background(), request); err != nil {
		return
	}
	return
}

func ListHelmChartInstallations(token string, hostClusterID *helmpackageapi.HostClusterID) (helmPackageInstallations *helmpackageapi.ListHelmPackageInstallationsResult, err error) {
	helmPackageService := httpclientwrapper.NewHelmPackage(config.GetHostScheme(), config.GetHost())

	request := &helmpackageapi.ListHelmPackageInstallationsRequest{
		Token:         token,
		HostClusterID: hostClusterID,
	}

	if helmPackageInstallations, err = helmPackageService.ListHelmPackageInstallations(context.Background(), request); err != nil {
		return
	}
	return
}

func DeleteHelmChart(token, chartName, chartVersion string) (err error) {
	helmPackageService := httpclientwrapper.NewHelmPackage(config.GetHostScheme(), config.GetHost())

	request := &helmpackageapi.DeleteHelmPackageRequest{
		Token:        token,
		ChartName:    chartName,
		ChartVersion: chartVersion,
	}

	if err = helmPackageService.DeleteHelmPackage(context.Background(), request); err != nil {
		return
	}
	return
}
