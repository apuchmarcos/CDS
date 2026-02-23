package new_bo

// AuthStore is the top-level data model for the CDS auth file.
// It maps credential keys to Credential values.
type AuthStore struct {
	Credentials map[string]Credential `json:"credentials"`
}

// NewAuthStore returns an empty AuthStore with an initialized map.
func NewAuthStore() AuthStore {
	return AuthStore{
		Credentials: make(map[string]Credential),
	}
}
