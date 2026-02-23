package new_config

import (
	"encoding/json"
	"testing"

	"github.com/amadeusitgroup/cds/internal/cos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockFS(t *testing.T) {
	t.Helper()
	cos.SetMockedFileSystem()
	t.Cleanup(func() {
		cos.SetRealFileSystem()
	})
}

// --- LocalFSSource Tests ---

func TestLocalFSSource_WriteAndRead(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	path := "/tmp/test/config.json"
	data := []byte(`{"key": "value"}`)

	err := src.Write(path, data)
	require.NoError(t, err)

	got, err := src.Read(path)
	require.NoError(t, err)
	assert.Equal(t, data, got)
}

func TestLocalFSSource_ReadMissingFile(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	_, err := src.Read("/nonexistent/file.json")
	assert.Error(t, err)
}

func TestLocalFSSource_ExistsTrue(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	path := "/tmp/test/exists.json"
	require.NoError(t, src.Write(path, []byte("data")))

	exists, err := src.Exists(path)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestLocalFSSource_ExistsFalse(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	exists, err := src.Exists("/nonexistent/file.json")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestLocalFSSource_DeleteExistingFile(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	path := "/tmp/test/todelete.json"
	require.NoError(t, src.Write(path, []byte("data")))

	err := src.Delete(path)
	require.NoError(t, err)

	exists, err := src.Exists(path)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestLocalFSSource_DeleteNonExistentFile(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	err := src.Delete("/nonexistent/file.json")
	assert.NoError(t, err, "deleting a non-existent file should be idempotent")
}

func TestLocalFSSource_WriteCreatesParentDirs(t *testing.T) {
	setupMockFS(t)
	src := &LocalFSSource{}

	path := "/tmp/deep/nested/dir/file.json"
	data := []byte("nested content")

	err := src.Write(path, data)
	require.NoError(t, err)

	got, err := src.Read(path)
	require.NoError(t, err)
	assert.Equal(t, data, got)
}

// --- ResolveSource Tests ---

func TestResolveSource_LocalFS(t *testing.T) {
	src, err := ResolveSource(SourceRef{Type: SourceTypeLocalFS, Path: "/any/path"})
	require.NoError(t, err)
	assert.IsType(t, &LocalFSSource{}, src)
}

func TestResolveSource_UnknownType(t *testing.T) {
	_, err := ResolveSource(SourceRef{Type: "s3", Path: "bucket/key"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported source type")
}

// --- SourceRef Tests ---

func TestSourceRef_JSONRoundTrip(t *testing.T) {
	ref := SourceRef{Type: SourceTypeLocalFS, Path: "/home/user/.xcds/inventory.json"}

	// Marshal
	data, err := json.Marshal(ref)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"type":"localfs"`)
	assert.Contains(t, string(data), `"path":"/home/user/.xcds/inventory.json"`)

	// Unmarshal
	var got SourceRef
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, ref, got)
}
