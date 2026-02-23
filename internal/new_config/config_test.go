package new_config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/amadeusitgroup/cds/internal/cos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig_AllLocalFS(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, SourceTypeLocalFS, cfg.Inventory.Type)
	assert.Equal(t, SourceTypeLocalFS, cfg.Auth.Type)
	assert.Equal(t, SourceTypeLocalFS, cfg.Profile.Type)
	assert.Equal(t, SourceTypeLocalFS, cfg.Recipes.Type)

	assert.Contains(t, cfg.Inventory.Path, "inventory.json")
	assert.Contains(t, cfg.Auth.Path, "auth.json")
	assert.Contains(t, cfg.Profile.Path, "profile.json")
	assert.Contains(t, cfg.Recipes.Path, "recipes")
}

func TestSaveAndLoadRootConfig(t *testing.T) {
	setupMockFS(t)

	// Set a known config path for tests
	t.Setenv("CDS_CONFIG_PATH", "/tmp/testcds")

	cfg := &CdsConfig{
		Inventory: SourceRef{Type: SourceTypeLocalFS, Path: "/tmp/testcds/.xcds/inventory.json"},
		Auth:      SourceRef{Type: SourceTypeLocalFS, Path: "/tmp/testcds/.xcds/auth.json"},
		Profile:   SourceRef{Type: SourceTypeLocalFS, Path: "/tmp/testcds/.xcds/profile.json"},
		Recipes:   SourceRef{Type: SourceTypeLocalFS, Path: "/tmp/testcds/.xcds/recipes/"},
	}

	err := SaveRootConfig(cfg)
	require.NoError(t, err)

	loaded, err := LoadRootConfig()
	require.NoError(t, err)

	assert.Equal(t, cfg.Inventory, loaded.Inventory)
	assert.Equal(t, cfg.Auth, loaded.Auth)
	assert.Equal(t, cfg.Profile, loaded.Profile)
	assert.Equal(t, cfg.Recipes, loaded.Recipes)
}

func TestLoadRootConfig_FileNotFound(t *testing.T) {
	setupMockFS(t)
	t.Setenv("CDS_CONFIG_PATH", "/tmp/testcds")

	_, err := LoadRootConfig()
	assert.Error(t, err)
}

func TestLoadRootConfig_InvalidJSON(t *testing.T) {
	setupMockFS(t)
	t.Setenv("CDS_CONFIG_PATH", "/tmp/testcds")

	path := RootConfigPath()
	require.NoError(t, cos.Fs.MkdirAll("/tmp/testcds/.xcds", os.FileMode(0700)))
	require.NoError(t, cos.WriteFile(path, []byte("not json"), os.FileMode(0600)))

	_, err := LoadRootConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestSaveRootConfig_JSONFormat(t *testing.T) {
	setupMockFS(t)
	t.Setenv("CDS_CONFIG_PATH", "/tmp/testcds")

	cfg := &CdsConfig{
		Inventory: SourceRef{Type: SourceTypeLocalFS, Path: "/data/inventory.json"},
		Auth:      SourceRef{Type: SourceTypeLocalFS, Path: "/data/auth.json"},
		Profile:   SourceRef{Type: SourceTypeLocalFS, Path: "/data/profile.json"},
		Recipes:   SourceRef{Type: SourceTypeLocalFS, Path: "/data/recipes/"},
	}

	require.NoError(t, SaveRootConfig(cfg))

	data, err := cos.ReadFile(RootConfigPath())
	require.NoError(t, err)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &parsed))

	inv, ok := parsed["inventory"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "localfs", inv["type"])
	assert.Equal(t, "/data/inventory.json", inv["path"])
}

func TestLoadOrCreateRootConfig_CreatesDefault(t *testing.T) {
	setupMockFS(t)
	t.Setenv("CDS_CONFIG_PATH", "/tmp/testcds")

	cfg, err := LoadOrCreateRootConfig()
	require.NoError(t, err)
	assert.Equal(t, SourceTypeLocalFS, cfg.Inventory.Type)

	// Verify file was written
	exists := cos.Exists(RootConfigPath())
	assert.True(t, exists)
}

func TestLoadOrCreateRootConfig_LoadsExisting(t *testing.T) {
	setupMockFS(t)
	t.Setenv("CDS_CONFIG_PATH", "/tmp/testcds")

	original := &CdsConfig{
		Inventory: SourceRef{Type: SourceTypeLocalFS, Path: "/custom/inventory.json"},
		Auth:      SourceRef{Type: SourceTypeLocalFS, Path: "/custom/auth.json"},
		Profile:   SourceRef{Type: SourceTypeLocalFS, Path: "/custom/profile.json"},
		Recipes:   SourceRef{Type: SourceTypeLocalFS, Path: "/custom/recipes/"},
	}
	require.NoError(t, SaveRootConfig(original))

	loaded, err := LoadOrCreateRootConfig()
	require.NoError(t, err)
	assert.Equal(t, "/custom/inventory.json", loaded.Inventory.Path)
}
