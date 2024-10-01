package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("test_key", "test_value")
	value := GetEnv("test_key", "default")
	assert.Equal(t, "test_value", value)
}

func TestGetEnvDefault(t *testing.T) {
	value := GetEnv("inexistent", "default")
	assert.Equal(t, "default", value)
}

func TestGetEnvAsBooleanTrue(t *testing.T) {
	os.Setenv("test_key", "true")
	value := GetEnvAsBoolean("test_key", "false")
	assert.True(t, value)
}

func TestGetEnvAsBooleanFalse(t *testing.T) {
	os.Setenv("test_key", "false")
	value := GetEnvAsBoolean("test_key", "true")
	assert.False(t, value)
}

func TestGetEnvAsBooleanDefault(t *testing.T) {
	os.Setenv("test_key", "")
	value := GetEnvAsBoolean("test_key", "true")
	assert.True(t, value)
}

func TestGetEnvAsBooleanInvalid(t *testing.T) {
	os.Setenv("test_key", "1")
	value := GetEnvAsBoolean("test_key", "true")
	assert.False(t, value)
}
