package instance

import (
	"fmt"
	"strings"
)

const (
	TerraformDeploymentType DeploymentType = "terraform"
)

type DeploymentType string

func getTerraformDeploymentName(resourceID, instanceID string) string {
	return strings.ToLower(fmt.Sprintf("tf-%s-%s", resourceID, instanceID))
}
