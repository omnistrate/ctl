package testutils

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/omnistrate-oss/ctl/cmd"
	"github.com/omnistrate-oss/ctl/cmd/instance"
	"github.com/pkg/errors"
)

func WaitForInstanceToReachStatus(ctx context.Context, instanceID string, status instance.InstanceStatusType, timeout time.Duration) error {
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

		if currentStatus == instance.InstanceStatusFailed {
			ticker.Stop()
			return errors.New("instance deployment failed")
		}

		if currentStatus == instance.InstanceStatusCancelled {
			ticker.Stop()
			return errors.New("instance deployment cancelled")
		}
	}

	return errors.New("instance did not reach the expected status")
}
