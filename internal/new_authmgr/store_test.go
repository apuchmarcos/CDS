package new_authmgr

import (
	"testing"

	"github.com/amadeusitgroup/cds/internal/cos"
	nc "github.com/amadeusitgroup/cds/internal/new_bo"
	nfg "github.com/amadeusitgroup/cds/internal/new_config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStore(t *testing.T) *Store {
	t.Helper()
	cos.SetMockedFileSystem()
	t.Cleanup(func() { cos.SetRealFileSystem() })

	src := &nfg.LocalFSSource{}
	ref := nfg.SourceRef{Type: nfg.SourceTypeLocalFS, Path: "/tmp/test/.xcds/auth.json"}
	return NewStore(src, ref)
}

func TestStore_LoadEmpty(t *testing.T) {
	store := setupTestStore(t)

	err := store.Load()
	require.NoError(t, err, "loading non-existent file should start empty")
	assert.Equal(t, 0, len(store.List()))
}

func TestStore_SetAndGet(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	cred := nc.NewSSHCredential("admin", "/keys/id_rsa", "/keys/id_rsa.pub")
	require.NoError(t, store.Set("host-myhost-ssh", cred))

	got, err := store.Get("host-myhost-ssh")
	require.NoError(t, err)
	assert.Equal(t, nc.CredSSH, got.Type)
	assert.Equal(t, "admin", got.Login)
	assert.Equal(t, "/keys/id_rsa", got.Metadata[nc.MetaKeyPath])
	assert.Equal(t, "/keys/id_rsa.pub", got.Metadata[nc.MetaPubKeyPath])
}

func TestStore_GetNotFound(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	_, err := store.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credential not found")
}

func TestStore_Has(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	assert.False(t, store.Has("key"))

	require.NoError(t, store.Set("key", nc.NewTokenCredential("tok123")))
	assert.True(t, store.Has("key"))
}

func TestStore_Delete(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	require.NoError(t, store.Set("key", nc.NewTokenCredential("tok")))
	assert.True(t, store.Has("key"))

	require.NoError(t, store.Delete("key"))
	assert.False(t, store.Has("key"))
}

func TestStore_DeleteNonExistent(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	err := store.Delete("nonexistent")
	assert.NoError(t, err, "deleting non-existent key should not error")
}

func TestStore_List(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	require.NoError(t, store.Set("a", nc.NewTokenCredential("t1")))
	require.NoError(t, store.Set("b", nc.NewPasswordCredential("user", "pass")))

	list := store.List()
	assert.Len(t, list, 2)
	assert.Equal(t, nc.CredToken, list["a"].Type)
	assert.Equal(t, nc.CredPassword, list["b"].Type)
}

func TestStore_SaveAndReload(t *testing.T) {
	cos.SetMockedFileSystem()
	t.Cleanup(func() { cos.SetRealFileSystem() })

	src := &nfg.LocalFSSource{}
	ref := nfg.SourceRef{Type: nfg.SourceTypeLocalFS, Path: "/tmp/test/.xcds/auth.json"}

	// Create and populate store
	store1 := NewStore(src, ref)
	require.NoError(t, store1.Load())
	require.NoError(t, store1.Set("ssh-key", nc.NewSSHCredential("root", "/k", "/k.pub")))
	require.NoError(t, store1.Set("tls-cert", nc.NewTLSCredential("/ca.pem", "/cert.pem", "/key.pem")))

	// Create a fresh store and load from same path
	store2 := NewStore(src, ref)
	require.NoError(t, store2.Load())

	assert.Len(t, store2.List(), 2)

	got, err := store2.Get("ssh-key")
	require.NoError(t, err)
	assert.Equal(t, "root", got.Login)

	got2, err := store2.Get("tls-cert")
	require.NoError(t, err)
	assert.Equal(t, nc.CredTLS, got2.Type)
	assert.Equal(t, "/ca.pem", got2.Metadata[nc.MetaCAPath])
}

func TestStore_SetOverwrites(t *testing.T) {
	store := setupTestStore(t)
	require.NoError(t, store.Load())

	require.NoError(t, store.Set("key", nc.NewTokenCredential("old")))
	require.NoError(t, store.Set("key", nc.NewTokenCredential("new")))

	got, err := store.Get("key")
	require.NoError(t, err)
	assert.Equal(t, "new", got.Token())
}
