package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDryRun(t *testing.T) {
	assert.True(t, IsDryRun(), "DryRun should be true for tests")
}

func TestDryRunModify(t *testing.T) {
	os.Setenv(DryRunEnv, "false")
	assert.False(t, IsDryRun(), "DryRun should be false")
	os.Setenv(DryRunEnv, "true")
	assert.True(t, IsDryRun(), "DryRun should be true")
}
