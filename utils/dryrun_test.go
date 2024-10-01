package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDryRun(t *testing.T) {
	t.Setenv(DryRunEnv, "true")
	assert.True(t, IsDryRun(), "DryRun should be true for tests")
}

func TestDryRunModify(t *testing.T) {
	t.Setenv(DryRunEnv, "false")
	assert.False(t, IsDryRun(), "DryRun should be false")
	t.Setenv(DryRunEnv, "true")
	assert.True(t, IsDryRun(), "DryRun should be true")
}
