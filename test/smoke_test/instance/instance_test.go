package instance

import (
	"fmt"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/instance"
	"github.com/omnistrate/ctl/test/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	Running   = "RUNNING"
	Stopped   = "STOPPED"
	Failed    = "FAILED"
	Cancelled = "CANCELLED"
)

func TestInstanceBasic(t *testing.T) {
	testutils.SmokeTest(t)

	defer testutils.Cleanup()

	// PASS: login
	testEmail, testPassword, err := testutils.GetTestAccount()
	require.NoError(t, err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// PASS: create instance 1 with param
	cmd.RootCmd.SetArgs([]string{"instance", "create",
		"--service=mysql",
		"--environment=dev",
		"--plan=mysql",
		"--version=latest",
		"--resource=mySQL",
		"--cloud-provider=aws",
		"--region=ca-central-1",
		"--param", `{"databaseName":"default","password":"a_secure_password","rootPassword":"a_secure_root_password","username":"user"}`})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)
	instanceID1 := instance.InstanceID
	require.NotEmpty(t, instanceID1)

	// PASS: create instance 2 with param file
	cmd.RootCmd.SetArgs([]string{"instance", "create",
		"--service=mysql",
		"--environment=dev",
		"--plan=mysql",
		"--version=latest",
		"--resource=mySQL",
		"--cloud-provider=aws",
		"--region=ca-central-1",
		"--param-file", "paramfiles/instance_create_param.json"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)
	instanceID2 := instance.InstanceID
	require.NotEmpty(t, instanceID2)

	// PASS: describe instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "describe", instanceID1})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// PASS: describe instance 2
	cmd.RootCmd.SetArgs([]string{"instance", "describe", instanceID2})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	err = WaitForInstanceToReachStatus(instanceID1, Running, 300*time.Second)
	require.NoError(t, err)

	// PASS: stop instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "stop", instanceID1})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	err = WaitForInstanceToReachStatus(instanceID1, Stopped, 300*time.Second)
	require.NoError(t, err)

	// PASS: start instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "start", instanceID1})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	err = WaitForInstanceToReachStatus(instanceID1, Running, 300*time.Second)
	require.NoError(t, err)

	// PASS: restart instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "restart", instanceID1})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	err = WaitForInstanceToReachStatus(instanceID1, Running, 300*time.Second)
	require.NoError(t, err)

	// PASS: update instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "update", instanceID1, "--param", `{"databaseName":"default","password":"updated_password","rootPassword":"updated_root_password","username":"user"}`})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	err = WaitForInstanceToReachStatus(instanceID1, Running, 300*time.Second)
	require.NoError(t, err)

	// PASS: update instance 2
	cmd.RootCmd.SetArgs([]string{"instance", "update", instanceID2, "--param-file", "paramfiles/instance_update_param.json"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	err = WaitForInstanceToReachStatus(instanceID2, Running, 300*time.Second)
	require.NoError(t, err)

	// PASS: instance list
	cmd.RootCmd.SetArgs([]string{"instance", "list"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// PASS: instance list with filters
	cmd.RootCmd.SetArgs([]string{"instance", "list", "-f", "environment:DEV,cloud_provider:gcp", "-f", "environment:Dev,cloud_provider:aws"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// PASS: delete instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "delete", instanceID1, "--yes"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)

	// PASS: delete instance 2
	cmd.RootCmd.SetArgs([]string{"instance", "delete", instanceID2, "--yes"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(t, err)
}

func WaitForInstanceToReachStatus(instanceID, status string, timeout time.Duration) error {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     10 * time.Second,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         10 * time.Second,
		MaxElapsedTime:      timeout,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	b.Reset()
	ticker := backoff.NewTicker(b)

	for range ticker.C {
		cmd.RootCmd.SetArgs([]string{"instance", "describe", instanceID})
		err := cmd.RootCmd.ExecuteContext(ctx)
		if err != nil {
			return err
		}
		currentStatus := instance.InstanceStatus

		if currentStatus == status {
			ticker.Stop()
			return nil
		}

		if currentStatus == string(Failed) {
			ticker.Stop()
			return errors.New("instance deployment failed")
		}

		if currentStatus == string(Cancelled) {
			ticker.Stop()
			return errors.New("instance deployment cancelled")
		}
	}

	return errors.New("instance did not reach the expected status")
}
