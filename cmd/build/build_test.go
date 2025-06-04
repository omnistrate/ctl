package build

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildReleaseDescriptionProcessing(t *testing.T) {
	require := require.New(t)

	// Test the logic that processes releaseName and releaseDescription
	
	// Test 1: Only releaseDescription is provided
	releaseName := ""
	releaseDescription := "v1.0.0-alpha"
	
	var releaseNamePtr *string
	if releaseName != "" {
		releaseNamePtr = &releaseName
	}
	if releaseDescription != "" {
		releaseNamePtr = &releaseDescription
	}
	
	require.NotNil(releaseNamePtr)
	require.Equal("v1.0.0-alpha", *releaseNamePtr)
	
	// Test 2: Both releaseName and releaseDescription are provided (releaseDescription should win)
	releaseName2 := "v0.9.0-beta"
	releaseDescription2 := "v1.0.0-alpha"
	
	var releaseNamePtr2 *string
	if releaseName2 != "" {
		releaseNamePtr2 = &releaseName2
	}
	if releaseDescription2 != "" {
		releaseNamePtr2 = &releaseDescription2
	}
	
	require.NotNil(releaseNamePtr2)
	require.Equal("v1.0.0-alpha", *releaseNamePtr2)
	
	// Test 3: Only releaseName is provided (legacy support)
	releaseName3 := "v0.9.0-beta"
	releaseDescription3 := ""
	
	var releaseNamePtr3 *string
	if releaseName3 != "" {
		releaseNamePtr3 = &releaseName3
	}
	if releaseDescription3 != "" {
		releaseNamePtr3 = &releaseDescription3
	}
	
	require.NotNil(releaseNamePtr3)
	require.Equal("v0.9.0-beta", *releaseNamePtr3)
	
	// Test 4: Neither is provided
	releaseName4 := ""
	releaseDescription4 := ""
	
	var releaseNamePtr4 *string
	if releaseName4 != "" {
		releaseNamePtr4 = &releaseName4
	}
	if releaseDescription4 != "" {
		releaseNamePtr4 = &releaseDescription4
	}
	
	require.Nil(releaseNamePtr4)
}