package new_bo

// OrchestrationInfo describes orchestration (e.g. k3s) state on a host.
type OrchestrationInfo struct {
	Name         string       `json:"name"`
	RegistryInfo RegistryInfo `json:"registry"`
	State        string       `json:"state"`
}

// RegistryInfo describes the container registry state on a host.
type RegistryInfo struct {
	State string `json:"state"`
	Port  int    `json:"port"`
}

// Host represents a remote or local machine managed by CDS.
// In the inventory hierarchy: Host → Agent → Projects.
type Host struct {
	Name          string             `json:"name"`
	CredentialRef string             `json:"credentialRef,omitempty"`
	IsDefault     bool               `json:"isDefault"`
	Agent         *Agent             `json:"agent,omitempty"`
	Orchestration *OrchestrationInfo `json:"orchestration,omitempty"`
}
