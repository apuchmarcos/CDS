package new_bo

import "encoding/base64"

// CredentialType identifies the kind of credential stored.
type CredentialType string

const (
	CredSSH      CredentialType = "ssh"
	CredPassword CredentialType = "password"
	CredToken    CredentialType = "token"
	CredTLS      CredentialType = "tls"
)

// Metadata keys for each credential type.
const (
	MetaKeyPath    = "keyPath"
	MetaPubKeyPath = "pubKeyPath"
	MetaCAPath     = "caPath"
	MetaCertPath   = "certPath"
	MetaTLSKeyPath = "keyPath"
)

// Credential represents a stored authentication credential.
// The Type field determines which fields are populated:
//   - SSH: Login (username), Metadata[keyPath, pubKeyPath]
//   - Password: Login, Raw (password bytes)
//   - Token: Raw (token bytes)
//   - TLS: Metadata[caPath, certPath, keyPath]
type Credential struct {
	Type     CredentialType    `json:"type"`
	Login    string            `json:"login,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Raw      []byte            `json:"raw,omitempty"`
}

// NewSSHCredential creates a credential for SSH key-based authentication.
func NewSSHCredential(username, keyPath, pubKeyPath string) Credential {
	return Credential{
		Type:  CredSSH,
		Login: username,
		Metadata: map[string]string{
			MetaKeyPath:    keyPath,
			MetaPubKeyPath: pubKeyPath,
		},
	}
}

// NewPasswordCredential creates a credential for username/password authentication.
// The password is stored as raw bytes.
func NewPasswordCredential(login, password string) Credential {
	return Credential{
		Type:  CredPassword,
		Login: login,
		Raw:   []byte(password),
	}
}

// NewTokenCredential creates a credential for token-based authentication.
// The token is stored as raw bytes.
func NewTokenCredential(token string) Credential {
	return Credential{
		Type: CredToken,
		Raw:  []byte(token),
	}
}

// NewTLSCredential creates a credential containing TLS certificate paths.
func NewTLSCredential(caPath, certPath, keyPath string) Credential {
	return Credential{
		Type: CredTLS,
		Metadata: map[string]string{
			MetaCAPath:     caPath,
			MetaCertPath:   certPath,
			MetaTLSKeyPath: keyPath,
		},
	}
}

// Password returns the password string from a Password credential's Raw field.
func (c Credential) Password() string {
	return string(c.Raw)
}

// Token returns the token string from a Token credential's Raw field.
func (c Credential) Token() string {
	return string(c.Raw)
}

// RawBase64 returns the Raw field encoded as base64 (for JSON display).
func (c Credential) RawBase64() string {
	if len(c.Raw) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(c.Raw)
}
