package new_db

import (
	"testing"

	nb "github.com/amadeusitgroup/cds/internal/new_bo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mgrWithHost returns a loaded manager with one host already added.
func mgrWithHost(t *testing.T, hostName string) *InventoryManager {
	t.Helper()
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())
	require.NoError(t, mgr.AddHost(nb.Host{Name: hostName}))
	return mgr
}

// ---------- Basic CRUD ----------

func TestProject_AddAndGet(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p1"}))

	p, err := mgr.GetProject("p1")
	require.NoError(t, err)
	assert.Equal(t, "p1", p.Name)
	assert.NotEmpty(t, p.ID, "ID should be auto-generated")
	assert.Contains(t, p.ID, "p1-", "ID should start with name-")
}

func TestProject_AddCreatesAgent(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	h, _ := mgr.GetHost("h")
	assert.Nil(t, h.Agent)

	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))

	h, _ = mgr.GetHost("h")
	require.NotNil(t, h.Agent)
}

func TestProject_AddDuplicate(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "dup"}))
	err := mgr.AddProject("h", nb.Project{Name: "dup"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestProject_AddDuplicateAcrossHosts(t *testing.T) {
	mgr := mgrWithHost(t, "h1")
	require.NoError(t, mgr.AddHost(nb.Host{Name: "h2"}))

	require.NoError(t, mgr.AddProject("h1", nb.Project{Name: "global-unique"}))
	err := mgr.AddProject("h2", nb.Project{Name: "global-unique"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	assert.Contains(t, err.Error(), "h1")
}

func TestProject_AddPreservesExplicitID(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p", ID: "custom-id-123"}))

	p, _ := mgr.GetProject("p")
	assert.Equal(t, "custom-id-123", p.ID)
}

func TestProject_AddHostNotFound(t *testing.T) {
	mgr := setupTestManager(t)
	require.NoError(t, mgr.Load())

	err := mgr.AddProject("nonexistent", nb.Project{Name: "p"})
	require.Error(t, err)
}

func TestProject_Remove(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))

	require.NoError(t, mgr.RemoveProject("p"))
	assert.Empty(t, mgr.ListProjectNames("h"))
}

func TestProject_RemoveNotFound(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	err := mgr.RemoveProject("ghost")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestProject_GetNotFound(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	_, err := mgr.GetProject("nope")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestProject_ListProjectNames(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "a"}))
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "b"}))
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "c"}))

	names := mgr.ListProjectNames("h")
	assert.ElementsMatch(t, []string{"a", "b", "c"}, names)
}

func TestProject_ListProjectsNilAgent(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	assert.Nil(t, mgr.ListProjects("h"))
}

func TestProject_ListAllProjectNames(t *testing.T) {
	mgr := mgrWithHost(t, "h1")
	require.NoError(t, mgr.AddHost(nb.Host{Name: "h2"}))
	require.NoError(t, mgr.AddProject("h1", nb.Project{Name: "p1"}))
	require.NoError(t, mgr.AddProject("h2", nb.Project{Name: "p2"}))

	names := mgr.ListAllProjectNames()
	assert.ElementsMatch(t, []string{"p1", "p2"}, names)
}

// ---------- InUse ----------

func TestProject_SetProjectInUse(t *testing.T) {
	mgr := mgrWithHost(t, "h1")
	require.NoError(t, mgr.AddHost(nb.Host{Name: "h2"}))
	require.NoError(t, mgr.AddProject("h1", nb.Project{Name: "p1"}))
	require.NoError(t, mgr.AddProject("h2", nb.Project{Name: "p2"}))

	require.NoError(t, mgr.SetProjectInUse("p2"))

	host, proj, err := mgr.GetProjectInUse()
	require.NoError(t, err)
	assert.Equal(t, "h2", host)
	assert.Equal(t, "p2", proj)

	// p1 should not be in use
	p1, _ := mgr.GetProject("p1")
	assert.False(t, p1.InUse)
}

func TestProject_SetProjectInUseClears(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "a"}))
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "b"}))

	// Set a, then switch to b
	require.NoError(t, mgr.SetProjectInUse("a"))
	require.NoError(t, mgr.SetProjectInUse("b"))

	a, _ := mgr.GetProject("a")
	assert.False(t, a.InUse)
	b, _ := mgr.GetProject("b")
	assert.True(t, b.InUse)
}

func TestProject_ClearProjectInUse(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.SetProjectInUse("p"))

	mgr.ClearProjectInUse()

	_, _, err := mgr.GetProjectInUse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no project currently in use")
}

func TestProject_GetProjectInUseNone(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))

	_, _, err := mgr.GetProjectInUse()
	require.Error(t, err)
}

// ---------- UpdateProject ----------

func TestProject_Update(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p", ConfDir: "/old"}))

	p, _ := mgr.GetProject("p")
	p.ConfDir = "/new"
	require.NoError(t, mgr.UpdateProject(p))

	updated, _ := mgr.GetProject("p")
	assert.Equal(t, "/new", updated.ConfDir)
}

func TestProject_UpdateNotFound(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	err := mgr.UpdateProject(nb.Project{Name: "ghost"})
	require.Error(t, err)
}

// ---------- ProjectHostName ----------

func TestProject_ProjectHostName(t *testing.T) {
	mgr := mgrWithHost(t, "h1")
	require.NoError(t, mgr.AddHost(nb.Host{Name: "h2"}))
	require.NoError(t, mgr.AddProject("h2", nb.Project{Name: "p"}))

	assert.Equal(t, "h2", mgr.ProjectHostName("p"))
	assert.Equal(t, "", mgr.ProjectHostName("missing"))
}

// ---------- ProjectConfig ----------

func TestProject_ProjectConfig(t *testing.T) {
	mgr := mgrWithHost(t, "h")

	require.NoError(t, mgr.AddProject("h", nb.Project{
		Name:    "p-confdir",
		ConfDir: "/conf",
	}))
	assert.Equal(t, "/conf", mgr.ProjectConfig("p-confdir"))

	require.NoError(t, mgr.AddProject("h", nb.Project{
		Name:    "p-flavour",
		Flavour: nb.FlavourInfo{LocalConfDir: "/flavour"},
	}))
	assert.Equal(t, "/flavour", mgr.ProjectConfig("p-flavour"))

	require.NoError(t, mgr.AddProject("h", nb.Project{
		Name:    "p-src",
		SrcRepo: nb.SrcRepoInfo{LocalConfDir: "/src"},
	}))
	assert.Equal(t, "/src", mgr.ProjectConfig("p-src"))

	assert.Equal(t, "", mgr.ProjectConfig("missing"))
}

// ---------- Container operations ----------

func TestProject_AddContainer(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))

	c := nb.Container{
		Name:       "web",
		Id:         "abc123",
		RemoteUser: "dev",
		Pmapping:   nb.PortMapping{nb.KSSHPortMapping: 2222},
	}
	require.NoError(t, mgr.AddContainer("p", c))

	got, err := mgr.GetContainer("p", "web")
	require.NoError(t, err)
	assert.Equal(t, "web", string(got.Name))
	assert.Equal(t, 2222, got.Pmapping[nb.KSSHPortMapping])
}

func TestProject_AddContainerDuplicate(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))

	require.NoError(t, mgr.AddContainer("p", nb.Container{Name: "c"}))
	err := mgr.AddContainer("p", nb.Container{Name: "c"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestProject_RemoveContainer(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{Name: "c"}))

	require.NoError(t, mgr.RemoveContainer("p", "c"))
	assert.Empty(t, mgr.ListContainerNames("p"))
}

func TestProject_RemoveContainerNotFound(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))

	err := mgr.RemoveContainer("p", "ghost")
	require.Error(t, err)
}

func TestProject_ContainerSSHPort(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{
		Name:     "c",
		Pmapping: nb.PortMapping{nb.KSSHPortMapping: 3333},
	}))

	assert.Equal(t, 3333, mgr.ContainerSSHPort("p", "c"))
	assert.Equal(t, -1, mgr.ContainerSSHPort("p", "missing"))
	assert.Equal(t, -1, mgr.ContainerSSHPort("missing-project", "c"))
}

func TestProject_ContainerRemoteUser(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{
		Name:       "c",
		RemoteUser: "testuser",
	}))

	assert.Equal(t, "testuser", mgr.ContainerRemoteUser("p", "c"))
	assert.Equal(t, "", mgr.ContainerRemoteUser("p", "missing"))
}

func TestProject_ClearContainers(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{Name: "c1"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{Name: "c2"}))

	require.NoError(t, mgr.ClearContainers("p"))
	assert.Empty(t, mgr.ListContainerNames("p"))
}

func TestProject_ListContainerNames(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{Name: "a"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{Name: "b"}))

	names := mgr.ListContainerNames("p")
	assert.ElementsMatch(t, []string{"a", "b"}, names)
}

func TestProject_SaveAndLoadWithContainers(t *testing.T) {
	mgr := mgrWithHost(t, "h")
	require.NoError(t, mgr.AddProject("h", nb.Project{Name: "p"}))
	require.NoError(t, mgr.AddContainer("p", nb.Container{
		Name:       "web",
		Id:         "id1",
		RemoteUser: "dev",
		Pmapping:   nb.PortMapping{nb.KSSHPortMapping: 4444},
	}))
	require.NoError(t, mgr.Save())

	mgr2 := NewInventoryManager(mgr.source, mgr.ref)
	require.NoError(t, mgr2.Load())

	c, err := mgr2.GetContainer("p", "web")
	require.NoError(t, err)
	assert.Equal(t, "id1", c.Id)
	assert.Equal(t, "dev", string(c.RemoteUser))
	assert.Equal(t, 4444, c.Pmapping[nb.KSSHPortMapping])
}
