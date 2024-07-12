package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New("")
	assert.Error(t, err)

	cfg, err := New("test.yml")
	assert.NoError(t, err)
	assert.Equal(t, "test.yml", cfg.FilePath)
	assert.Empty(t, cfg.AuthConfigs)
}

func TestConfigDir(t *testing.T) {
	assert.Equal(t, DefaultDir, ConfigDir())
}

func TestEnsureFile(t *testing.T) {
	filePath, err := EnsureFile()
	assert.NoError(t, err)
	assert.FileExists(t, filePath)

	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}

func TestFileExists(t *testing.T) {
	filePath, err := EnsureFile()
	assert.NoError(t, err)

	assert.True(t, fileExists())

	os.Remove(filePath)
	assert.False(t, fileExists())
}

func TestSaveAndLoad(t *testing.T) {
	dir := ConfigDir()

	cfg, err := New(filepath.Join(dir, DefaultFile))
	assert.NoError(t, err)

	cfg.AuthConfigs = append(cfg.AuthConfigs, AuthConfig{
		Email: "test@example.com",
		Token: "token123",
	})

	err = cfg.save()
	assert.NoError(t, err)

	loadedCfg, err := New(filepath.Join(dir, DefaultFile))
	assert.NoError(t, err)

	err = loadedCfg.load()
	assert.NoError(t, err)

	assert.Equal(t, cfg.AuthConfigs, loadedCfg.AuthConfigs)

	err = RemoveAuthConfig()
	assert.NoError(t, err)
}

func TestCreateOrUpdateAuthConfig(t *testing.T) {
	authConfig := AuthConfig{
		Email: "test@example.com",
		Token: "token123",
	}

	err := CreateOrUpdateAuthConfig(authConfig)
	assert.NoError(t, err)

	loadedConfig, err := LookupAuthConfig()
	assert.NoError(t, err)
	assert.Equal(t, authConfig, loadedConfig)

	err = RemoveAuthConfig()
	assert.NoError(t, err)
}

func TestLookupAuthConfig(t *testing.T) {
	_, err := LookupAuthConfig()
	assert.Error(t, err)

	authConfig := AuthConfig{
		Email: "test@example.com",
		Token: "token123",
	}

	err = CreateOrUpdateAuthConfig(authConfig)
	assert.NoError(t, err)

	loadedConfig, err := LookupAuthConfig()
	assert.NoError(t, err)
	assert.Equal(t, authConfig, loadedConfig)

	err = RemoveAuthConfig()
	assert.NoError(t, err)
}

func TestRemoveAuthConfig(t *testing.T) {
	authConfig := AuthConfig{
		Email: "test@example.com",
		Token: "token123",
	}

	err := CreateOrUpdateAuthConfig(authConfig)
	assert.NoError(t, err)

	err = RemoveAuthConfig()
	assert.NoError(t, err)

	_, err = LookupAuthConfig()
	assert.Error(t, err)
}

func TestLoadNonExistentFile(t *testing.T) {
	dir := ConfigDir()
	cfg, err := New(filepath.Join(dir, "non_existent.yml"))
	assert.NoError(t, err)

	err = cfg.load()
	assert.Error(t, err)
}

func TestLoadInvalidYaml(t *testing.T) {
	dir := ConfigDir()
	filePath := filepath.Join(dir, DefaultFile)
	err := os.WriteFile(filePath, []byte("invalid_yaml: [abc,"), 0600)
	assert.NoError(t, err)

	cfg, err := New(filePath)
	assert.NoError(t, err)

	err = cfg.load()
	assert.Error(t, err)
}
