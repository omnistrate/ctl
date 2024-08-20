package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTruncateString(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name     string
		input    string
		max      int
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			max:      10,
			expected: "",
		},
		{
			name:     "Short string, no truncation needed",
			input:    "Short",
			max:      10,
			expected: "Short",
		},
		{
			name:     "Exact length, no truncation needed",
			input:    "Exact length",
			max:      12,
			expected: "Exact length",
		},
		{
			name:     "Truncation at word boundary",
			input:    "This is a long sentence that should be truncated.",
			max:      25,
			expected: "This is a long sentence...",
		},
		{
			name:     "Truncation of a long word without spaces",
			input:    "Thisisaverylongwordwithoutspaces",
			max:      10,
			expected: "Thisisaver...",
		},
		{
			name:     "No truncation needed, sentence ends with punctuation",
			input:    "This sentence ends with punctuation!",
			max:      36,
			expected: "This sentence ends with punctuation!",
		},
		{
			name:     "Truncation with punctuation",
			input:    "Another sentence; with punctuation.",
			max:      25,
			expected: "Another sentence; with...",
		},
		{
			name:     "Truncation with trailing spaces and punctuation",
			input:    "Trailing spaces and punctuations; ",
			max:      30,
			expected: "Trailing spaces and...",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TruncateString(test.input, test.max)
			require.Equal(test.expected, got, "TruncateString(%q, %d) = %q; want %q", test.input, test.max, got, test.expected)
		})
	}
}
