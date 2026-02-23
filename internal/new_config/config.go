package new_config

import (
	"encoding/json"
	"fmt"

	"github.com/amadeusitgroup/cds/internal/cenv"
	"github.com/amadeusitgroup/cds/internal/cos"
	cg "github.com/amadeusitgroup/cds/internal/global"
)

const (
	rootConfigFile  = "cds.json"
	inventoryFile   = "inventory.json"
	authFile        = "auth.json"
	profileFile     = "profile.json"
	recipesDir      = "recipes"
)

// CdsConfig is the root configuration for CDS.
// It holds SourceRefs pointing to the 4 main config files.
type CdsConfig struct {
	Inventory SourceRef `json:"inventory"`
	Auth      SourceRef `json:"auth"`
	Profile   SourceRef `json:"profile"`
	Recipes   SourceRef `json:"recipes"`
}

// RootConfigPath returns the absolute path to the root config file.
func RootConfigPath() string {
	return cenv.ConfigFile(rootConfigFile)
}

// LoadRootConfig reads and parses the root CDS config from ~/.xcds/cds.json.
func LoadRootConfig() (*CdsConfig, error) {
	path := RootConfigPath()

	data, err := cos.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read root config at %s: %w", path, err)
	}

	var cfg CdsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse root config at %s: %w", path, err)
	}

	return &cfg, nil
}

// SaveRootConfig writes the root CDS config to ~/.xcds/cds.json.
func SaveRootConfig(cfg *CdsConfig) error {
	path := RootConfigPath()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal root config: %w", err)
	}

	if err := cenv.EnsureFile(path, cg.KPermFile); err != nil {
		return fmt.Errorf("failed to ensure root config file at %s: %w", path, err)
	}

	if err := cos.WriteFile(path, data, cg.KPermFile); err != nil {
		return fmt.Errorf("failed to write root config at %s: %w", path, err)
	}

	return nil
}

// DefaultConfig returns a CdsConfig with all 4 sources pointing to
// localfs files under the default ~/.xcds/ directory.
func DefaultConfig() *CdsConfig {
	return &CdsConfig{
		Inventory: SourceRef{
			Type: SourceTypeLocalFS,
			Path: cenv.ConfigFile(inventoryFile),
		},
		Auth: SourceRef{
			Type: SourceTypeLocalFS,
			Path: cenv.ConfigFile(authFile),
		},
		Profile: SourceRef{
			Type: SourceTypeLocalFS,
			Path: cenv.ConfigFile(profileFile),
		},
		Recipes: SourceRef{
			Type: SourceTypeLocalFS,
			Path: cenv.ConfigDir(recipesDir),
		},
	}
}

// LoadOrCreateRootConfig attempts to load the root config.
// If it does not exist, it creates a default config, saves it, and returns it.
func LoadOrCreateRootConfig() (*CdsConfig, error) {
	path := RootConfigPath()

	if cos.Exists(path) {
		return LoadRootConfig()
	}

	cfg := DefaultConfig()
	if err := SaveRootConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to create default root config: %w", err)
	}

	return cfg, nil
}
