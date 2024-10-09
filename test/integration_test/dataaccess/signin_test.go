package dataaccess

import (
	"testing"

	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/test/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignIn(t *testing.T) {
	testutils.IntegrationTest(t)

	testEmail, testPassword, err := testutils.GetTestAccount()
	require.NoError(t, err)

	tests := []struct {
		name           string
		email          string
		password       string
		wantErr        bool
		expectedErrMsg string
	}{
		{
			"valid login",
			testEmail,
			testPassword,
			false,
			"",
		},
		{
			"missing email",
			"",
			"",
			false,
			"",
		},
		{
			"missing password",
			"xzhang+cli@omnistrate.com",
			"",
			true,
			"must provide a non-empty password via --password or --password-stdin",
		},
		{
			"invalid password",
			"--email=xzhang+cli@omnistrate.com",
			"--password=wrong_password",
			true,
			"wrong user email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			token, err := dataaccess.LoginWithPassword("pberton@omnistrate.com", "invalidpassword")

			if tt.wantErr {
				assert.Equal(err.Error(), "bad_request\nDetail: "+tt.expectedErrMsg)
				assert.Empty(token)
			} else {
				require.NoError(err)
				assert.NotEmpty(token)
			}
		})
	}
}
