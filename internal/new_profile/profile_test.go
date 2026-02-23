package new_profile

import (
	"testing"

	"github.com/amadeusitgroup/cds/internal/cos"
	nfg "github.com/amadeusitgroup/cds/internal/new_config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestProfile(t *testing.T) *ProfileManager {
	t.Helper()
	cos.SetMockedFileSystem()
	t.Cleanup(func() { cos.SetRealFileSystem() })

	src := &nfg.LocalFSSource{}
	ref := nfg.SourceRef{Type: nfg.SourceTypeLocalFS, Path: "/tmp/test/.xcds/profile.json"}
	return NewProfileManager(src, ref)
}

func TestProfileManager_LoadEmpty(t *testing.T) {
	pm := setupTestProfile(t)

	err := pm.Load()
	require.NoError(t, err)
	assert.Equal(t, Profile{}, pm.Data())
}

func TestProfileManager_SaveAndLoad(t *testing.T) {
	pm := setupTestProfile(t)
	require.NoError(t, pm.Load())

	pm.SetData(Profile{Holder: "my-profile"})
	require.NoError(t, pm.Save())

	pm2 := NewProfileManager(pm.source, pm.ref)
	require.NoError(t, pm2.Load())

	assert.Equal(t, "my-profile", pm2.Data().Holder)
}

func TestProfileManager_LoadInvalidJSON(t *testing.T) {
	pm := setupTestProfile(t)

	require.NoError(t, pm.source.Write(pm.ref.Path, []byte("not json")))

	err := pm.Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse profile")
}

func TestProfileManager_SetData(t *testing.T) {
	pm := setupTestProfile(t)
	require.NoError(t, pm.Load())

	pm.SetData(Profile{Holder: "updated"})
	assert.Equal(t, "updated", pm.Data().Holder)
}
