package new_bo

// Inventory is the top-level data model for the CDS inventory file.
// It contains a list of hosts, each of which may have an agent with projects.
type Inventory struct {
	Hosts []Host `json:"hosts"`
}
