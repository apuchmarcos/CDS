package new_authmgr

import (
	"encoding/json"
	"fmt"
	"sync"

	nc "github.com/amadeusitgroup/cds/internal/new_bo"
	nfg "github.com/amadeusitgroup/cds/internal/new_config"
)

// Store manages the credential repository backed by a Source.
// It is the single owner of the auth file — no other package should
// read or write credentials directly.
type Store struct {
	source nfg.Source
	ref    nfg.SourceRef
	data   nc.AuthStore
	mu     sync.Mutex
}

// NewStore creates a Store that will load/save credentials through the given Source.
func NewStore(src nfg.Source, ref nfg.SourceRef) *Store {
	return &Store{
		source: src,
		ref:    ref,
		data:   nc.NewAuthStore(),
	}
}

// Load reads and parses the auth file from the Source.
// If the file does not exist, the store starts empty.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	exists, err := s.source.Exists(s.ref.Path)
	if err != nil {
		return fmt.Errorf("failed to check auth file existence: %w", err)
	}
	if !exists {
		s.data = nc.NewAuthStore()
		return nil
	}

	raw, err := s.source.Read(s.ref.Path)
	if err != nil {
		return fmt.Errorf("failed to read auth file at %s: %w", s.ref.Path, err)
	}

	var store nc.AuthStore
	if err := json.Unmarshal(raw, &store); err != nil {
		return fmt.Errorf("failed to parse auth file at %s: %w", s.ref.Path, err)
	}

	if store.Credentials == nil {
		store.Credentials = make(map[string]nc.Credential)
	}

	s.data = store
	return nil
}

// Save writes the current credential store to the Source.
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth store: %w", err)
	}

	if err := s.source.Write(s.ref.Path, raw); err != nil {
		return fmt.Errorf("failed to write auth file at %s: %w", s.ref.Path, err)
	}

	return nil
}

// Get returns the credential stored under the given key.
// Returns an error if the key does not exist.
func (s *Store) Get(key string) (nc.Credential, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cred, ok := s.data.Credentials[key]
	if !ok {
		return nc.Credential{}, fmt.Errorf("credential not found: %q", key)
	}
	return cred, nil
}

// Set stores a credential under the given key and persists to the Source.
func (s *Store) Set(key string, cred nc.Credential) error {
	s.mu.Lock()
	s.data.Credentials[key] = cred
	s.mu.Unlock()

	return s.Save()
}

// Delete removes a credential by key and persists to the Source.
// No error if the key does not exist.
func (s *Store) Delete(key string) error {
	s.mu.Lock()
	delete(s.data.Credentials, key)
	s.mu.Unlock()

	return s.Save()
}

// Has checks whether a credential exists for the given key.
func (s *Store) Has(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data.Credentials[key]
	return ok
}

// List returns a copy of all stored credentials.
func (s *Store) List() map[string]nc.Credential {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string]nc.Credential, len(s.data.Credentials))
	for k, v := range s.data.Credentials {
		result[k] = v
	}
	return result
}
