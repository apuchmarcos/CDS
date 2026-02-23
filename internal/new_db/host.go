package new_db

import (
	"fmt"

	nb "github.com/amadeusitgroup/cds/internal/new_bo"
)

// AddHost adds a host to the inventory. If a host with the same name
// already exists, it returns an error.
func (m *InventoryManager) AddHost(host nb.Host) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.findHost(host.Name) != nil {
		return fmt.Errorf("host %q already exists", host.Name)
	}
	m.data.Hosts = append(m.data.Hosts, host)
	return nil
}

// RemoveHost removes a host by name. Returns an error if not found.
func (m *InventoryManager) RemoveHost(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.data.Hosts {
		if m.data.Hosts[i].Name == name {
			m.data.Hosts = append(m.data.Hosts[:i], m.data.Hosts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("host %q not found", name)
}

// GetHost returns a copy of the host with the given name.
func (m *InventoryManager) GetHost(name string) (nb.Host, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(name)
	if h == nil {
		return nb.Host{}, fmt.Errorf("host %q not found", name)
	}
	return *h, nil
}

// ListHosts returns a copy of all hosts.
func (m *InventoryManager) ListHosts() []nb.Host {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]nb.Host, len(m.data.Hosts))
	copy(result, m.data.Hosts)
	return result
}

// ListHostNames returns the names of all hosts.
func (m *InventoryManager) ListHostNames() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	names := make([]string, 0, len(m.data.Hosts))
	for _, h := range m.data.Hosts {
		names = append(names, h.Name)
	}
	return names
}

// HasHost checks whether a host with the given name exists.
func (m *InventoryManager) HasHost(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.findHost(name) != nil
}

// GetDefaultHost returns the host marked as default.
// Returns an error if no default is set.
func (m *InventoryManager) GetDefaultHost() (nb.Host, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, h := range m.data.Hosts {
		if h.IsDefault {
			return h, nil
		}
	}
	return nb.Host{}, fmt.Errorf("no default host set")
}

// GetDefaultHostName returns the name of the default host, or empty string.
func (m *InventoryManager) GetDefaultHostName() string {
	h, err := m.GetDefaultHost()
	if err != nil {
		return ""
	}
	return h.Name
}

// SetDefault marks the given host as default and clears default on all others.
func (m *InventoryManager) SetDefault(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	for i := range m.data.Hosts {
		if m.data.Hosts[i].Name == name {
			m.data.Hosts[i].IsDefault = true
			found = true
		} else {
			m.data.Hosts[i].IsDefault = false
		}
	}
	if !found {
		return fmt.Errorf("host %q not found", name)
	}
	return nil
}

// UpdateHostCredentialRef updates the credential reference for a host.
func (m *InventoryManager) UpdateHostCredentialRef(hostName, credRef string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return fmt.Errorf("host %q not found", hostName)
	}
	h.CredentialRef = credRef
	return nil
}

// SetOrcInfoName sets the orchestration name for a host.
func (m *InventoryManager) SetOrcInfoName(hostName, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return fmt.Errorf("host %q not found", hostName)
	}
	if h.Orchestration == nil {
		h.Orchestration = &nb.OrchestrationInfo{}
	}
	h.Orchestration.Name = name
	return nil
}

// SetOrcInfoState sets the orchestration state for a host.
func (m *InventoryManager) SetOrcInfoState(hostName, state string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return fmt.Errorf("host %q not found", hostName)
	}
	if h.Orchestration == nil {
		h.Orchestration = &nb.OrchestrationInfo{}
	}
	h.Orchestration.State = state
	return nil
}

// SetOrcInfoRegistryState sets the orchestration registry state for a host.
func (m *InventoryManager) SetOrcInfoRegistryState(hostName, state string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return fmt.Errorf("host %q not found", hostName)
	}
	if h.Orchestration == nil {
		h.Orchestration = &nb.OrchestrationInfo{}
	}
	h.Orchestration.RegistryInfo.State = state
	return nil
}

// SetOrcInfoRegistryPort sets the orchestration registry port for a host.
func (m *InventoryManager) SetOrcInfoRegistryPort(hostName string, port int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return fmt.Errorf("host %q not found", hostName)
	}
	if h.Orchestration == nil {
		h.Orchestration = &nb.OrchestrationInfo{}
	}
	h.Orchestration.RegistryInfo.Port = port
	return nil
}

// GetRegistryInfo returns the registry info for a host.
func (m *InventoryManager) GetRegistryInfo(hostName string) (nb.RegistryInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return nb.RegistryInfo{}, fmt.Errorf("host %q not found", hostName)
	}
	if h.Orchestration == nil {
		return nb.RegistryInfo{}, nil
	}
	return h.Orchestration.RegistryInfo, nil
}

// EnsureAgent ensures the host has an Agent. Creates one if nil.
// Caller must hold m.mu.
func (m *InventoryManager) ensureAgent(h *nb.Host) {
	if h.Agent == nil {
		h.Agent = &nb.Agent{Projects: []nb.Project{}}
	}
	if h.Agent.Projects == nil {
		h.Agent.Projects = []nb.Project{}
	}
}
