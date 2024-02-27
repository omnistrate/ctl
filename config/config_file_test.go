package config

import (
	"errors"
	"os"
	"regexp"
	"testing"
)

func Test_LookupAuthConfig_WithNoConfigFile(t *testing.T) {
	configDir, err := os.MkdirTemp("", "omnistrate-cli-file-test")
	if err != nil {
		t.Fatalf("can not create test config directory: %s", err)
	}
	defer os.RemoveAll(configDir)

	os.Setenv(ConfigLocationEnv, configDir)
	defer os.Unsetenv(ConfigLocationEnv)

	_, err = LookupAuthConfig()
	if err == nil {
		t.Errorf("Error was not returned")
	}

	if !errors.Is(err, ErrConfigNotFound) {
		t.Errorf("Error was not ErrConfigNotFound")
	}

	r := regexp.MustCompile(`(?m:config file not found)`)
	if !r.MatchString(err.Error()) {
		t.Errorf("Error not matched: %s", err.Error())
	}
}

func Test_UpdateAuthConfig_Insert(t *testing.T) {
	configDir, err := os.MkdirTemp("", "omnistrate-cli-file-test")
	if err != nil {
		t.Fatalf("can not create test config directory: %s", err)
	}
	defer os.RemoveAll(configDir)

	os.Setenv(ConfigLocationEnv, configDir)
	defer os.Unsetenv(ConfigLocationEnv)

	email := "test@abcd.com"
	token := "token"
	err = UpdateAuthConfig(AuthConfig{
		Email: email,
		Token: token,
		Auth:  JWTAuthType,
	})
	if err != nil {
		t.Fatalf("unexpected error when updating auth config: %s", err)
	}

	authConfig, err := LookupAuthConfig()
	if err != nil {
		t.Errorf("got error %s", err.Error())
		t.Errorf(authConfig.Token)
	}

	if authConfig.Email != email || authConfig.Token != token {
		t.Errorf("got email %s and token %s, expected %s %s", authConfig.Email, authConfig.Token, email, token)
	}
}

func Test_UpdateAuthConfig_Update(t *testing.T) {
	configDir, err := os.MkdirTemp("", "omnistrate-cli-file-test")
	if err != nil {
		t.Fatalf("can not create test config directory: %s", err)
	}
	defer os.RemoveAll(configDir)

	os.Setenv(ConfigLocationEnv, configDir)
	defer os.Unsetenv(ConfigLocationEnv)

	email := "test@abcd.com"
	token := "token"
	err = UpdateAuthConfig(AuthConfig{
		Email: email,
		Token: token,
		Auth:  JWTAuthType,
	})
	if err != nil {
		t.Fatalf("unexpected error when updating auth config: %s", err)
	}

	authConfig, err := LookupAuthConfig()
	if err != nil {
		t.Errorf("got error %s", err.Error())
	}

	if authConfig.Email != email || authConfig.Token != token {
		t.Errorf("got email %s and token %s, expected %s %s", authConfig.Email, authConfig.Token, email, token)
	}

	email = "test2@abcd.com"
	token = "token2"
	err = UpdateAuthConfig(AuthConfig{
		Email: email,
		Token: token,
		Auth:  JWTAuthType,
	})
	if err != nil {
		t.Fatalf("unexpected error when updating auth config: %s", err)
	}

	authConfig, err = LookupAuthConfig()
	if err != nil {
		t.Errorf("got error %s", err.Error())
	}

	if authConfig.Email != email || authConfig.Token != token {
		t.Errorf("got email %s and token %s, expected %s %s", authConfig.Email, authConfig.Token, email, token)
	}
}

func Test_New_NoFile(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Error("expected to fail on empty file path")
	}
}

func Test_EnsureFile(t *testing.T) {
	configDir, err := os.MkdirTemp("", "omnistrate-cli-file-test")
	if err != nil {
		t.Fatalf("can not create test config directory: %s", err)
	}
	defer os.RemoveAll(configDir)

	os.Setenv(ConfigLocationEnv, configDir)
	defer os.Unsetenv(ConfigLocationEnv)

	cfg, err := EnsureFile()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = os.Stat(cfg)
	if os.IsNotExist(err) {
		t.Errorf("expected config at %s", cfg)
	}
}

func Test_RemoveAuthConfig(t *testing.T) {
	configDir, err := os.MkdirTemp("", "omnistrate-cli-file-test")
	if err != nil {
		t.Fatalf("can not create test config directory: %s", err)
	}
	defer os.RemoveAll(configDir)

	os.Setenv(ConfigLocationEnv, configDir)
	defer os.Unsetenv(ConfigLocationEnv)

	email := "test@abcd.com"
	token := "token"
	err = UpdateAuthConfig(AuthConfig{
		Email: email,
		Token: token,
		Auth:  JWTAuthType,
	})
	if err != nil {
		t.Fatalf("unexpected error when updating auth config: %s", err)
	}

	err = RemoveAuthConfig()
	if err != nil {
		t.Errorf("got error %s", err.Error())
	}

	_, err = LookupAuthConfig()
	if err == nil {
		t.Fatal("Error was not returned")
	}
	r := regexp.MustCompile(`(?m:no auth config found)`)
	if !r.MatchString(err.Error()) {
		t.Errorf("Error not matched: %s", err.Error())
	}
}

func Test_RemoveAuthConfig_WithNoConfigFile(t *testing.T) {
	configDir, err := os.MkdirTemp("", "omnistrate-cli-file-test")
	if err != nil {
		t.Fatalf("can not create test config directory: %s", err)
	}
	defer os.RemoveAll(configDir)

	os.Setenv(ConfigLocationEnv, configDir)
	defer os.Unsetenv(ConfigLocationEnv)

	err = RemoveAuthConfig()
	if err == nil {
		t.Errorf("Error was not returned")
	}

	if !errors.Is(err, ErrConfigNotFound) {
		t.Errorf("Error was not ErrConfigNotFound")
	}

	r := regexp.MustCompile(`(?m:config file not found)`)
	if !r.MatchString(err.Error()) {
		t.Errorf("Error not matched: %s", err.Error())
	}
}

func Test_ConfigDir(t *testing.T) {

	cases := []struct {
		name         string
		env          map[string]string
		expectedPath string
	}{
		{
			name: "override value is returned",
			env: map[string]string{
				"OMNISTRATE_CONFIG": "/tmp/foo",
			},
			expectedPath: "/tmp/foo",
		},
		{
			name: "override value is returned, when CI is set but false",
			env: map[string]string{
				"OMNISTRATE_CONFIG": "/tmp/foo",
				"CI":                "false",
			},
			expectedPath: "/tmp/foo",
		},
		{
			name: "override value is returned even when CI is set",
			env: map[string]string{
				"OMNISTRATE_CONFIG": "/tmp/foo",
				"CI":                "true",
			},
			expectedPath: "/tmp/foo",
		},
		{
			name: "when CI is true, return the default CI directory",
			env: map[string]string{
				"CI": "true",
			},
			expectedPath: DefaultCIDir,
		},
		{
			name: "when CI is false, return the default directory",
			env: map[string]string{
				"CI": "false",
			},
			expectedPath: DefaultDir,
		},
		{
			name:         "when no other env variables are set, the default path is returned",
			expectedPath: DefaultDir,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			for name, value := range tc.env {
				os.Setenv(name, value)
				defer os.Unsetenv(name)
			}

			path := ConfigDir()
			if path != tc.expectedPath {
				t.Fatalf("expected config path '%s', got '%s'", tc.expectedPath, path)
			}
		})
	}

}
