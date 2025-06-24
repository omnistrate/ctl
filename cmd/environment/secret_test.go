package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEnvironmentType(t *testing.T) {
	// Test valid environment types
	validTypes := []string{"dev", "qa", "staging", "canary", "prod", "private"}
	for _, envType := range validTypes {
		err := validateEnvironmentType(envType)
		assert.NoError(t, err, "Expected %s to be valid", envType)
	}

	// Test invalid environment types
	invalidTypes := []string{"invalid", "test", "demo", ""}
	for _, envType := range invalidTypes {
		err := validateEnvironmentType(envType)
		assert.Error(t, err, "Expected %s to be invalid", envType)
	}
}

func TestSecretCommands(t *testing.T) {
	// Test that secret commands are properly initialized
	assert.NotNil(t, secretCmd)
	assert.Equal(t, "secret", secretCmd.Name())
	assert.Equal(t, "Manage environment secrets", secretCmd.Short)

	// Test that all secret subcommands are added
	subcommands := secretCmd.Commands()
	expectedCommands := []string{"create", "delete", "describe", "list", "update"}
	
	assert.Len(t, subcommands, len(expectedCommands))
	
	commandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		commandNames[cmd.Name()] = true
	}
	
	for _, expectedCmd := range expectedCommands {
		assert.True(t, commandNames[expectedCmd], "Expected command %s to be present", expectedCmd)
	}
}