package testutils

import (
	"context"
	"github.com/cenkalti/backoff/v4"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/instance"
	"github.com/pkg/errors"
	"time"
)

const (
	Running   = "RUNNING"
	Stopped   = "STOPPED"
	Failed    = "FAILED"
	Cancelled = "CANCELLED"
)

func WaitForInstanceToReachStatus(ctx context.Context, instanceID, status string, timeout time.Duration) error {
	b := &backoff.ConstantBackOff{
		Interval: 10 * time.Second,
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
