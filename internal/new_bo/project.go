package new_bo

// FlavourInfo describes the devcontainer flavour configuration for a project.
type FlavourInfo struct {
	Name         string `json:"name"`
	OverrideDir  string `json:"overrideDir,omitempty"`
	LocalConfDir string `json:"localConfDir,omitempty"`
}

// SrcRepoInfo describes source repository configuration for a project.
type SrcRepoInfo struct {
	LocalConfDir string `json:"localConfDir,omitempty"`
	ToClone      bool   `json:"toClone"`
	URI          string `json:"uri,omitempty"`
	Reference    string `json:"reference,omitempty"`
}

// OrchestrationUsage describes how a project uses orchestration (cluster/registry).
type OrchestrationUsage struct {
	Cluster  ClusterUsage  `json:"cluster"`
	Registry RegistryUsage `json:"registry"`
}

// ClusterUsage indicates whether a project uses a cluster.
type ClusterUsage struct {
	Use bool `json:"use"`
}

// RegistryUsage indicates whether a project uses a container registry.
type RegistryUsage struct {
	Use     bool `json:"use"`
	Secured bool `json:"secured,omitempty"`
}

// Project represents a development project managed by CDS.
// Projects live under a Host's Agent in the inventory hierarchy.
// Project names (and IDs) are globally unique across all hosts.
type Project struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	InUse              bool               `json:"inUse"`
	ConfDir            string             `json:"confDir,omitempty"`
	Containers         Containers         `json:"containers,omitempty"`
	Flavour            FlavourInfo        `json:"flavour"`
	SrcRepo            SrcRepoInfo        `json:"srcRepo"`
	UseSshTunnel       bool               `json:"useSshTunnel"`
	OverrideImageTag   string             `json:"overrideImageTag,omitempty"`
	OrchestrationUsage OrchestrationUsage `json:"orchestrationUsage"`
}
