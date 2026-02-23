package new_config

import "fmt"

// Source abstracts file I/O for configuration storage.
// Implementations: LocalFSSource (now), GitSource, S3Source, CyberArkSource (future).
type Source interface {
	// Read returns the contents of the file at the given path.
	// Returns an error if the file does not exist or cannot be read.
	Read(path string) ([]byte, error)

	// Write persists data to the given path, creating parent directories as needed.
	Write(path string, data []byte) error

	// Exists checks whether a file exists at the given path.
	// Returns (false, nil) when the file simply does not exist.
	// Returns a non-nil error only for unexpected failures (e.g. permission issues).
	Exists(path string) (bool, error)

	// Delete removes the file at the given path.
	// Idempotent: returns nil if the file does not exist.
	Delete(path string) error
}

// SourceType identifies the transport/storage backend for a Source.
const (
	SourceTypeLocalFS = "localfs"
)

// SourceRef is a serializable reference to a Source-backed file.
// It specifies the transport type and the path within that transport.
type SourceRef struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// ResolveSource creates a concrete Source implementation from a SourceRef.
// Currently only "localfs" is supported. Future types (git, s3, etc.) will be added here.
func ResolveSource(ref SourceRef) (Source, error) {
	switch ref.Type {
	case SourceTypeLocalFS:
		return &LocalFSSource{}, nil
	default:
		return nil, fmt.Errorf("unsupported source type: %q", ref.Type)
	}
}
