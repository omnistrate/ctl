package serviceplan

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestGetReleaseDescription(t *testing.T) {
	require := require.New(t)

	// Test empty description returns nil
	result := getReleaseDescription("")
	require.Nil(result)

	// Test non-empty description returns pointer to string
	description := "v1.0.0-alpha"
	result = getReleaseDescription(description)
	require.NotNil(result)
	require.Equal(description, *result)

	// Test another description
	description2 := "Release with custom description"
	result2 := getReleaseDescription(description2)
	require.NotNil(result2)
	require.Equal(description2, *result2)
}

func TestValidateReleaseArguments(t *testing.T) {
	require := require.New(t)

	// Test valid arguments
	err := validateReleaseArguments([]string{"service", "plan"}, "", "")
	require.NoError(err)

	// Test valid service and plan IDs
	err = validateReleaseArguments([]string{}, "service-id", "plan-id")
	require.NoError(err)

	// Test invalid - missing arguments and IDs
	err = validateReleaseArguments([]string{}, "", "")
	require.Error(err)

	// Test invalid - wrong number of arguments
	err = validateReleaseArguments([]string{"service"}, "", "")
	require.Error(err)

	// Test invalid - too many arguments
	err = validateReleaseArguments([]string{"service", "plan", "extra"}, "", "")
	require.Error(err)
}

func TestReleaseCmdFlags(t *testing.T) {
	require := require.New(t)

	// Test that the release command has the correct flags defined
	releaseDescFlag := releaseCmd.Flags().Lookup("release-description")
	require.NotNil(releaseDescFlag)
	require.Equal("string", releaseDescFlag.Value.Type())

	releaseAsPreferredFlag := releaseCmd.Flags().Lookup("release-as-preferred")
	require.NotNil(releaseAsPreferredFlag)
	require.Equal("bool", releaseAsPreferredFlag.Value.Type())

	dryrunFlag := releaseCmd.Flags().Lookup("dryrun")
	require.NotNil(dryrunFlag)
	require.Equal("bool", dryrunFlag.Value.Type())
}

func TestReleaseCmdFlagParsing(t *testing.T) {
	require := require.New(t)

	// Create a test command
	cmd := &cobra.Command{}
	cmd.Flags().String("release-description", "", "Test flag")

	// Set the flag value
	err := cmd.Flags().Set("release-description", "v1.0.0-alpha")
	require.NoError(err)

	// Retrieve the flag value
	value, err := cmd.Flags().GetString("release-description")
	require.NoError(err)
	require.Equal("v1.0.0-alpha", value)

	// Test getReleaseDescription with the parsed value
	result := getReleaseDescription(value)
	require.NotNil(result)
	require.Equal("v1.0.0-alpha", *result)
}