package new_profile

import (
	"encoding/json"
	"fmt"
	"sync"

	nfg "github.com/amadeusitgroup/cds/internal/new_config"
)

// Profile is a stub for the CDS profile specification.
// The Holder field is a placeholder — replace with real spec fields later.
type Profile struct {
	Holder string `json:"holder"`
}

// ProfileManager manages a profile file backed by a Source.
type ProfileManager struct {
	source nfg.Source
	ref    nfg.SourceRef
	data   Profile
	mu     sync.Mutex
}

// NewProfileManager creates a ProfileManager that will load/save
// the profile through the given Source at ref.Path.
func NewProfileManager(src nfg.Source, ref nfg.SourceRef) *ProfileManager {
	return &ProfileManager{
		source: src,
		ref:    ref,
	}
}

// Load reads and parses the profile file from the Source.
// If the file does not exist, the profile starts empty.
func (pm *ProfileManager) Load() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	exists, err := pm.source.Exists(pm.ref.Path)
	if err != nil {
		return fmt.Errorf("failed to check profile file existence: %w", err)
	}
	if !exists {
		pm.data = Profile{}
		return nil
	}

	raw, err := pm.source.Read(pm.ref.Path)
	if err != nil {
		return fmt.Errorf("failed to read profile file at %s: %w", pm.ref.Path, err)
	}

	var p Profile
	if err := json.Unmarshal(raw, &p); err != nil {
		return fmt.Errorf("failed to parse profile file at %s: %w", pm.ref.Path, err)
	}

	pm.data = p
	return nil
}

// Save writes the current profile data to the Source.
func (pm *ProfileManager) Save() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	raw, err := json.MarshalIndent(pm.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := pm.source.Write(pm.ref.Path, raw); err != nil {
		return fmt.Errorf("failed to write profile file at %s: %w", pm.ref.Path, err)
	}

	return nil
}

// Data returns the current profile.
func (pm *ProfileManager) Data() Profile {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.data
}

// SetData replaces the entire profile.
func (pm *ProfileManager) SetData(p Profile) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.data = p
}
