package serviceplan

import (
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	"github.com/omnistrate/ctl/internal/utils"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterLatestNVersions(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		name         string
		servicePlans []openapiclientfleet.ServicePlanSearchRecord
		latestN      int
		expected     []openapiclientfleet.ServicePlanSearchRecord
	}{
		{
			name: "latestN is -1, return all service plans",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
			latestN: -1,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
		},
		{
			name: "latestN is greater than available plans, return all service plans",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
			latestN: 5,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
		},
		{
			name: "latestN is 1, return only the latest service plan",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
			latestN: 1,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
			},
		},
		{
			name: "latestN is 2, return the latest 2 service plans",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-06-01T00:00:00Z")},
			},
			latestN: 2,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
		},
		{
			name: "plans with nil release dates, sort properly",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
				{ReleasedAt: nil},
			},
			latestN: 2,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
		},
		{
			name: "mix of plans with and without release dates, return correct latest plans",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: nil},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
			latestN: 2,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: utils.ToPtr("2023-08-01T00:00:00Z")},
				{ReleasedAt: utils.ToPtr("2023-07-01T00:00:00Z")},
			},
		},
		{
			name: "all plans have nil release dates, return the first N plans",
			servicePlans: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: nil},
				{ReleasedAt: nil},
			},
			latestN: 1,
			expected: []openapiclientfleet.ServicePlanSearchRecord{
				{ReleasedAt: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterLatestNVersions(tt.servicePlans, tt.latestN)
			require.Equal(len(tt.expected), len(result))
			for i := range tt.expected {
				if tt.expected[i].ReleasedAt == nil {
					require.Nil(result[i].ReleasedAt)
				} else {
					require.NotNil(result[i].ReleasedAt)
					require.Equal(tt.expected[i].ReleasedAt, result[i].ReleasedAt)
				}
			}
		})
	}
}
