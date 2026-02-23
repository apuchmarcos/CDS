package new_authmgr

import (
	"fmt"

	nc "github.com/amadeusitgroup/cds/internal/new_bo"
)

// Resolver provides typed credential lookups for consumers.
// Packages like host, tls, ar, and scm use Resolver instead of
// accessing the auth file directly.
type Resolver struct {
	store *Store
}

// NewResolver creates a Resolver backed by the given Store.
func NewResolver(store *Store) *Resolver {
	return &Resolver{store: store}
}

// ResolveSSH looks up an SSH credential and returns its fields.
func (r *Resolver) ResolveSSH(credRef string) (username, keyPath, pubKeyPath string, err error) {
	cred, err := r.store.Get(credRef)
	if err != nil {
		return "", "", "", fmt.Errorf("resolving SSH credential %q: %w", credRef, err)
	}
	if cred.Type != nc.CredSSH {
		return "", "", "", fmt.Errorf("credential %q is type %q, expected %q", credRef, cred.Type, nc.CredSSH)
	}

	keyPath = cred.Metadata[nc.MetaKeyPath]
	pubKeyPath = cred.Metadata[nc.MetaPubKeyPath]
	return cred.Login, keyPath, pubKeyPath, nil
}

// ResolvePassword looks up a password credential and returns its fields.
func (r *Resolver) ResolvePassword(credRef string) (login, password string, err error) {
	cred, err := r.store.Get(credRef)
	if err != nil {
		return "", "", fmt.Errorf("resolving password credential %q: %w", credRef, err)
	}
	if cred.Type != nc.CredPassword {
		return "", "", fmt.Errorf("credential %q is type %q, expected %q", credRef, cred.Type, nc.CredPassword)
	}

	return cred.Login, cred.Password(), nil
}

// ResolveToken looks up a token credential and returns the token string.
func (r *Resolver) ResolveToken(credRef string) (token string, err error) {
	cred, err := r.store.Get(credRef)
	if err != nil {
		return "", fmt.Errorf("resolving token credential %q: %w", credRef, err)
	}
	if cred.Type != nc.CredToken {
		return "", fmt.Errorf("credential %q is type %q, expected %q", credRef, cred.Type, nc.CredToken)
	}

	return cred.Token(), nil
}

// ResolveTLS looks up a TLS credential and returns the certificate paths.
func (r *Resolver) ResolveTLS(credRef string) (caPath, certPath, keyPath string, err error) {
	cred, err := r.store.Get(credRef)
	if err != nil {
		return "", "", "", fmt.Errorf("resolving TLS credential %q: %w", credRef, err)
	}
	if cred.Type != nc.CredTLS {
		return "", "", "", fmt.Errorf("credential %q is type %q, expected %q", credRef, cred.Type, nc.CredTLS)
	}

	caPath = cred.Metadata[nc.MetaCAPath]
	certPath = cred.Metadata[nc.MetaCertPath]
	keyPath = cred.Metadata[nc.MetaTLSKeyPath]
	return caPath, certPath, keyPath, nil
}
