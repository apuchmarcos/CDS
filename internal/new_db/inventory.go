package new_db

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	nb "github.com/amadeusitgroup/cds/internal/new_bo"
	nfg "github.com/amadeusitgroup/cds/internal/new_config"
)

// InventoryManager manages the Host→Agent→Project hierarchy,
// backed by a Source for persistence.
type InventoryManager struct {
	source nfg.Source
	ref    nfg.SourceRef
	data   nb.Inventory
	mu     sync.Mutex
}

// NewInventoryManager creates an InventoryManager that will load/save
// the inventory through the given Source at ref.Path.
func NewInventoryManager(src nfg.Source, ref nfg.SourceRef) *InventoryManager {
	return &InventoryManager{
		source: src,
		ref:    ref,
		data:   nb.Inventory{Hosts: []nb.Host{}},
	}
}

// Load reads and parses the inventory file from the Source.
// If the file does not exist, the inventory starts empty.
func (m *InventoryManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	exists, err := m.source.Exists(m.ref.Path)
	if err != nil {
		return fmt.Errorf("failed to check inventory file existence: %w", err)
	}
	if !exists {
		m.data = nb.Inventory{Hosts: []nb.Host{}}
		return nil
	}

	raw, err := m.source.Read(m.ref.Path)
	if err != nil {
		return fmt.Errorf("failed to read inventory file at %s: %w", m.ref.Path, err)
	}

	var inv nb.Inventory
	if err := json.Unmarshal(raw, &inv); err != nil {
		return fmt.Errorf("failed to parse inventory file at %s: %w", m.ref.Path, err)
	}

	if inv.Hosts == nil {
		inv.Hosts = []nb.Host{}
	}

	m.data = inv
	return nil
}

// Save writes the current inventory to the Source.
func (m *InventoryManager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	raw, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}

	if err := m.source.Write(m.ref.Path, raw); err != nil {
		return fmt.Errorf("failed to write inventory file at %s: %w", m.ref.Path, err)
	}

	return nil
}

// findHost returns a pointer to the host with the given name, or nil.
// Caller must hold m.mu.
func (m *InventoryManager) findHost(name string) *nb.Host {
	for i := range m.data.Hosts {
		if m.data.Hosts[i].Name == name {
			return &m.data.Hosts[i]
		}
	}
	return nil
}

// findProjectGlobal locates a project by name across all hosts.
// Returns pointers to the owning host and the project, or nil if not found.
// Caller must hold m.mu.
func (m *InventoryManager) findProjectGlobal(projectName string) (*nb.Host, *nb.Project) {
	for i := range m.data.Hosts {
		h := &m.data.Hosts[i]
		if h.Agent == nil {
			continue
		}
		for j := range h.Agent.Projects {
			if h.Agent.Projects[j].Name == projectName {
				return h, &h.Agent.Projects[j]
			}
		}
	}
	return nil, nil
}

// generateProjectID creates a unique project ID as name-<random hex suffix>.
func generateProjectID(name string) string {
	b := make([]byte, 4) // 8 hex chars
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s-%s", name, hex.EncodeToString(b))
}
