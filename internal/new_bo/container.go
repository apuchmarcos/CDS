package new_bo

import (
	"fmt"
	"slices"
)

// PortMapping maps a container port (e.g. "22/tcp") to a host port.
type PortMapping map[string]int

// ContainerID is a unique identifier for a container.
type ContainerID = string

// ContainerName is the human-readable name of a container.
type ContainerName = string

// ContainerStatus represents the lifecycle state of a container.
type ContainerStatus int

// ContainerUser is the user defined in the container image.
type ContainerUser = string

// ContainerRemoteUser is the user used for remote connections.
type ContainerRemoteUser = string

const (
	KContainerStatusRunning ContainerStatus = iota
	KContainerStatusExited
	KContainerStatusDeleted
	KContainerStatusUnknown
	KContainerStatusArchived
)

const (
	KSSHPortMapping = "22/tcp"
)

var containerStatuses = map[ContainerStatus]string{
	KContainerStatusRunning:  "running",
	KContainerStatusExited:   "exited",
	KContainerStatusDeleted:  "deleted",
	KContainerStatusUnknown:  "unknown",
	KContainerStatusArchived: "archived",
}

// Container represents a development container within a project.
type Container struct {
	Id             ContainerID         `json:"id"`
	Name           ContainerName       `json:"name"`
	Pmapping       PortMapping         `json:"portMapping,omitempty"`
	Status         ContainerStatus     `json:"status"`
	ExpectedStatus ContainerStatus     `json:"expectedStatus"`
	RemoteUser     ContainerRemoteUser `json:"remoteUser,omitempty"`
	User           ContainerUser       `json:"user,omitempty"`
}

func (status ContainerStatus) ToString() string {
	return FContainerStatus(status)
}

// FContainerStatus formats a ContainerStatus as a human-readable string.
func FContainerStatus(cs ContainerStatus) string {
	fstatus, ok := containerStatuses[cs]
	if !ok {
		return containerStatuses[KContainerStatusUnknown]
	}
	return fstatus
}

// SContainerStatus parses a string into a ContainerStatus.
func SContainerStatus(containerStatusStr string) ContainerStatus {
	for status, fStatus := range containerStatuses {
		if containerStatusStr == fStatus {
			return status
		}
	}
	return KContainerStatusUnknown
}

// AddPort adds a port mapping entry. Returns an error if the key already exists.
func (c *Container) AddPort(key string, val int) error {
	if c.Pmapping == nil {
		c.Pmapping = make(PortMapping)
	}
	if _, ok := c.Pmapping[key]; ok {
		return fmt.Errorf("attempt to update existing key (%v)'s value from (%v) to (%v)",
			key, c.Pmapping[key], val)
	}
	c.Pmapping[key] = val
	return nil
}

// PortMap returns the container's port mapping.
func (c *Container) PortMap() PortMapping {
	return c.Pmapping
}

// Containers is a slice of Container with helper methods.
type Containers []Container

// Contains checks if a container with the given name exists in the slice.
func (cs *Containers) Contains(name ContainerName) bool {
	return slices.ContainsFunc(*cs, func(c Container) bool {
		return c.Name == name
	})
}

// ContainsId checks if a container with the given ID exists in the slice.
func (cs *Containers) ContainsId(id ContainerID) bool {
	return slices.ContainsFunc(*cs, func(c Container) bool {
		return c.Id == id
	})
}

// Get returns the container with the given name, or an empty Container if not found.
func (cs *Containers) Get(name ContainerName) Container {
	for _, container := range *cs {
		if container.Name == name {
			return container
		}
	}
	return Container{}
}

// GetById returns the container with the given ID, or an empty Container if not found.
func (cs *Containers) GetById(id ContainerID) Container {
	for _, container := range *cs {
		if container.Id == id {
			return container
		}
	}
	return Container{}
}
