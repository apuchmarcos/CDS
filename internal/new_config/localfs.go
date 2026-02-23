package new_config

import (
	"path/filepath"

	"github.com/amadeusitgroup/cds/internal/cenv"
	"github.com/amadeusitgroup/cds/internal/cos"
	cg "github.com/amadeusitgroup/cds/internal/global"
)

// LocalFSSource implements Source using the local filesystem via cos (afero).
// All paths are absolute.
type LocalFSSource struct{}

func (s *LocalFSSource) Read(path string) ([]byte, error) {
	return cos.ReadFile(path)
}

func (s *LocalFSSource) Write(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := cenv.EnsureDir(dir, cg.KPermDir); err != nil {
		return err
	}
	return cos.WriteFile(path, data, cg.KPermFile)
}

func (s *LocalFSSource) Exists(path string) (bool, error) {
	exists := cos.Exists(path)
	return exists, nil
}

func (s *LocalFSSource) Delete(path string) error {
	if cos.NotExist(path) {
		return nil
	}
	return cos.Fs.Remove(path)
}

// Compile-time check that LocalFSSource implements Source.
var _ Source = (*LocalFSSource)(nil)
