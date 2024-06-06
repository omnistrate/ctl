package cmd

import (
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
)

func Test_login(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	tests := []struct {
		Args           []string
		WantErr        bool
		ExpectedErrMsg string
	}{
		{[]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=Test@1234"}, false, ""},
		{[]string{"login"}, true, "must provide --email or -e"},
		{[]string{"login", "--email=xzhang+cli@omnistrate.com"}, true, "must provide a non-empty password via --password or --password-stdin"},
		{[]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=wrong_password"}, true, "wrong user email or password"},
	}

	for _, tt := range tests {
		rootCmd.SetArgs(tt.Args)
		err = rootCmd.Execute()
		if tt.WantErr {
			require.Error(err, tt.ExpectedErrMsg)
			require.Contains(err.Error(), tt.ExpectedErrMsg)
		} else {
			require.NoError(err)
		}
	}
}
