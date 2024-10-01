package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestTruncateString(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	str1 := "A-tisket a-tasket A green and yellow basket"
	str2 := "Peter Piper picked a peck of pickled peppers"

	assert.Equal("A-tisket...", TruncateString(str1, 11))
	assert.Equal("Peter Piper...", TruncateString(str2, 14))
	assert.Equal("A-tisket a-tasket A green and yellow basket", TruncateString(str1, len(str1)))
	assert.Equal("A-tisket a-tasket A green and yellow basket", TruncateString(str1, len(str1)+2))
	assert.Equal("A...", TruncateString("A-", 1))
	assert.Equal("Ab...", TruncateString("Absolutely Longer", 2))
}

func TestTruncateStringAndMax(t *testing.T) {
	t.Parallel()

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
			max:      26,
			expected: "This is a long sentence...",
		},
		{
			name:     "Truncation of a long word without spaces",
			input:    "Thisisaverylongwordwithoutspaces",
			max:      10,
			expected: "Thisisa...",
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
			require.LessOrEqual(len(test.expected), test.max)
			got := TruncateString(test.input, test.max)
			require.LessOrEqual(len(got), test.max)
			require.Equal(test.expected, got, "TruncateString(%q, %d) = %q; want %q", test.input, test.max, got, test.expected)
		})
	}
}

func TestCutString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		str    string
		length int
		want   string
	}{
		// Basic cases
		{"empty string", "", 5, ""},
		{"truncate to zero", "hello", 0, ""},
		{"negative length", "hello", -1, ""},
		{"no truncation needed", "hello", 5, "hello"},
		{"truncate to less", "hello", 3, "hel"},
		{"longer than string", "hi", 10, "hi"},

		// Unicode and special characters
		{"unicode characters", "こんにちは", 3, "こんに"},
		{"mixed ascii and unicode", "hello😊world", 8, "hello😊wo"},

		// Edge cases
		{"exact length", "hey", 3, "hey"},
		{"single character", "a", 1, "a"},
		{"truncate to one", "world", 1, "w"},
		{"empty string, zero length", "", 0, ""},
		{"empty string, negative length", "", -1, ""},
		{"non-zero string, zero length", "hello", 0, ""},
		{"non-zero string, negative length", "hello", -1, ""},
		{"service name", Keyify("service-name-omnitest-motorhead-simple-2024-04-15 00:44:00-cb750f1e-29a0-433a-a8b9-ea9f5eaeb6f1"), 60, "service-name-omnitest-motorhead-simple-2024-04-15-00-44-00-c"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := CutString(tc.str, tc.length); got != tc.want {
				t.Errorf("CutString(%q, %d) = %q; want %q", tc.str, tc.length, got, tc.want)
			}
		})
	}
}

func TestIsValidDirectoryPath(t *testing.T) {
	t.Parallel()

	// These are the test cases that we will use
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/usr/local/bin", true},
		{"/", true},
		{"/var/log/system.log", true},
		{"/etc/nginx/conf.d", true},
		{"//example.com/path/to/resource", true},
		{"/etc/passwd/test", true},
		{"/var/log/system.log/test", true},

		{"./test", false},
		{"../test", false},
		{"C:\\Windows\\System32", false},
		{"C:/Windows/System32", false},
		{"~/Documents/Files", false},
		{"./test/../test", false},
		{"http://example.com", false},
		{"ftp://example.com", false},
		{"file:///etc/passwd", false},
		{"mailto:user@example.com", false},
		{"data:image/png;base64,iVBORw0KGg....", false},
		{"./test/../test/../test", false},
		{"../test/../test/../test", false},
	}

	// Iterate over the test cases and test the function
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.path), func(t *testing.T) {
			err := ValidateIfIsAValidUnixPath(tc.path)

			result := err == nil

			if result != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, result)
			}
		})
	}
}

func TestKeyify(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	assert.Equal("test", Keyify("test"))
	assert.Equal("test", Keyify("TEST"))
	assert.Equal("test-1234", Keyify("test_1234"))
	assert.Equal("test-1234", Keyify("test 1234"))
	assert.Equal("test-1234", Keyify("TEST 1234"))
}

func TestRemoveDashes(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	assert.Equal("test", RemoveDashes("test"))
	assert.Equal("test", RemoveDashes("TEST"))
	assert.Equal("test1234", RemoveDashes("test-1234"))
	assert.Equal("test 1234", RemoveDashes("test 1234"))
	assert.Equal("test 1234", RemoveDashes("TEST 1234"))
}

func TestValidateURL(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	// PASS: Test for valid secure URL
	err := ValidateURL("https://www.google.com")
	assert.NoError(err)

	// PASS: Test for valid insecure URL
	err = ValidateURL("http://www.google.com")
	assert.NoError(err)

	// PASS: Test for valid URL with port
	err = ValidateURL("http://www.google.com:8080")
	assert.NoError(err)

	// PASS: Test for valid URL with path
	err = ValidateURL("http://www.google.com:8080/path")
	assert.NoError(err)

	// PASS: Test for valid URL with query
	err = ValidateURL("http://www.google.com:8080/path?query=1")
	assert.NoError(err)

	// PASS: Test for valid URL with fragment
	err = ValidateURL("http://www.google.com:8080/path?query=1#fragment")
	assert.NoError(err)

	// PASS: Test for valid URL with username and password
	err = ValidateURL("http://user:password@google.com:8080/path?query=1#fragment")
	assert.NoError(err)

	// PASS: Test for valid URL without scheme
	err = ValidateURL("www.google.com:8080/path?query=1#fragment")
	assert.NoError(err)

	// PASS: Test for valid URL without scheme and port
	err = ValidateURL("www.google.com/path?query=1#fragment")
	assert.NoError(err)

	// FAIL: Test for invalid URL with empty scheme
	err = ValidateURL("://www.google.com")
	assert.Error(err)

	// FAIL: Test for invalid URL with non-numeric port
	err = ValidateURL("http://www.google.com:port/path?query=1#fragment")
	assert.Error(err)
}

func TestToLowerCamelCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"All lowercase", "hello", "hello"},
		{"All uppercase", "HELLO", "hello"},
		{"Mixed case", "Hello", "hello"},
		{"Already camelCase", "helloWorld", "helloWorld"},
		{"Underscores", "hello_world", "helloWorld"},
		{"Hyphens", "hello-world", "hello-world"},
		{"Spaces", "hello world", "helloWorld"},
		{"Punctuation", "hello, world! How are!you?", "helloWorldHowAreYou"},
		{"Consecutive uppercase", "HTTPServer", "hTTPServer"},
		{"Acronym in sentence", "my CPU is overheating", "myCpuIsOverheating"},
		{"Starts with numbers", "123startHere", "123startHere"},
		{"Only non-letters", "123456!@#$%^", "123456"},
		{"OMNISTRATE_HOSTED", "OMNISTRATE_HOSTED", "omnistrateHosted"},
		{"CUSTOMER_HOSTED", "CUSTOMER_HOSTED", "customerHosted"},
		{"dedicated tenancy", "dedicated tenancy", "dedicatedTenancy"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToLowerCamelCase(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"All lowercase", "hello", "Hello"},
		{"All uppercase", "HELLO", "Hello"},
		{"Mixed case", "Hello", "Hello"},
		{"Already camelCase", "helloWorld", "HelloWorld"},
		{"Underscores", "hello_world", "HelloWorld"},
		{"Hyphens", "hello-world", "Hello-world"},
		{"Spaces", "hello world", "HelloWorld"},
		{"Punctuation", "hello, world! How are_you?", "HelloWorldHowAreYou"},
		{"Consecutive uppercase", "HTTPServer", "HTTPServer"},
		{"Acronym in sentence", "my CPU is overheating", "MyCpuIsOverheating"},
		{"Starts with numbers", "123startHere", "123startHere"},
		{"Only non-letters", "123456!@#$%^", "123456"},
		{"OMNISTRATE_HOSTED", "OMNISTRATE_HOSTED", "OmnistrateHosted"},
		{"CUSTOMER_HOSTED", "CUSTOMER_HOSTED", "CustomerHosted"},
		{"dedicated tenancy", "dedicated tenancy", "DedicatedTenancy"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToCamelCase(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestValidateKey(t *testing.T) {
	t.Parallel()

	// Valid
	require.Nil(t, ValidateKey("key_test"))
	require.Nil(t, ValidateKey("key1"))

	// Invalid
	require.NotNil(t, ValidateKey("1key"))
	require.NotNil(t, ValidateKey("_key"))
	require.NotNil(t, ValidateKey("key&"))
}

func TestConvertSnakeToLowerCamelCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"All lowercase", "hello", "hello"},
		{"All uppercase", "HELLO", "hello"},
		{"Mixed case", "Hello", "hello"},
		{"Underscores", "hello_world", "helloWorld"},
		{"Underscores with numbers", "hello_world_123", "helloWorld123"},
		{"Underscores with special characters", "hello_world_!@#", "helloWorld"},
		{"Starts with numbers", "123_start_here", "123StartHere"},
		{"Only non-letters", "123456!@#$%^", "123456"},
		{"OMNISTRATE_HOSTED", "OMNISTRATE_HOSTED", "omnistrateHosted"},
		{"CUSTOMER_HOSTED", "CUSTOMER_HOSTED", "customerHosted"},
		{"dedicated tenancy", "dedicated_tenancy", "dedicatedTenancy"},
		{"dedicated tenancy", "dedicated_tenancy_123", "dedicatedTenancy123"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertSnakeToLowerCamelCase(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestCheckIfNonAlphanumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", false},
		{"Alphanumeric", "hello", false},
		{"Alphanumeric with numbers", "hello 123 ą", false},
		{"Alphanumeric with special characters", "hello!@#", true},
		{"Special characters", "_//+", true},
		{"Numbers", "123", false},
		{"Numbers with special characters", "123!@#", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckIfNonAlphanumeric(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %t, got %t", tc.expected, result)
			}
		})
	}
}
