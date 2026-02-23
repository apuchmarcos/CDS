package new_bo

// Agent represents a CDS agent running on a host.
// It holds the agent's address, its TLS credential reference,
// and the list of projects managed through it.
type Agent struct {
	Address       string    `json:"address"`
	CredentialRef string    `json:"credentialRef,omitempty"`
	Projects      []Project `json:"projects,omitempty"`
}
