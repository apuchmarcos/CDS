package new_db

import (
	"fmt"

	nb "github.com/amadeusitgroup/cds/internal/new_bo"
)

// AddContainer adds a container to the specified project.
// Project names are globally unique, so no host name is needed.
func (m *InventoryManager) AddContainer(projectName string, c nb.Container) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return fmt.Errorf("project %q not found", projectName)
	}

	// Check for duplicate name
	for _, existing := range p.Containers {
		if existing.Name == c.Name {
			return fmt.Errorf("container %q already exists in project %q", c.Name, projectName)
		}
	}

	p.Containers = append(p.Containers, c)
	return nil
}

// RemoveContainer removes a container by name from the specified project.
func (m *InventoryManager) RemoveContainer(projectName, containerName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return fmt.Errorf("project %q not found", projectName)
	}

	for i := range p.Containers {
		if p.Containers[i].Name == nb.ContainerName(containerName) {
			p.Containers = append(p.Containers[:i], p.Containers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("container %q not found in project %q", containerName, projectName)
}

// GetContainer returns a copy of the container with the given name.
func (m *InventoryManager) GetContainer(projectName, containerName string) (nb.Container, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return nb.Container{}, fmt.Errorf("project %q not found", projectName)
	}

	for _, c := range p.Containers {
		if c.Name == nb.ContainerName(containerName) {
			return c, nil
		}
	}
	return nb.Container{}, fmt.Errorf("container %q not found in project %q", containerName, projectName)
}

// ListContainerNames returns the names of all containers in the specified project.
func (m *InventoryManager) ListContainerNames(projectName string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return nil
	}

	names := make([]string, 0, len(p.Containers))
	for _, c := range p.Containers {
		names = append(names, string(c.Name))
	}
	return names
}

// ContainerSSHPort returns the SSH port for the given container, or -1 if not found.
func (m *InventoryManager) ContainerSSHPort(projectName, containerName string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return -1
	}

	for _, c := range p.Containers {
		if c.Name == nb.ContainerName(containerName) {
			if port, ok := c.Pmapping[nb.KSSHPortMapping]; ok {
				return port
			}
			return -1
		}
	}
	return -1
}

// ContainerRemoteUser returns the remote user for the given container.
func (m *InventoryManager) ContainerRemoteUser(projectName, containerName string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return ""
	}

	for _, c := range p.Containers {
		if c.Name == nb.ContainerName(containerName) {
			return string(c.RemoteUser)
		}
	}
	return ""
}

// ClearContainers removes all containers from the specified project.
func (m *InventoryManager) ClearContainers(projectName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, p := m.findProjectGlobal(projectName)
	if p == nil {
		return fmt.Errorf("project %q not found", projectName)
	}

	p.Containers = nil
	return nil
}
