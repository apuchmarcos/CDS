package new_db

import (
	"fmt"

	nb "github.com/amadeusitgroup/cds/internal/new_bo"
)

// AddProject adds a project under the given host's agent.
// Creates the agent if the host doesn't have one yet.
// The project name must be globally unique across all hosts.
// If the project has no ID, one is auto-generated as name-<random suffix>.
func (m *InventoryManager) AddProject(hostName string, project nb.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil {
		return fmt.Errorf("host %q not found", hostName)
	}

	// Validate global uniqueness
	if existingHost, _ := m.findProjectGlobal(project.Name); existingHost != nil {
		return fmt.Errorf("project %q already exists (on host %q)", project.Name, existingHost.Name)
	}

	// Auto-generate ID if empty
	if project.ID == "" {
		project.ID = generateProjectID(project.Name)
	}

	m.ensureAgent(h)
	h.Agent.Projects = append(h.Agent.Projects, project)
	return nil
}

// RemoveProject removes a project by name. The name is globally unique,
// so no host name is needed.
func (m *InventoryManager) RemoveProject(projectName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.data.Hosts {
		h := &m.data.Hosts[i]
		if h.Agent == nil {
			continue
		}
		for j := range h.Agent.Projects {
			if h.Agent.Projects[j].Name == projectName {
				h.Agent.Projects = append(h.Agent.Projects[:j], h.Agent.Projects[j+1:]...)
				return nil
			}
		}
	}
	return fmt.Errorf("project %q not found", projectName)
}

// GetProject returns a copy of the project with the given name.
// Project names are globally unique, so no host name is needed.
func (m *InventoryManager) GetProject(projectName string) (nb.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return nb.Project{}, fmt.Errorf("project %q not found", projectName)
	}
	return *p, nil
}

// ListProjects returns all projects under the given host's agent.
func (m *InventoryManager) ListProjects(hostName string) []nb.Project {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil || h.Agent == nil {
		return nil
	}

	result := make([]nb.Project, len(h.Agent.Projects))
	copy(result, h.Agent.Projects)
	return result
}

// ListProjectNames returns the names of all projects under the given host's agent.
func (m *InventoryManager) ListProjectNames(hostName string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := m.findHost(hostName)
	if h == nil || h.Agent == nil {
		return nil
	}

	names := make([]string, 0, len(h.Agent.Projects))
	for _, p := range h.Agent.Projects {
		names = append(names, p.Name)
	}
	return names
}

// ListAllProjectNames returns the names of all projects across all hosts.
func (m *InventoryManager) ListAllProjectNames() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var names []string
	for _, h := range m.data.Hosts {
		if h.Agent == nil {
			continue
		}
		for _, p := range h.Agent.Projects {
			names = append(names, p.Name)
		}
	}
	return names
}

// SetProjectInUse marks the given project as in-use and clears InUse on
// all other projects across all hosts.
func (m *InventoryManager) SetProjectInUse(projectName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear all InUse flags first
	for hi := range m.data.Hosts {
		if m.data.Hosts[hi].Agent == nil {
			continue
		}
		for pi := range m.data.Hosts[hi].Agent.Projects {
			m.data.Hosts[hi].Agent.Projects[pi].InUse = false
		}
	}

	// Set the target project (global lookup)
	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return fmt.Errorf("project %q not found", projectName)
	}
	p.InUse = true
	return nil
}

// ClearProjectInUse clears InUse on all projects across all hosts.
func (m *InventoryManager) ClearProjectInUse() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for hi := range m.data.Hosts {
		if m.data.Hosts[hi].Agent == nil {
			continue
		}
		for pi := range m.data.Hosts[hi].Agent.Projects {
			m.data.Hosts[hi].Agent.Projects[pi].InUse = false
		}
	}
}

// GetProjectInUse finds the project currently marked as in-use across all hosts.
// Returns the host name, project name, and the project itself.
// Returns an error if no project is in use.
func (m *InventoryManager) GetProjectInUse() (hostName string, projectName string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, h := range m.data.Hosts {
		if h.Agent == nil {
			continue
		}
		for _, p := range h.Agent.Projects {
			if p.InUse {
				return h.Name, p.Name, nil
			}
		}
	}
	return "", "", fmt.Errorf("no project currently in use")
}

// UpdateProject replaces the project with matching name (globally unique).
func (m *InventoryManager) UpdateProject(project nb.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(project.Name)
	if p == nil {
		return fmt.Errorf("project %q not found", project.Name)
	}
	*p = project
	return nil
}

// ProjectHostName returns the host name that owns the given project.
// Returns empty string if not found.
func (m *InventoryManager) ProjectHostName(projectName string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	h, _ := m.findProjectGlobal(projectName)
	if h == nil {
		return ""
	}
	return h.Name
}

// ProjectConfig returns the config directory for a project, checking
// ConfDir, Flavour.LocalConfDir, and SrcRepo.LocalConfDir in order.
func (m *InventoryManager) ProjectConfig(projectName string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return ""
	}
	switch {
	case len(p.ConfDir) != 0:
		return p.ConfDir
	case len(p.Flavour.LocalConfDir) != 0:
		return p.Flavour.LocalConfDir
	case len(p.SrcRepo.LocalConfDir) != 0:
		return p.SrcRepo.LocalConfDir
	}
	return ""
}
