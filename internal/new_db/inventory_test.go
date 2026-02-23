package new_db

import (
	"testing"

	"github.com/amadeusitgroup/cds/internal/cos"
	nb "github.com/amadeusitgroup/cds/internal/new_bo"
	nfg "github.com/amadeusitgroup/cds/internal/new_config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testInventoryPath = "/tmp/test/.xcds/inventory.json"

func setupTestManager(t *testing.T) *InventoryManager {
	t.Helper()
	cos.SetMockedFileSystem()
	t.Cleanup(func() { cos.SetRealFileSystem() })

	src := &nfg.LocalFSSource{}
	ref := nfg.SourceRef{Type: nfg.SourceTypeLocalFS, Path: testInventoryPath}
	return NewInventoryManager(src, ref)
}

func TestInventoryManager_LoadEmpty(t *testing.T) {
	mgr := setupTestManager(t)

	err := mgr.Load()
	require.NoError(t, err, "loading non-existent file should start empty")
	assert.Empty(t, mgr.ListHosts())
}

func TestInventoryManager_SaveAndLoad(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "h1", CredentialRef: "cred-h1"}))
	require.NoError(t, mgr.AddHost(nb.Host{Name: "h2"}))
	require.NoError(t, mgr.Save())

	// reload into a fresh manager
	mgr2 := NewInventoryManager(mgr.source, mgr.ref)
	require.NoError(t, mgr2.Load())
	assert.Equal(t, 2, len(mgr2.ListHosts()))
	h, err := mgr2.GetHost("h1")
	require.NoError(t, err)
	assert.Equal(t, "cred-h1", h.CredentialRef)
}

func TestInventoryManager_LoadInvalidJSON(t *testing.T) {
	mgr := setupTestManager(t)

	require.NoError(t, mgr.source.Write(mgr.ref.Path, []byte("not json")))

	err := mgr.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse inventory")
}

func TestInventoryManager_SaveRoundTrip(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	host := nb.Host{
		Name:          "dev-box",
		CredentialRef: "ssh-dev",
		IsDefault:     true,
		Agent: &nb.Agent{
			Address:       "localhost:9090",
			CredentialRef: "tls-agent",
			Projects: []nb.Project{
				{Name: "proj-a", InUse: true},
			},
		},
	}
	require.NoError(t, mgr.AddHost(host))
	require.NoError(t, mgr.Save())

	mgr2 := NewInventoryManager(mgr.source, mgr.ref)
	require.NoError(t, mgr2.Load())

	h, err := mgr2.GetHost("dev-box")
	require.NoError(t, err)
	assert.True(t, h.IsDefault)
	require.NotNil(t, h.Agent)
	assert.Equal(t, "localhost:9090", h.Agent.Address)
	assert.Equal(t, 1, len(h.Agent.Projects))
	assert.Equal(t, "proj-a", h.Agent.Projects[0].Name)
	assert.True(t, h.Agent.Projects[0].InUse)
}
