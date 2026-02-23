package new_db

import (
	"testing"

	nb "github.com/amadeusitgroup/cds/internal/new_bo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Basic CRUD ----------

func TestHost_AddAndGet(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "alpha"}))
	require.NoError(t, mgr.AddHost(nb.Host{Name: "beta"}))

	h, err := mgr.GetHost("alpha")
	require.NoError(t, err)
	assert.Equal(t, "alpha", h.Name)
}

func TestHost_AddDuplicate(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "dup"}))
	err := mgr.AddHost(nb.Host{Name: "dup"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestHost_Remove(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "rm-me"}))
	require.NoError(t, mgr.RemoveHost("rm-me"))

	assert.False(t, mgr.HasHost("rm-me"))
}

func TestHost_RemoveNotFound(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	err := mgr.RemoveHost("ghost")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHost_GetNotFound(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	_, err := mgr.GetHost("nope")
	require.Error(t, err)
}

func TestHost_ListHostNames(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "x"}))
	require.NoError(t, mgr.AddHost(nb.Host{Name: "y"}))

	names := mgr.ListHostNames()
	assert.ElementsMatch(t, []string{"x", "y"}, names)
}

func TestHost_Has(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	assert.False(t, mgr.HasHost("missing"))
	require.NoError(t, mgr.AddHost(nb.Host{Name: "present"}))
	assert.True(t, mgr.HasHost("present"))
}

// ---------- Default host ----------

func TestHost_DefaultHost(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "a"}))
	require.NoError(t, mgr.AddHost(nb.Host{Name: "b"}))

	require.NoError(t, mgr.SetDefault("b"))

	h, err := mgr.GetDefaultHost()
	require.NoError(t, err)
	assert.Equal(t, "b", h.Name)

	assert.Equal(t, "b", mgr.GetDefaultHostName())
}

func TestHost_SetDefaultClears(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "a", IsDefault: true}))
	require.NoError(t, mgr.AddHost(nb.Host{Name: "b"}))

	require.NoError(t, mgr.SetDefault("b"))

	a, _ := mgr.GetHost("a")
	assert.False(t, a.IsDefault)

	b, _ := mgr.GetHost("b")
	assert.True(t, b.IsDefault)
}

func TestHost_SetDefaultNotFound(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	err := mgr.SetDefault("missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHost_NoDefault(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "nodef"}))

	_, err := mgr.GetDefaultHost()
	require.Error(t, err)
	assert.Equal(t, "", mgr.GetDefaultHostName())
}

// ---------- Credential & orchestration ----------

func TestHost_UpdateCredentialRef(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "h"}))
	require.NoError(t, mgr.UpdateHostCredentialRef("h", "new-cred"))

	h, _ := mgr.GetHost("h")
	assert.Equal(t, "new-cred", h.CredentialRef)
}

func TestHost_OrcInfoSetters(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "h"}))

	require.NoError(t, mgr.SetOrcInfoName("h", "k3s"))
	require.NoError(t, mgr.SetOrcInfoState("h", "running"))
	require.NoError(t, mgr.SetOrcInfoRegistryState("h", "active"))
	require.NoError(t, mgr.SetOrcInfoRegistryPort("h", 5000))

	h, _ := mgr.GetHost("h")
	require.NotNil(t, h.Orchestration)
	assert.Equal(t, "k3s", h.Orchestration.Name)
	assert.Equal(t, "running", h.Orchestration.State)
	assert.Equal(t, "active", h.Orchestration.RegistryInfo.State)
	assert.Equal(t, 5000, h.Orchestration.RegistryInfo.Port)
}

func TestHost_GetRegistryInfo(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	require.NoError(t, mgr.AddHost(nb.Host{Name: "h"}))

	// Before setting, should be empty
	ri, err := mgr.GetRegistryInfo("h")
	require.NoError(t, err)
	assert.Equal(t, "", ri.State)

	require.NoError(t, mgr.SetOrcInfoRegistryState("h", "ready"))
	require.NoError(t, mgr.SetOrcInfoRegistryPort("h", 5000))

	ri, err = mgr.GetRegistryInfo("h")
	require.NoError(t, err)
	assert.Equal(t, "ready", ri.State)
	assert.Equal(t, 5000, ri.Port)
}

func TestHost_GetRegistryInfoNotFound(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	_, err := mgr.GetRegistryInfo("missing")
	require.Error(t, err)
}
