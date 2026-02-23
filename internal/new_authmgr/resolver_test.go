package new_authmgr

import (
	"testing"

	"github.com/amadeusitgroup/cds/internal/cos"
	nc "github.com/amadeusitgroup/cds/internal/new_bo"
	nfg "github.com/amadeusitgroup/cds/internal/new_config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestResolver(t *testing.T) *Resolver {
	t.Helper()
	cos.SetMockedFileSystem()
	t.Cleanup(func() { cos.SetRealFileSystem() })

	src := &nfg.LocalFSSource{}
	ref := nfg.SourceRef{Type: nfg.SourceTypeLocalFS, Path: "/tmp/test/.xcds/auth.json"}
	store := NewStore(src, ref)
	require.NoError(t, store.Load())
	return NewResolver(store)
}

func TestResolver_ResolveSSH(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("ssh1", nc.NewSSHCredential("admin", "/keys/id", "/keys/id.pub")))

	user, key, pub, err := r.ResolveSSH("ssh1")
	require.NoError(t, err)
	assert.Equal(t, "admin", user)
	assert.Equal(t, "/keys/id", key)
	assert.Equal(t, "/keys/id.pub", pub)
}

func TestResolver_ResolveSSH_NotFound(t *testing.T) {
	r := setupTestResolver(t)

	_, _, _, err := r.ResolveSSH("missing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credential not found")
}

func TestResolver_ResolveSSH_WrongType(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("token1", nc.NewTokenCredential("tok")))

	_, _, _, err := r.ResolveSSH("token1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected \"ssh\"")
}

func TestResolver_ResolvePassword(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("pwd1", nc.NewPasswordCredential("user@co.com", "secret123")))

	login, pwd, err := r.ResolvePassword("pwd1")
	require.NoError(t, err)
	assert.Equal(t, "user@co.com", login)
	assert.Equal(t, "secret123", pwd)
}

func TestResolver_ResolvePassword_WrongType(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("ssh1", nc.NewSSHCredential("u", "/k", "/k.pub")))

	_, _, err := r.ResolvePassword("ssh1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected \"password\"")
}

func TestResolver_ResolveToken(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("tok1", nc.NewTokenCredential("my-api-token")))

	token, err := r.ResolveToken("tok1")
	require.NoError(t, err)
	assert.Equal(t, "my-api-token", token)
}

func TestResolver_ResolveToken_WrongType(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("pwd1", nc.NewPasswordCredential("u", "p")))

	_, err := r.ResolveToken("pwd1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected \"token\"")
}

func TestResolver_ResolveTLS(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("tls1", nc.NewTLSCredential("/ca.pem", "/cert.pem", "/key.pem")))

	ca, cert, key, err := r.ResolveTLS("tls1")
	require.NoError(t, err)
	assert.Equal(t, "/ca.pem", ca)
	assert.Equal(t, "/cert.pem", cert)
	assert.Equal(t, "/key.pem", key)
}

func TestResolver_ResolveTLS_WrongType(t *testing.T) {
	r := setupTestResolver(t)
	require.NoError(t, r.store.Set("tok1", nc.NewTokenCredential("t")))

	_, _, _, err := r.ResolveTLS("tok1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected \"tls\"")
}
