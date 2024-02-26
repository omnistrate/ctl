package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_login(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		Args           []string
		WantErr        bool
		ExpectedErrMsg string
	}{
		{[]string{"login"}, true, "must provide --email or -e"},
		{[]string{"login", "--email=xzhang+ctltest@omnistrate.com"}, true, "must provide a non-empty password via --password or --password-stdin"},
		{[]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=wrong_password"}, true, "unable to login, either email or password is incorrect"},
		{[]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"}, false, ""},
	}

	for _, tt := range tests {
		rootCmd.SetArgs(tt.Args)
		err := rootCmd.Execute()
		if tt.WantErr {
			require.Error(err)
			require.Contains(err.Error(), tt.ExpectedErrMsg)
		} else {
			require.NoError(err)
		}
	}
}
